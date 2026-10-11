package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestFormatJSONKeepsKeyOrder(t *testing.T) {
	app := NewApp()
	in := `{"zebra":1,"apple":2,"mango":3}`
	out, err := app.FormatJSON(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"zebra"`) {
		t.Fatalf("missing key in output: %s", out)
	}
	zi := strings.Index(out, "zebra")
	ai := strings.Index(out, "apple")
	mi := strings.Index(out, "mango")
	if !(zi < ai && ai < mi) {
		t.Fatalf("key order not preserved: %s", out)
	}
	if !strings.Contains(out, "\n  ") {
		t.Fatalf("expected 2-space indent: %s", out)
	}
}

func TestFormatJSONPreservesBigNumberPrecision(t *testing.T) {
	app := NewApp()
	in := `{"id":12345678901234567890,"pi":3.141592653589793238462643383279}`
	out, err := app.FormatJSON(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "12345678901234567890") {
		t.Fatalf("big integer lost precision: %s", out)
	}
	if !strings.Contains(out, "3.141592653589793238462643383279") {
		t.Fatalf("high-precision decimal lost: %s", out)
	}
}

func TestFormatJSONNestedAndDuplicateKeys(t *testing.T) {
	app := NewApp()
	in := `{"b":{"y":[1,2,{"q":true}],"x":null},"a":"<html> & \"转义\"","b2":[]}`
	out, err := app.FormatJSON(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, `\u003c`) || !strings.Contains(out, "<html>") {
		t.Fatalf("HTML should not be escaped: %s", out)
	}
	if !strings.Contains(out, "[]") {
		t.Fatalf("empty array should stay compact: %s", out)
	}
}

func TestMinifyJSON(t *testing.T) {
	app := NewApp()
	in := "{\n  \"a\": 1,\n  \"b\": [ 1, 2 ]\n}"
	out, err := app.MinifyJSON(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != `{"a":1,"b":[1,2]}` {
		t.Fatalf("minify wrong: %s", out)
	}
}

func TestAnalyzeJSONErrorPosition(t *testing.T) {
	app := NewApp()
	in := "{\n  \"a\": 1,\n  \"b\": oops\n}"
	r := app.AnalyzeJSON(in)
	if r.Valid {
		t.Fatal("expected invalid")
	}
	if r.Line != 3 {
		t.Fatalf("expected error on line 3, got %d (%s)", r.Line, r.Message)
	}
	if r.Col != 8 {
		t.Fatalf("expected col 8 (position of 'oops'), got %d", r.Col)
	}
	if r.Stats != nil {
		t.Fatal("stats should be nil on error")
	}
}

func TestAnalyzeJSONTruncated(t *testing.T) {
	app := NewApp()
	r := app.AnalyzeJSON(`{"a": [1, 2`)
	if r.Valid {
		t.Fatal("expected invalid")
	}
	if r.Line == 0 || r.Col == 0 {
		t.Fatalf("expected a position for truncated input, got line %d col %d (%s)", r.Line, r.Col, r.Message)
	}
}

func TestAnalyzeJSONStats(t *testing.T) {
	app := NewApp()
	in := `{"a":{"b":[1,"x",true,null]},"c":{}}`
	r := app.AnalyzeJSON(in)
	if !r.Valid {
		t.Fatalf("expected valid: %s", r.Message)
	}
	s := r.Stats
	if s == nil {
		t.Fatal("stats missing")
	}
	// 键: a,b,c = 3；节点: 根对象+a对象+数组+1+x+true+null+c对象 = 8
	if s.Keys != 3 {
		t.Errorf("keys = %d, want 3", s.Keys)
	}
	if s.Nodes != 8 {
		t.Errorf("nodes = %d, want 8", s.Nodes)
	}
	if s.Depth != 3 {
		t.Errorf("depth = %d, want 3", s.Depth)
	}
	if s.Objects != 3 || s.Arrays != 1 || s.Strings != 1 || s.Numbers != 1 || s.Bools != 1 || s.Nulls != 1 {
		t.Errorf("type counts wrong: %+v", s)
	}
}

func TestAnalyzeJSONEmptyAndTrailing(t *testing.T) {
	app := NewApp()
	if r := app.AnalyzeJSON("   "); r.Valid || r.Line != 0 {
		t.Fatalf("empty input handling wrong: %+v", r)
	}
	if r := app.AnalyzeJSON(`{} {}`); r.Valid {
		t.Fatal("trailing content should be invalid")
	}
}

func TestValidateJSONBackwardCompatible(t *testing.T) {
	app := NewApp()
	ok, msg := app.ValidateJSON(`{"a":1}`)
	if !ok || !strings.Contains(msg, "正确") {
		t.Fatalf("valid case wrong: %v %s", ok, msg)
	}
	ok, _ = app.ValidateJSON(`{`)
	if ok {
		t.Fatal("invalid case should return false")
	}
}

// ---------- 单遍直写序列化器：与旧实现（Encoder + 每值 Encoder）行为对等 ----------

// legacyObj/legacyEncode 忠实复刻优化前的序列化路径，用于对等性对比与基准。
type legacyObj struct{ o *Obj }

func (l legacyObj) MarshalJSON() ([]byte, error) {
	var b bytes.Buffer
	b.WriteByte('{')
	for i, m := range l.o.Members {
		if i > 0 {
			b.WriteByte(',')
		}
		k, err := legacyMarshalNoEscape(m.Key)
		if err != nil {
			return nil, err
		}
		b.Write(k)
		b.WriteByte(':')
		v, err := legacyMarshalNoEscape(m.Value)
		if err != nil {
			return nil, err
		}
		b.Write(v)
	}
	b.WriteByte('}')
	return b.Bytes(), nil
}

func legacyMarshalNoEscape(v any) ([]byte, error) {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return bytes.TrimRight(b.Bytes(), "\n"), nil
}

// legacyWrap 把解析树深转换为旧路径类型，使嵌套对象也走每值 Encoder
func legacyWrap(v any) any {
	switch t := v.(type) {
	case *Obj:
		members := make([]Member, len(t.Members))
		for i, m := range t.Members {
			members[i] = Member{Key: m.Key, Value: legacyWrap(m.Value)}
		}
		return legacyObj{o: &Obj{Members: members}}
	case []any:
		out := make([]any, len(t))
		for i, e := range t {
			out[i] = legacyWrap(e)
		}
		return out
	default:
		return v
	}
}

func legacyEncode(v any, indent string) (string, error) {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetEscapeHTML(false)
	if indent != "" {
		enc.SetIndent("", indent)
	}
	if err := enc.Encode(legacyWrap(v)); err != nil {
		return "", err
	}
	return strings.TrimRight(b.String(), "\n"), nil
}

func TestEncodeParityWithLegacyPath(t *testing.T) {
	docs := []string{
		`{"a":"<html> & \"转义\"","b":"tab\there\nnl\u0001ctl","c":"line\u2028sep\u2029para"}`,
		`{"嵌套":{"数组":[1,-2.5e3,true,null,"中文🎉"],"空对象":{},"空数组":[]}}`,
		`"顶层字符串 <>&"`,
		`12345678901234567890`,
		`[]`,
		`{}`,
	}
	for _, in := range docs {
		val, err := parseOrdered([]byte(in))
		if err != nil {
			t.Fatalf("parse %s: %v", in, err)
		}
		for _, indent := range []string{"", "  "} {
			got, err := encodeValue(val, indent)
			if err != nil {
				t.Fatalf("new encode %s: %v", in, err)
			}
			want, err := legacyEncode(val, indent)
			if err != nil {
				t.Fatalf("legacy encode %s: %v", in, err)
			}
			if got != want {
				t.Fatalf("parity mismatch for %s (indent=%q):\n new=%q\nold=%q", in, indent, got, want)
			}
		}
	}
}

func TestEncodeInvalidUTF8Parity(t *testing.T) {
	// 非法 UTF-8 字节必须与旧路径一样替换为原始 U+FFFD 字符
	val := &Obj{Members: []Member{{Key: "k\x00\xff", Value: "bad\xff\xfe尾"}}}
	got, err := encodeValue(val, "")
	if err != nil {
		t.Fatal(err)
	}
	want, err := legacyEncode(val, "")
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("invalid-UTF8 parity mismatch:\n new=%q\nold=%q", got, want)
	}
}

// benchDoc 构造约 300KB 的混合结构文档
func benchDoc(b *testing.B) any {
	var sb strings.Builder
	sb.WriteString(`{"items":[`)
	for i := 0; i < 2000; i++ {
		if i > 0 {
			sb.WriteByte(',')
		}
		fmt.Fprintf(&sb,
			`{"id":%d,"name":"项目 <item> %d 中文🎉","tags":["a%d","b%d","c%d"],`+
				`"meta":{"active":%v,"score":%d.5,"note":null,"nested":{"deep":{"x":[1,2,3]}}}}`,
			i, i, i, i, i, i%2 == 0, i)
	}
	sb.WriteString(`],"total":2000}`)
	val, err := parseOrdered([]byte(sb.String()))
	if err != nil {
		b.Fatal(err)
	}
	return val
}

func BenchmarkEncodeValue(b *testing.B) {
	val := benchDoc(b)
	b.ReportAllocs()
	b.Run("new", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if _, err := encodeValue(val, "  "); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("legacy", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if _, err := legacyEncode(val, "  "); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func TestParseCache(t *testing.T) {
	app := NewApp()
	var sb strings.Builder
	sb.WriteString(`{"items":[`)
	for i := 0; i < 1500; i++ {
		if i > 0 {
			sb.WriteByte(',')
		}
		fmt.Fprintf(&sb, `{"id":%d,"name":"item-%d","tags":["a","b"]}`, i, i)
	}
	sb.WriteString("]}")
	big := sb.String()
	if len(big) < parseCacheMinSize {
		t.Fatalf("doc too small to exercise cache: %d bytes", len(big))
	}

	v1, err := app.parseCached(big)
	if err != nil {
		t.Fatalf("first parse: %v", err)
	}
	v2, err := app.parseCached(big)
	if err != nil {
		t.Fatalf("second parse: %v", err)
	}
	if v1 != v2 {
		t.Fatal("expected cache hit: same parse tree pointer for identical input")
	}

	// 小文档不缓存：每次都是新解析
	small := `{"a":1}`
	w1, err := app.parseCached(small)
	if err != nil {
		t.Fatal(err)
	}
	w2, err := app.parseCached(small)
	if err != nil {
		t.Fatal(err)
	}
	if w1 == w2 {
		t.Fatal("docs below parseCacheMinSize must not be cached")
	}

	// 解析失败应使缓存失效，之后相同内容仍全新解析
	if _, err := app.parseCached(big + " {}"); err == nil {
		t.Fatal("expected parse error for corrupted input")
	}
	v3, err := app.parseCached(big)
	if err != nil {
		t.Fatal(err)
	}
	if v3 == v1 {
		t.Fatal("cache must be invalidated after a failed parse")
	}

	// 上层 API 共享缓存：Analyze 之后 Format 不再重复解析，结果仍正确
	res := app.AnalyzeJSON(big)
	if !res.Valid {
		t.Fatalf("analyze failed: %s", res.Message)
	}
	out, err := app.FormatJSON(big)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, `"items"`) || !strings.Contains(out, `"item-1499"`) {
		t.Fatal("formatted output missing expected keys")
	}
}
