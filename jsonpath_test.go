package main

import (
	"encoding/json"
	"strings"
	"testing"
)

const jpDoc = `{
  "store": {
    "book": [
      {"title": "Sayings", "price": 8.99, "author": "Nigel"},
      {"title": "Moby Dick", "price": 8.99, "author": "Herman"},
      {"title": "The Odyssey", "price": 22.99, "author": "Emily"}
    ],
    "bicycle": {"color": "red", "price": 19.95}
  },
  "expensive": 10
}`

func jpQuery(t *testing.T, path string) *JSONPathResult {
	t.Helper()
	app := NewApp()
	res, err := app.JSONPathQuery(jpDoc, path)
	if err != nil {
		t.Fatalf("查询 %q 失败: %v", path, err)
	}
	return res
}

func TestJSONPathBasics(t *testing.T) {
	cases := []struct {
		path  string
		count int
		want  string // 输出中必须包含的片段
	}{
		{"$.store.bicycle.color", 1, `"red"`},
		{"$['store']['bicycle']['price']", 1, "19.95"},
		{"$.store.book[0].title", 1, "Sayings"},
		{"$.store.book[-1].title", 1, "Odyssey"},
		{"$.store.book[*].author", 3, "Herman"},
		{"$..author", 3, "Emily"},
		{"$..price", 4, "8.99"},
		{"$.store.book[0:2]", 2, "Moby Dick"},
		{"$.store.book[1:]", 2, "Odyssey"},
		{"$.store.book[-1:].title", 1, "Odyssey"},
		{"$.store.book[?(@.price<10)].title", -1, ""}, // 不支持的语法 → 报错
	}
	for _, c := range cases {
		app := NewApp()
		res, err := app.JSONPathQuery(jpDoc, c.path)
		if c.count < 0 {
			if err == nil {
				t.Errorf("查询 %q 应报语法错误", c.path)
			}
			continue
		}
		if err != nil {
			t.Errorf("查询 %q 失败: %v", c.path, err)
			continue
		}
		if res.Count != c.count {
			t.Errorf("查询 %q 匹配数 = %d, 期望 %d", c.path, res.Count, c.count)
		}
		if c.want != "" && !strings.Contains(res.Output, c.want) {
			t.Errorf("查询 %q 输出缺少 %q:\n%s", c.path, c.want, res.Output)
		}
	}
}

func TestJSONPathSingleVsArray(t *testing.T) {
	// 单个匹配直接输出值，多个匹配输出数组
	res := jpQuery(t, "$.store.bicycle.color")
	if strings.HasPrefix(strings.TrimSpace(res.Output), "[") {
		t.Errorf("单个匹配不应包装为数组: %s", res.Output)
	}
	res = jpQuery(t, "$..author")
	if !strings.HasPrefix(strings.TrimSpace(res.Output), "[") {
		t.Errorf("多个匹配应包装为数组: %s", res.Output)
	}
	// 包装后的数组必须是合法 JSON
	var arr []any
	if err := json.Unmarshal([]byte(res.Output), &arr); err != nil {
		t.Errorf("多匹配输出不是合法 JSON: %v", err)
	}
}

func TestJSONPathNoMatchAndErrors(t *testing.T) {
	app := NewApp()
	// [::-1] 展开为多个匹配且顺序反转（JSONPath 语义：切片结果展开为匹配集）
	res := jpQuery(t, "$.store.book[::-1]")
	if res.Count != 3 {
		t.Fatalf("[::-1] 应匹配 3 个, 实际 %d", res.Count)
	}
	if strings.Index(res.Output, "Odyssey") > strings.Index(res.Output, "Sayings") {
		t.Errorf("[::-1] 未反转顺序:\n%s", res.Output)
	}
	res, err := app.JSONPathQuery(jpDoc, "$.missing.deep.path")
	if err != nil {
		t.Fatalf("无匹配应正常返回: %v", err)
	}
	if res.Count != 0 || res.Output != "" {
		t.Errorf("无匹配结果异常: %+v", res)
	}
	if _, err := app.JSONPathQuery(jpDoc, "store.book"); err == nil {
		t.Error("缺少 $ 前缀应报错（或按根相对解析后仍报错）")
	}
	if _, err := app.JSONPathQuery(jpDoc, "$.store.book["); err == nil {
		t.Error("未闭合括号应报错")
	}
	if _, err := app.JSONPathQuery("{bad", "$"); err == nil {
		t.Error("非法 JSON 应报错")
	}
}

func TestJSONPathQuotedKeyWithDot(t *testing.T) {
	app := NewApp()
	res, err := app.JSONPathQuery(`{"a.b": {"c": 1}}`, `$['a.b'].c`)
	if err != nil {
		t.Fatalf("带点键查询失败: %v", err)
	}
	if res.Count != 1 || !strings.Contains(res.Output, "1") {
		t.Errorf("带点键结果异常: %+v", res)
	}
}
