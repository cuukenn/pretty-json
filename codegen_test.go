package main

import (
	"strings"
	"testing"
)

const cgDoc = `{
  "user_id": 42,
  "name": "alice",
  "score": 9.5,
  "active": true,
  "nickname": null,
  "tags": ["a", "b"],
  "profile": {"age": 18, "city": "hz"},
  "items": [
    {"id": 1, "sku": "x"},
    {"id": 2, "sku": "y", "note": "extra"}
  ],
  "mixed": [1, "two"]
}`

func TestGenerateGo(t *testing.T) {
	app := NewApp()
	out, err := app.GenerateCode(cgDoc, "go")
	if err != nil {
		t.Fatalf("Go 生成失败: %v", err)
	}
	for _, want := range []string{
		"type Root struct {",
		"UserID int64 `json:\"user_id\"`",
		"Score float64 `json:\"score\"`",
		"Nickname interface{} `json:\"nickname\"`", // null → interface{}
		"Tags []string `json:\"tags\"`",
		"Items []Item `json:\"items\"`", // 数组对象元素单数化命名
		"type Item struct {",
		"Note string `json:\"note\"`", // 部分元素缺失的字段仍按已有类型生成
		"Profile Profile `json:\"profile\"`",
		"Mixed []interface{} `json:\"mixed\"`", // 混合类型数组
	} {
		if !strings.Contains(out, want) {
			t.Errorf("Go 输出缺少 %q:\n%s", want, out)
		}
	}
}

func TestGenerateTypeScript(t *testing.T) {
	app := NewApp()
	out, err := app.GenerateCode(cgDoc, "typescript")
	if err != nil {
		t.Fatalf("TS 生成失败: %v", err)
	}
	for _, want := range []string{
		"interface Root {",
		"userId: number;",
		"nickname: null;",
		"tags: string[];",
		"items: Item[];",
		"interface Item {",
		"note?: string;", // 部分元素缺失的键 → 可选
		"mixed: any[];",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("TS 输出缺少 %q:\n%s", want, out)
		}
	}
}

func TestGenerateJava(t *testing.T) {
	app := NewApp()
	out, err := app.GenerateCode(cgDoc, "java")
	if err != nil {
		t.Fatalf("Java 生成失败: %v", err)
	}
	for _, want := range []string{
		"public class Root {",
		"public long userId;",
		"public List<String> tags;",
		"public List<Item> items;",
		"public Object nickname;",
		"public static class Item {",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("Java 输出缺少 %q:\n%s", want, out)
		}
	}
	// Java 顶层非对象应报错
	if _, err := app.GenerateCode(`[1,2]`, "java"); err == nil {
		t.Error("Java 顶层数组应报错")
	}
}

func TestGenerateDedupeAndErrors(t *testing.T) {
	app := NewApp()
	// 两个分支下同构对象应复用同名类型且只定义一次
	out, err := app.GenerateCode(`{"a": {"items": [{"x":1}]}, "b": {"items": [{"x":1}]}}`, "go")
	if err != nil {
		t.Fatalf("Go 生成失败: %v", err)
	}
	if n := strings.Count(out, "type Item struct"); n != 1 {
		t.Errorf("同构类型应去重，出现 %d 次:\n%s", n, out)
	}
	if _, err := app.GenerateCode(cgDoc, "rust"); err == nil {
		t.Error("不支持的语言应报错")
	}
	if _, err := app.GenerateCode("{bad", "go"); err == nil {
		t.Error("非法 JSON 应报错")
	}
	// 顶层标量/数组：Go/TS 用类型别名
	out, err = app.GenerateCode(`[{"k": 1}]`, "go")
	if err != nil {
		t.Fatalf("顶层数组生成失败: %v", err)
	}
	if !strings.Contains(out, "type Root = []RootItem") || !strings.Contains(out, "type RootItem struct") {
		t.Errorf("顶层数组应生成别名+元素类型:\n%s", out)
	}
}
