package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"unicode/utf8"
)

// Member 一个键值对，用于保留 JSON 对象的原始键顺序
type Member struct {
	Key   string
	Value any
}

// Obj 保留键插入顺序的 JSON 对象。
// 解析结果中的对象一律使用该类型，避免 Go map 按字母排序破坏原始顺序。
type Obj struct {
	Members []Member
}

// MarshalJSON 输出紧凑 JSON；外层 Encoder 的 SetIndent 会负责缩进美化
func (o *Obj) MarshalJSON() ([]byte, error) {
	var b bytes.Buffer
	b.WriteByte('{')
	for i, m := range o.Members {
		if i > 0 {
			b.WriteByte(',')
		}
		k, err := marshalNoEscape(m.Key)
		if err != nil {
			return nil, err
		}
		b.Write(k)
		b.WriteByte(':')
		v, err := marshalNoEscape(m.Value)
		if err != nil {
			return nil, err
		}
		b.Write(v)
	}
	b.WriteByte('}')
	return b.Bytes(), nil
}

// marshalNoEscape 序列化单个值，且关闭 HTML 转义（保留 < > & 原样输出）
func marshalNoEscape(v any) ([]byte, error) {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return bytes.TrimRight(b.Bytes(), "\n"), nil
}

// parseOrdered 按原始顺序解析 JSON。
// 对象 -> *Obj（保序），数组 -> []any，数字 -> json.Number（保精度）。
func parseOrdered(data []byte) (any, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	val, err := decodeValue(dec)
	if err != nil {
		return nil, err
	}
	if dec.More() {
		return nil, fmt.Errorf("JSON 文档结束后存在多余内容")
	}
	return val, nil
}

func decodeValue(dec *json.Decoder) (any, error) {
	tok, err := dec.Token()
	if err != nil {
		return nil, err
	}
	switch t := tok.(type) {
	case json.Delim:
		switch t {
		case '{':
			obj := &Obj{}
			for dec.More() {
				kt, err := dec.Token()
				if err != nil {
					return nil, err
				}
				key, ok := kt.(string)
				if !ok {
					return nil, fmt.Errorf("对象键必须是字符串")
				}
				v, err := decodeValue(dec)
				if err != nil {
					return nil, err
				}
				obj.Members = append(obj.Members, Member{Key: key, Value: v})
			}
			if _, err := dec.Token(); err != nil { // 消费 '}'
				return nil, err
			}
			return obj, nil
		case '[':
			arr := []any{}
			for dec.More() {
				v, err := decodeValue(dec)
				if err != nil {
					return nil, err
				}
				arr = append(arr, v)
			}
			if _, err := dec.Token(); err != nil { // 消费 ']'
				return nil, err
			}
			return arr, nil
		default:
			return nil, fmt.Errorf("意外的分隔符 %v", t)
		}
	default:
		// string / json.Number / bool / nil
		return t, nil
	}
}

// encodeValue 序列化解析结果；indent 非空时进行缩进美化
func encodeValue(v any, indent string) (string, error) {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetEscapeHTML(false)
	if indent != "" {
		enc.SetIndent("", indent)
	}
	if err := enc.Encode(v); err != nil {
		return "", err
	}
	return strings.TrimRight(b.String(), "\n"), nil
}

// locate 将字节偏移换算为 1 起始的行号与列号（列按字符计数）
func locate(data []byte, offset int) (line, col int) {
	line, col = 1, 1
	if offset > len(data) {
		offset = len(data)
	}
	for i := 0; i < offset; {
		r, size := utf8.DecodeRune(data[i:])
		if r == '\n' {
			line++
			col = 1
		} else {
			col++
		}
		i += size
	}
	return line, col
}

// JSONStats JSON 文档的结构统计
type JSONStats struct {
	Keys    int `json:"keys"`    // 对象键总数
	Nodes   int `json:"nodes"`   // 节点总数（含容器与叶子）
	Depth   int `json:"depth"`   // 最大嵌套深度
	Objects int `json:"objects"` // 对象个数
	Arrays  int `json:"arrays"`  // 数组个数
	Strings int `json:"strings"` // 字符串值个数
	Numbers int `json:"numbers"` // 数字个数
	Bools   int `json:"bools"`   // 布尔值个数
	Nulls   int `json:"nulls"`   // null 个数
}

func collectStats(v any, depth int, s *JSONStats) {
	s.Nodes++
	switch t := v.(type) {
	case *Obj:
		s.Objects++
		if depth > s.Depth {
			s.Depth = depth
		}
		for _, m := range t.Members {
			s.Keys++
			collectStats(m.Value, depth+1, s)
		}
	case []any:
		s.Arrays++
		if depth > s.Depth {
			s.Depth = depth
		}
		for _, e := range t {
			collectStats(e, depth+1, s)
		}
	case string:
		s.Strings++
	case json.Number:
		s.Numbers++
	case bool:
		s.Bools++
	case nil:
		s.Nulls++
	}
}

// AnalyzeResult 校验 + 统计的聚合结果，供前端状态栏与编辑器错误标记使用
type AnalyzeResult struct {
	Valid   bool       `json:"valid"`
	Message string     `json:"message"`
	Line    int        `json:"line"`  // 出错行（1 起始），无错误时为 0
	Col     int        `json:"col"`   // 出错列（1 起始），无错误时为 0
	Stats   *JSONStats `json:"stats"` // 校验通过时非空
}

// AnalyzeJSON 校验 JSON 并返回统计信息；失败时携带错误行列位置
func (a *App) AnalyzeJSON(jsonStr string) AnalyzeResult {
	if strings.TrimSpace(jsonStr) == "" {
		return AnalyzeResult{Valid: false, Message: "请输入JSON内容"}
	}
	data := []byte(jsonStr)
	val, err := parseOrdered(data)
	if err != nil {
		// 注意：部分 Go 版本（如 1.25.x）中 Decoder.Token() 报告的 SyntaxError.Offset
		// 是相对当前 token 的偏移，不可靠；用 json.Unmarshal 复核以获得准确的绝对字节偏移。
		line, col := 0, 0
		var discard any
		if uerr := json.Unmarshal(data, &discard); uerr != nil {
			if use, ok := uerr.(*json.SyntaxError); ok {
				off := int(use.Offset) - 1 // Offset 指向非法字符的下一字节
				if off < 0 {
					off = 0
				}
				line, col = locate(data, off)
			}
		}
		if line == 0 {
			// Unmarshal 无法定位时，退回 Token 错误的偏移或文档末尾
			var se *json.SyntaxError
			if asSyntaxError(err, &se) {
				off := int(se.Offset) - 1
				if off < 0 {
					off = 0
				}
				line, col = locate(data, off)
			} else {
				line, col = locate(data, len(data))
			}
		}
		return AnalyzeResult{
			Valid:   false,
			Message: fmt.Sprintf("第 %d 行第 %d 列附近: %v", line, col, err),
			Line:    line,
			Col:     col,
		}
	}
	var s JSONStats
	collectStats(val, 1, &s)
	return AnalyzeResult{Valid: true, Message: "JSON格式正确 ✅", Stats: &s}
}

// asSyntaxError 在错误链中查找 *json.SyntaxError
func asSyntaxError(err error, target **json.SyntaxError) bool {
	for err != nil {
		if se, ok := err.(*json.SyntaxError); ok {
			*target = se
			return true
		}
		u, ok := err.(interface{ Unwrap() error })
		if !ok {
			return false
		}
		err = u.Unwrap()
	}
	return false
}

// FormatJSON 格式化 JSON（2 空格缩进，保留键顺序与数字精度）
func (a *App) FormatJSON(jsonStr string) (string, error) {
	if strings.TrimSpace(jsonStr) == "" {
		return "", fmt.Errorf("请输入JSON内容")
	}
	val, err := parseOrdered([]byte(jsonStr))
	if err != nil {
		return "", fmt.Errorf("JSON格式错误: %v", err)
	}
	return encodeValue(val, "  ")
}

// MinifyJSON 压缩 JSON（去除空白，保留键顺序与数字精度）
func (a *App) MinifyJSON(jsonStr string) (string, error) {
	if strings.TrimSpace(jsonStr) == "" {
		return "", fmt.Errorf("请输入JSON内容")
	}
	val, err := parseOrdered([]byte(jsonStr))
	if err != nil {
		return "", fmt.Errorf("JSON格式错误: %v", err)
	}
	return encodeValue(val, "")
}

// ValidateJSON 校验 JSON 格式（保留旧接口，供按钮调用）
func (a *App) ValidateJSON(jsonStr string) (bool, string) {
	r := a.AnalyzeJSON(jsonStr)
	return r.Valid, r.Message
}
