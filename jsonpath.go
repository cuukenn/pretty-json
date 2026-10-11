package main

import (
	"fmt"
	"strconv"
	"strings"
)

// JSONPathResult JSONPath 查询结果
type JSONPathResult struct {
	Count  int    `json:"count"`
	Output string `json:"output"` // 单个匹配直接输出该值；多个匹配输出 JSON 数组
}

// ---------- 路径解析 ----------

type jpSegment struct {
	kind      int // 0=键名, 1=通配符, 2=数组下标, 3=切片
	recursive bool
	name      string
	index     int
	slice     [3]*int // start, end, step（nil 表示省略）
}

// parseJSONPath 解析 JSONPath 子集：$ .name ..name ['k'] [n] [*] [a:b:step]
func parseJSONPath(path string) ([]jpSegment, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, fmt.Errorf("JSONPath 不能为空")
	}
	if path[0] == '$' {
		path = path[1:]
	}
	var segs []jpSegment
	for path != "" {
		switch path[0] {
		case '.':
			path = path[1:]
			recursive := false
			if strings.HasPrefix(path, ".") { // ".." 递归下降
				recursive = true
				path = path[1:]
			}
			if path == "" {
				return nil, fmt.Errorf("JSONPath 以 '.' 结尾")
			}
			if path[0] == '*' {
				segs = append(segs, jpSegment{kind: 1, recursive: recursive})
				path = path[1:]
				continue
			}
			if path[0] == '[' { // $..['k'] 形式
				continue
			}
			// 读取标识符（允许字母数字下划线与中划线，直到下一个分隔符）
			i := 0
			for i < len(path) && path[i] != '.' && path[i] != '[' {
				i++
			}
			name := path[:i]
			if name == "" {
				return nil, fmt.Errorf("JSONPath 语法错误：空的键名")
			}
			segs = append(segs, jpSegment{kind: 0, name: name, recursive: recursive})
			path = path[i:]
		case '[':
			end := jpBracketEnd(path)
			if end < 0 {
				return nil, fmt.Errorf("JSONPath 缺少匹配的 ']'")
			}
			inner := strings.TrimSpace(path[1:end])
			path = path[end+1:]
			switch {
			case inner == "*":
				segs = append(segs, jpSegment{kind: 1})
			case strings.HasPrefix(inner, "'") || strings.HasPrefix(inner, "\""):
				q := inner[0]
				if len(inner) < 2 || inner[len(inner)-1] != q {
					return nil, fmt.Errorf("JSONPath 引号未闭合: %s", inner)
				}
				segs = append(segs, jpSegment{kind: 0, name: inner[1 : len(inner)-1]})
			case strings.Contains(inner, ":"):
				seg, err := jpParseSlice(inner)
				if err != nil {
					return nil, err
				}
				segs = append(segs, seg)
			default:
				n, err := strconv.Atoi(inner)
				if err != nil {
					return nil, fmt.Errorf("JSONPath 中括号内无法解析: %s", inner)
				}
				segs = append(segs, jpSegment{kind: 2, index: n})
			}
		default:
			return nil, fmt.Errorf("JSONPath 语法错误：意外的字符 %q", path[0])
		}
	}
	return segs, nil
}

// jpBracketEnd 返回与开头 '[' 匹配的 ']' 下标（跳过引号内的内容）
func jpBracketEnd(s string) int {
	depth := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '\'', '"':
			q := s[i]
			for i++; i < len(s); i++ {
				if s[i] == '\\' {
					i++
					continue
				}
				if s[i] == q {
					break
				}
			}
		case '[':
			depth++
		case ']':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

func jpParseSlice(inner string) (jpSegment, error) {
	parts := strings.Split(inner, ":")
	if len(parts) > 3 {
		return jpSegment{}, fmt.Errorf("JSONPath 切片语法错误: %s", inner)
	}
	seg := jpSegment{kind: 3}
	parseInt := func(s string) (*int, error) {
		s = strings.TrimSpace(s)
		if s == "" {
			return nil, nil
		}
		n, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("JSONPath 切片语法错误: %s", inner)
		}
		return &n, nil
	}
	var err error
	if seg.slice[0], err = parseInt(parts[0]); err != nil {
		return seg, err
	}
	if seg.slice[1], err = parseInt(parts[1]); err != nil {
		return seg, err
	}
	if len(parts) == 3 {
		if seg.slice[2], err = parseInt(parts[2]); err != nil {
			return seg, err
		}
	}
	return seg, nil
}

// ---------- 求值 ----------

// jpWalk 收集节点自身及全部后代（用于 .. 递归下降）
func jpWalk(v any, out *[]any) {
	*out = append(*out, v)
	switch t := v.(type) {
	case *Obj:
		for _, m := range t.Members {
			jpWalk(m.Value, out)
		}
	case []any:
		for _, e := range t {
			jpWalk(e, out)
		}
	}
}

func jpApply(seg jpSegment, cur []any) []any {
	if seg.recursive {
		var all []any
		for _, n := range cur {
			jpWalk(n, &all)
		}
		cur = all
	}
	var next []any
	for _, n := range cur {
		switch seg.kind {
		case 0: // 键名
			if obj, ok := n.(*Obj); ok {
				for _, m := range obj.Members {
					if m.Key == seg.name {
						next = append(next, m.Value)
					}
				}
			}
		case 1: // 通配符
			switch t := n.(type) {
			case *Obj:
				for _, m := range t.Members {
					next = append(next, m.Value)
				}
			case []any:
				next = append(next, t...)
			}
		case 2: // 下标
			if arr, ok := n.([]any); ok {
				i := seg.index
				if i < 0 {
					i += len(arr)
				}
				if i >= 0 && i < len(arr) {
					next = append(next, arr[i])
				}
			}
		case 3: // 切片
			if arr, ok := n.([]any); ok {
				next = append(next, jpSlice(arr, seg.slice)...)
			}
		}
	}
	return next
}

// jpSlice 按 Python 风格切片（支持负索引与步长；负步长时默认从末尾开始）
func jpSlice(arr []any, sl [3]*int) []any {
	n := len(arr)
	if n == 0 {
		return nil
	}
	var out []any
	step := 1
	if sl[2] != nil {
		step = *sl[2]
	}
	if step == 0 {
		return nil
	}
	if step > 0 {
		start, end := 0, n
		if sl[0] != nil {
			start = *sl[0]
			if start < 0 {
				start += n
			}
			if start < 0 {
				start = 0
			}
		}
		if sl[1] != nil {
			end = *sl[1]
			if end < 0 {
				end += n
			}
			if end > n {
				end = n
			}
		}
		for i := start; i < end; i += step {
			out = append(out, arr[i])
		}
		return out
	}
	// 负步长：默认 start=n-1, end=-1（遍历到索引 0）
	start, end := n-1, -1
	if sl[0] != nil {
		start = *sl[0]
		if start < 0 {
			start += n
		}
		if start > n-1 {
			start = n - 1
		}
	}
	if sl[1] != nil {
		end = *sl[1]
		if end < 0 {
			end += n
		}
		if end < -1 {
			end = -1
		}
	}
	for i := start; i > end; i += step {
		if i >= 0 {
			out = append(out, arr[i])
		}
	}
	return out
}

// JSONPathQuery 在 JSON 文本上执行 JSONPath 查询
func (a *App) JSONPathQuery(jsonStr, path string) (*JSONPathResult, error) {
	if strings.TrimSpace(jsonStr) == "" {
		return nil, fmt.Errorf("请输入JSON内容")
	}
	val, err := parseOrdered([]byte(jsonStr))
	if err != nil {
		return nil, fmt.Errorf("JSON格式错误: %v", err)
	}
	segs, err := parseJSONPath(path)
	if err != nil {
		return nil, err
	}
	matches := []any{val}
	for _, s := range segs {
		matches = jpApply(s, matches)
	}
	res := &JSONPathResult{Count: len(matches)}
	switch len(matches) {
	case 0:
		res.Output = ""
	case 1:
		out, err := encodeValue(matches[0], "  ")
		if err != nil {
			return nil, err
		}
		res.Output = out
	default:
		out, err := encodeValue(matches, "  ")
		if err != nil {
			return nil, err
		}
		res.Output = out
	}
	return res, nil
}
