package main

import (
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
