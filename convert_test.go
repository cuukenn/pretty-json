package main

import (
	"strings"
	"testing"
)

func TestJSONToYAMLOrderAndPrecision(t *testing.T) {
	app := NewApp()
	in := `{"z": 123456789012345678901234567890, "a": 1.50, "s": "1.5", "b": true, "n": null, "arr": [1, {"k": "v"}], "nested": {"y": 1, "x": 2}}`
	out, err := app.JSONToYAML(in)
	if err != nil {
		t.Fatalf("JSONToYAML 失败: %v", err)
	}
	// 键顺序保留
	iZ, iA := strings.Index(out, "z:"), strings.Index(out, "a:")
	if iZ < 0 || iA < 0 || iZ > iA {
		t.Errorf("键顺序未保留:\n%s", out)
	}
	// 大数字与原始精度
	if !strings.Contains(out, "123456789012345678901234567890") {
		t.Errorf("大数字精度丢失:\n%s", out)
	}
	if !strings.Contains(out, "1.50") {
		t.Errorf("浮点原始精度丢失:\n%s", out)
	}
	// 数字字符串必须加引号，避免被误读为数字
	if !strings.Contains(out, `"1.5"`) && !strings.Contains(out, "'1.5'") {
		t.Errorf("数字字符串未加引号:\n%s", out)
	}
	// 嵌套映射缩进
	if !strings.Contains(out, "nested:") || !strings.Contains(out, "y: 1") {
		t.Errorf("嵌套结构缺失:\n%s", out)
	}
}

func TestYAMLToJSONRoundTrip(t *testing.T) {
	app := NewApp()
	yamlIn := "zeta: 42\nalpha: true\nlist:\n  - 1\n  - 2\nnested:\n  b: [1.5, \"x\"]\n  a: null\n"
	jsonOut, err := app.YAMLToJSON(yamlIn)
	if err != nil {
		t.Fatalf("YAMLToJSON 失败: %v", err)
	}
	// 键顺序保留：zeta 在 alpha 前
	if strings.Index(jsonOut, "zeta") > strings.Index(jsonOut, "alpha") {
		t.Errorf("YAML 键顺序未保留: %s", jsonOut)
	}
	if !strings.Contains(jsonOut, "true") {
		t.Errorf("YAML 布尔 true 未转换: %s", jsonOut)
	}
	// 回转 YAML 应合法
	back, err := app.JSONToYAML(jsonOut)
	if err != nil {
		t.Fatalf("回转 JSONToYAML 失败: %v", err)
	}
	if !strings.Contains(back, "zeta: 42") {
		t.Errorf("回转结果异常:\n%s", back)
	}
}

func TestJSONToTOML(t *testing.T) {
	app := NewApp()
	in := `{"title": "demo", "count": 3, "rate": 1.25, "on": true, "server": {"host": "localhost", "port": 8080}, "items": [{"n": 1}, {"n": 2}]}`
	out, err := app.JSONToTOML(in)
	if err != nil {
		t.Fatalf("JSONToTOML 失败: %v", err)
	}
	for _, want := range []string{`title = "demo"`, "count = 3", "[server]", "host = \"localhost\"", "[[items]]"} {
		if !strings.Contains(out, want) {
			t.Errorf("TOML 输出缺少 %q:\n%s", want, out)
		}
	}

	// null 不支持
	if _, err := app.JSONToTOML(`{"a": null}`); err == nil {
		t.Error("null 值应报错")
	}
	// 顶层非对象不支持
	if _, err := app.JSONToTOML(`[1,2]`); err == nil {
		t.Error("顶层数组应报错")
	}
}

func TestTOMLToJSONRoundTrip(t *testing.T) {
	app := NewApp()
	tomlIn := "b = 2\na = \"x\"\n[tbl]\nk = 1.5\n"
	jsonOut, err := app.TOMLToJSON(tomlIn)
	if err != nil {
		t.Fatalf("TOMLToJSON 失败: %v", err)
	}
	// 键按字母序稳定输出
	if strings.Index(jsonOut, `"a"`) > strings.Index(jsonOut, `"b"`) {
		t.Errorf("TOML 转 JSON 键应排序: %s", jsonOut)
	}
	if !strings.Contains(jsonOut, "1.5") || !strings.Contains(jsonOut, "tbl") {
		t.Errorf("TOML 转 JSON 内容缺失: %s", jsonOut)
	}
	// 回转 TOML 再解析应等价
	back, err := app.JSONToTOML(jsonOut)
	if err != nil {
		t.Fatalf("回转 JSONToTOML 失败: %v", err)
	}
	again, err := app.TOMLToJSON(back)
	if err != nil {
		t.Fatalf("再解析失败: %v", err)
	}
	if again != jsonOut {
		t.Errorf("往返不一致:\n%s\n---\n%s", jsonOut, again)
	}
}

func TestConvertEmptyAndBadInput(t *testing.T) {
	app := NewApp()
	if _, err := app.JSONToYAML("  "); err == nil {
		t.Error("空输入应报错")
	}
	if _, err := app.JSONToYAML("{oops"); err == nil {
		t.Error("非法 JSON 应报错")
	}
	if _, err := app.YAMLToJSON("a: [1,\n b"); err == nil {
		t.Error("非法 YAML 应报错")
	}
	if _, err := app.TOMLToJSON("a = = 1"); err == nil {
		t.Error("非法 TOML 应报错")
	}
}
