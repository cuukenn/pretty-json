package main

import (
	"encoding/json"
	"fmt"
)

// TreeNode 供前端可折叠树视图渲染的 JSON 节点，保留键顺序与数字精度。
type TreeNode struct {
	Key      string      `json:"key,omitempty"`      // 对象键；数组元素由前端按序号渲染，此处留空
	Kind     string      `json:"kind"`               // object | array | string | number | bool | null
	Value    string      `json:"value,omitempty"`    // 标量的原始文本（数字保精度）
	Count    int         `json:"count,omitempty"`    // 容器成员数（折叠时用于预览）
	Children []*TreeNode `json:"children,omitempty"` // 子节点（保持原始顺序）
}

// treeMaxChars 树视图文档上限，防止过大的桥接数据拖慢界面。
const treeMaxChars = 5 << 20 // 5 MB

// BuildJSONTree 将 JSON 文本解析为树结构，复用解析缓存。
func (a *App) BuildJSONTree(jsonStr string) (*TreeNode, error) {
	if len(jsonStr) > treeMaxChars {
		return nil, fmt.Errorf("文档过大（%s），树视图上限为 %s", formatSize(int64(len(jsonStr))), formatSize(treeMaxChars))
	}
	v, err := a.parseCached(jsonStr)
	if err != nil {
		return nil, err
	}
	return buildTreeNode("", v), nil
}

func buildTreeNode(key string, v any) *TreeNode {
	switch val := v.(type) {
	case *Obj:
		node := &TreeNode{Key: key, Kind: "object", Count: len(val.Members)}
		for _, m := range val.Members {
			node.Children = append(node.Children, buildTreeNode(m.Key, m.Value))
		}
		return node
	case []any:
		node := &TreeNode{Key: key, Kind: "array", Count: len(val)}
		for _, el := range val {
			node.Children = append(node.Children, buildTreeNode("", el))
		}
		return node
	case json.Number:
		return &TreeNode{Key: key, Kind: "number", Value: val.String()}
	case string:
		return &TreeNode{Key: key, Kind: "string", Value: val}
	case bool:
		return &TreeNode{Key: key, Kind: "bool", Value: fmt.Sprintf("%t", val)}
	default: // nil → JSON null
		return &TreeNode{Key: key, Kind: "null"}
	}
}
