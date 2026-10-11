package main

import (
	"strings"
	"testing"
)

func TestBuildJSONTree(t *testing.T) {
	app := &App{}
	doc := `{"b":1,"a":[true,null,"x",9007199254740993],"o":{"k":{}}}`
	root, err := app.BuildJSONTree(doc)
	if err != nil {
		t.Fatalf("BuildJSONTree 失败: %v", err)
	}
	if root.Kind != "object" || root.Count != 3 {
		t.Fatalf("根节点异常: kind=%s count=%d", root.Kind, root.Count)
	}
	// 键顺序保留：b, a, o
	wantKeys := []string{"b", "a", "o"}
	for i, want := range wantKeys {
		if root.Children[i].Key != want {
			t.Fatalf("子节点 %d 键应为 %q，实际 %q", i, want, root.Children[i].Key)
		}
	}
	arr := root.Children[1]
	if arr.Kind != "array" || arr.Count != 4 || len(arr.Children) != 4 {
		t.Fatalf("数组节点异常: kind=%s count=%d len=%d", arr.Kind, arr.Count, len(arr.Children))
	}
	for _, el := range arr.Children {
		if el.Key != "" {
			t.Fatalf("数组元素键应为空（前端按序号渲染），实际 %q", el.Key)
		}
	}
	if got := arr.Children[0]; got.Kind != "bool" || got.Value != "true" {
		t.Fatalf("布尔元素异常: %+v", got)
	}
	if got := arr.Children[1]; got.Kind != "null" || got.Value != "" {
		t.Fatalf("null 元素异常: %+v", got)
	}
	if got := arr.Children[2]; got.Kind != "string" || got.Value != "x" {
		t.Fatalf("字符串元素异常: %+v", got)
	}
	// 大数字保精度
	if got := arr.Children[3]; got.Kind != "number" || got.Value != "9007199254740993" {
		t.Fatalf("大数字精度丢失: %+v", got)
	}
	obj := root.Children[2]
	if obj.Kind != "object" || obj.Count != 1 || obj.Children[0].Kind != "object" || obj.Children[0].Count != 0 {
		t.Fatalf("嵌套对象异常: %+v", obj)
	}
}

func TestBuildJSONTreeInvalid(t *testing.T) {
	app := &App{}
	if _, err := app.BuildJSONTree(`{"a":`); err == nil {
		t.Fatal("非法 JSON 应返回错误")
	}
}

func TestBuildJSONTreeTooLarge(t *testing.T) {
	app := &App{}
	big := strings.Repeat(" ", treeMaxChars+1) + "1" // 前导空白合法，仅触发大小守卫
	if _, err := app.BuildJSONTree(big); err == nil {
		t.Fatal("超过大小上限应返回错误")
	}
}
