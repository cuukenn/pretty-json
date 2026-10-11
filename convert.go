package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

// ---------- JSON -> YAML ----------

// toYAMLNode 将保序解析树转换为 yaml.Node，保留对象键顺序
func toYAMLNode(v any) (*yaml.Node, error) {
	switch t := v.(type) {
	case *Obj:
		node := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
		for _, m := range t.Members {
			k := &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: m.Key}
			vn, err := toYAMLNode(m.Value)
			if err != nil {
				return nil, err
			}
			node.Content = append(node.Content, k, vn)
		}
		return node, nil
	case []any:
		node := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
		for _, e := range t {
			en, err := toYAMLNode(e)
			if err != nil {
				return nil, err
			}
			node.Content = append(node.Content, en)
		}
		return node, nil
	case json.Number:
		return numberNode(string(t)), nil
	case string:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: t}, nil
	case bool:
		val := "false"
		if t {
			val = "true"
		}
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!bool", Value: val}, nil
	case nil:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!null", Value: "null"}, nil
	default:
		return nil, fmt.Errorf("不支持的类型 %T", v)
	}
}

// numberNode 按字面量输出数字，保持原始精度（整数用 !!int，其余用 !!float）
func numberNode(lit string) *yaml.Node {
	tag := "!!float"
	if _, ok := new(big.Int).SetString(lit, 10); ok {
		tag = "!!int"
	}
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: tag, Value: lit}
}

// JSONToYAML 将 JSON 文本转换为 YAML（2 空格缩进，保留键顺序与数字精度）
func (a *App) JSONToYAML(jsonStr string) (string, error) {
	if strings.TrimSpace(jsonStr) == "" {
		return "", fmt.Errorf("请输入JSON内容")
	}
	val, err := a.parseCached(jsonStr)
	if err != nil {
		return "", fmt.Errorf("JSON格式错误: %v", err)
	}
	node, err := toYAMLNode(val)
	if err != nil {
		return "", err
	}
	var b strings.Builder
	enc := yaml.NewEncoder(&b)
	enc.SetIndent(2)
	if err := enc.Encode(node); err != nil {
		return "", fmt.Errorf("生成YAML失败: %v", err)
	}
	if err := enc.Close(); err != nil {
		return "", fmt.Errorf("生成YAML失败: %v", err)
	}
	return strings.TrimRight(b.String(), "\n"), nil
}

// ---------- YAML -> JSON ----------

// fromYAMLNode 将 yaml.Node 解析树转换为保序 JSON 结构
func fromYAMLNode(n *yaml.Node) (any, error) {
	switch n.Kind {
	case yaml.DocumentNode:
		if len(n.Content) == 0 {
			return nil, fmt.Errorf("YAML 文档为空")
		}
		return fromYAMLNode(n.Content[0])
	case yaml.MappingNode:
		obj := &Obj{}
		for i := 0; i+1 < len(n.Content); i += 2 {
			k := n.Content[i]
			if k.Kind != yaml.ScalarNode {
				return nil, fmt.Errorf("不支持的复杂键（第 %d 行）", k.Line)
			}
			v, err := fromYAMLNode(n.Content[i+1])
			if err != nil {
				return nil, err
			}
			obj.Members = append(obj.Members, Member{Key: k.Value, Value: v})
		}
		return obj, nil
	case yaml.SequenceNode:
		arr := make([]any, 0, len(n.Content))
		for _, c := range n.Content {
			v, err := fromYAMLNode(c)
			if err != nil {
				return nil, err
			}
			arr = append(arr, v)
		}
		return arr, nil
	case yaml.AliasNode:
		if n.Alias == nil {
			return nil, fmt.Errorf("无效的别名引用（第 %d 行）", n.Line)
		}
		return fromYAMLNode(n.Alias)
	case yaml.ScalarNode:
		switch n.Tag {
		case "!!str", "!!binary":
			return n.Value, nil
		case "!!int", "!!float":
			return json.Number(n.Value), nil
		case "!!bool":
			return yamlBool(n.Value), nil
		case "!!null":
			return nil, nil
		default:
			// 未知标签：按字符串处理
			return n.Value, nil
		}
	default:
		return nil, fmt.Errorf("不支持的 YAML 节点类型（第 %d 行）", n.Line)
	}
}

// yamlBool 兼容 YAML 1.1 的多种布尔写法
func yamlBool(s string) bool {
	switch strings.ToLower(s) {
	case "true", "yes", "on", "y":
		return true
	default:
		return false
	}
}

// YAMLToJSON 将 YAML 文本转换为 JSON（2 空格缩进，保留键顺序）
func (a *App) YAMLToJSON(yamlStr string) (string, error) {
	if strings.TrimSpace(yamlStr) == "" {
		return "", fmt.Errorf("请输入YAML内容")
	}
	var node yaml.Node
	if err := yaml.Unmarshal([]byte(yamlStr), &node); err != nil {
		return "", fmt.Errorf("YAML格式错误: %v", err)
	}
	val, err := fromYAMLNode(&node)
	if err != nil {
		return "", err
	}
	return encodeValue(val, "  ")
}

// ---------- JSON -> TOML ----------

// toTOMLValue 将保序解析树转换为 toml 编码器可处理的值
func toTOMLValue(v any) (any, error) {
	switch t := v.(type) {
	case *Obj:
		m := make(map[string]any, len(t.Members))
		for _, mem := range t.Members {
			if mem.Value == nil {
				return nil, fmt.Errorf("TOML 不支持 null 值（键 %q）", mem.Key)
			}
			cv, err := toTOMLValue(mem.Value)
			if err != nil {
				return nil, err
			}
			m[mem.Key] = cv
		}
		return m, nil
	case []any:
		arr := make([]any, 0, len(t))
		for _, e := range t {
			if e == nil {
				return nil, fmt.Errorf("TOML 不支持 null 值（数组元素）")
			}
			cv, err := toTOMLValue(e)
			if err != nil {
				return nil, err
			}
			arr = append(arr, cv)
		}
		return arr, nil
	case json.Number:
		if i, err := t.Int64(); err == nil {
			return i, nil
		}
		f, err := t.Float64()
		if err != nil || math.IsInf(f, 0) || math.IsNaN(f) {
			return nil, fmt.Errorf("数字 %s 超出 TOML 支持范围", t.String())
		}
		return f, nil
	case string:
		return t, nil
	case bool:
		return t, nil
	default:
		return nil, fmt.Errorf("不支持的类型 %T", v)
	}
}

// JSONToTOML 将 JSON 文本转换为 TOML（顶层必须是对象，不支持 null）
func (a *App) JSONToTOML(jsonStr string) (string, error) {
	if strings.TrimSpace(jsonStr) == "" {
		return "", fmt.Errorf("请输入JSON内容")
	}
	val, err := a.parseCached(jsonStr)
	if err != nil {
		return "", fmt.Errorf("JSON格式错误: %v", err)
	}
	obj, ok := val.(*Obj)
	if !ok {
		return "", fmt.Errorf("TOML 顶层必须是对象，当前是 %T", val)
	}
	m, err := toTOMLValue(obj)
	if err != nil {
		return "", err
	}
	var b strings.Builder
	enc := toml.NewEncoder(&b)
	if err := enc.Encode(m); err != nil {
		return "", fmt.Errorf("生成TOML失败: %v", err)
	}
	return strings.TrimRight(b.String(), "\n"), nil
}

// ---------- TOML -> JSON ----------

// fromTOMLValue 将 toml 解码结果转换为保序 JSON 结构（对象键按字母排序以保证稳定输出）
func fromTOMLValue(v any) any {
	switch t := v.(type) {
	case map[string]any:
		keys := make([]string, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		obj := &Obj{}
		for _, k := range keys {
			obj.Members = append(obj.Members, Member{Key: k, Value: fromTOMLValue(t[k])})
		}
		return obj
	case map[string]string: // toml 对某些基本类型表返回的强类型 map
		keys := make([]string, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		obj := &Obj{}
		for _, k := range keys {
			obj.Members = append(obj.Members, Member{Key: k, Value: t[k]})
		}
		return obj
	case []any:
		arr := make([]any, 0, len(t))
		for _, e := range t {
			arr = append(arr, fromTOMLValue(e))
		}
		return arr
	case int64:
		return json.Number(strconv.FormatInt(t, 10))
	case float64:
		return json.Number(strconv.FormatFloat(t, 'g', -1, 64))
	case bool:
		return t
	case string:
		return t
	case time.Time:
		return t.Format(time.RFC3339Nano)
	default:
		return fmt.Sprint(t)
	}
}

// TOMLToJSON 将 TOML 文本转换为 JSON（2 空格缩进）
func (a *App) TOMLToJSON(tomlStr string) (string, error) {
	if strings.TrimSpace(tomlStr) == "" {
		return "", fmt.Errorf("请输入TOML内容")
	}
	var m map[string]any
	if _, err := toml.Decode(tomlStr, &m); err != nil {
		return "", fmt.Errorf("TOML格式错误: %v", err)
	}
	return encodeValue(fromTOMLValue(m), "  ")
}
