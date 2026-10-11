package main

import (
	"strings"
	"testing"
)

const cgDoc = `{
	"userId": 7,
	"api_key": "k",
	"active": true,
	"score": 1.5,
	"meta": null,
	"tags": ["a", "b"],
	"items": [{"id": 1, "name": "x"}, {"id": 2, "name": "y"}],
	"address": {"city": "SH", "zip": "200000"}
}`

func TestGenerateGo(t *testing.T) {
	app := NewApp()
	out, err := app.GenerateCode(cgDoc, "go", nil)
	if err != nil {
		t.Fatalf("Go 生成失败: %v", err)
	}
	for _, want := range []string{
		"type Root struct",
		"UserId int64 `json:\"userId\"`",
		"APIKey string `json:\"api_key\"`", // 缩写词大写
		"Active bool `json:\"active\"`",
		"Score float64 `json:\"score\"`",
		"Meta interface{} `json:\"meta\"`",
		"Tags []string `json:\"tags\"`",
		"Items []Item `json:\"items\"`", // 数组元素命名用单数
		"Address Address `json:\"address\"`",
		"type Item struct",
		"ID int64 `json:\"id\"`",
		"type Address struct",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("Go 输出缺少 %q\n%s", want, out)
		}
	}
}

func TestGenerateTS(t *testing.T) {
	app := NewApp()
	out, err := app.GenerateCode(cgDoc, "typescript", nil)
	if err != nil {
		t.Fatalf("TS 生成失败: %v", err)
	}
	for _, want := range []string{
		"interface Root {",
		"userId: number;",
		"apiKey: string;",
		"active: boolean;",
		"meta: null;",
		"tags: string[];",
		"items: Item[];",
		"interface Item {",
		"id: number;",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("TS 输出缺少 %q\n%s", want, out)
		}
	}
}

func TestGenerateJava(t *testing.T) {
	app := NewApp()
	out, err := app.GenerateCode(cgDoc, "java", nil)
	if err != nil {
		t.Fatalf("Java 生成失败: %v", err)
	}
	for _, want := range []string{
		"import com.fasterxml.jackson.annotation.JsonIgnoreProperties;",
		"import com.fasterxml.jackson.annotation.JsonProperty;",
		"import java.util.List;",
		"@JsonIgnoreProperties(ignoreUnknown = true)",
		"public class Root {",
		"@JsonProperty(\"userId\")",
		"private long userId;",
		"public long getUserId() {",
		"public void setUserId(long userId) {",
		"private String apiKey;",
		"public boolean isActive() {", // boolean 基本类型用 is 前缀
		"public void setActive(boolean active) {",
		"private List<String> tags;",
		"private List<Item> items;",
		"    @JsonIgnoreProperties(ignoreUnknown = true)\n    public static class Item {", // 内部类注解缩进
	} {
		if !strings.Contains(out, want) {
			t.Errorf("Java 输出缺少 %q\n%s", want, out)
		}
	}
	// 无 List 字段时不应导入 java.util.List
	out, err = app.GenerateCode(`{"a": 1}`, "java", nil)
	if err != nil {
		t.Fatalf("Java 生成失败: %v", err)
	}
	if strings.Contains(out, "import java.util.List;") {
		t.Errorf("无 List 字段不应导入 java.util.List:\n%s", out)
	}
	// 顶层非对象应报错
	if _, err := app.GenerateCode(`[1,2]`, "java", nil); err == nil {
		t.Error("顶层数组应报错")
	}
}

func TestGenerateJavaLombok(t *testing.T) {
	app := NewApp()
	out, err := app.GenerateCode(cgDoc, "java", map[string]any{"javaLombok": true})
	if err != nil {
		t.Fatalf("Java(Lombok) 生成失败: %v", err)
	}
	for _, want := range []string{
		"import lombok.Data;",
		"@Data",
		"@JsonIgnoreProperties(ignoreUnknown = true)",
		"@JsonProperty(\"userId\")",
		"private long userId;",
		"    @Data\n    @JsonIgnoreProperties(ignoreUnknown = true)\n    public static class Item {",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("Lombok 输出缺少 %q\n%s", want, out)
		}
	}
	for _, bad := range []string{"public long getUserId()", "public void setUserId", "public boolean isActive()"} {
		if strings.Contains(out, bad) {
			t.Errorf("Lombok 模式不应生成访问器 %q\n%s", bad, out)
		}
	}
	// opts 为 nil 或 false 时仍生成显式访问器
	out, _ = app.GenerateCode(cgDoc, "java", map[string]any{"javaLombok": false})
	if !strings.Contains(out, "public long getUserId()") {
		t.Errorf("javaLombok=false 应生成访问器\n%s", out)
	}
}

func TestGenerateDedupeAndErrors(t *testing.T) {
	app := NewApp()
	// 两个分支下同构对象应复用同名类型且只定义一次
	out, err := app.GenerateCode(`{"a": {"items": [{"x":1}]}, "b": {"items": [{"x":1}]}}`, "go", nil)
	if err != nil {
		t.Fatalf("Go 生成失败: %v", err)
	}
	if n := strings.Count(out, "type Item struct"); n != 1 {
		t.Errorf("同构类型应去重，出现 %d 次\n%s", n, out)
	}
	if _, err := app.GenerateCode(cgDoc, "rust", nil); err == nil {
		t.Error("不支持的语言应报错")
	}
	if _, err := app.GenerateCode("{bad", "go", nil); err == nil {
		t.Error("非法 JSON 应报错")
	}
	// 顶层标量/数组：Go/TS 用类型别名
	out, err = app.GenerateCode(`[{"k": 1}]`, "go", nil)
	if err != nil {
		t.Fatalf("顶层数组生成失败: %v", err)
	}
	if !strings.Contains(out, "type Root = []RootItem") || !strings.Contains(out, "type RootItem struct") {
		t.Errorf("顶层数组应生成别名+元素类型:\n%s", out)
	}
}
