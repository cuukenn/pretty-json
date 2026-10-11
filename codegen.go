package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"unicode"
)

// ---------- 类型推断 ----------

type jtype struct {
	kind     string // object | array | string | int | float | bool | null | any
	fields   []*jfield
	elem     *jtype
	nullable bool
	name     string // 对象类型生成的类型名
}

type jfield struct {
	key      string
	typ      *jtype
	optional bool
}

func inferType(v any) *jtype {
	switch t := v.(type) {
	case *Obj:
		jt := &jtype{kind: "object"}
		for _, m := range t.Members {
			jt.fields = append(jt.fields, &jfield{key: m.Key, typ: inferType(m.Value)})
		}
		return jt
	case []any:
		jt := &jtype{kind: "array"}
		if len(t) == 0 {
			jt.elem = &jtype{kind: "any"}
			return jt
		}
		var elems []*jtype
		for _, e := range t {
			elems = append(elems, inferType(e))
		}
		jt.elem = unionType(elems)
		return jt
	case json.Number:
		if _, ok := new(big.Int).SetString(t.String(), 10); ok {
			return &jtype{kind: "int"}
		}
		return &jtype{kind: "float"}
	case string:
		return &jtype{kind: "string"}
	case bool:
		return &jtype{kind: "bool"}
	case nil:
		return &jtype{kind: "null", nullable: true}
	default:
		return &jtype{kind: "any"}
	}
}

// unionType 合并多个推断类型（对象合并字段、数组合并元素、int+float→float，其余冲突→any）
func unionType(ts []*jtype) *jtype {
	if len(ts) == 0 {
		return &jtype{kind: "any"}
	}
	res := &jtype{}
	var objs []*jtype
	var arrs []*jtype
	var scalars []*jtype
	for _, t := range ts {
		if t == nil {
			continue
		}
		if t.nullable || t.kind == "null" {
			res.nullable = true
			if t.kind == "null" {
				continue
			}
		}
		switch t.kind {
		case "object":
			objs = append(objs, t)
		case "array":
			arrs = append(arrs, t)
		default:
			scalars = append(scalars, t)
		}
	}
	nk := 0
	for _, g := range [][]*jtype{objs, arrs, scalars} {
		if len(g) > 0 {
			nk++
		}
	}
	if nk != 1 {
		res.kind = "any"
		return res
	}
	switch {
	case len(objs) > 0:
		res.kind = "object"
		byKey := map[string]*jfield{}
		var order []string
		for _, o := range objs {
			for _, f := range o.fields {
				if ex, ok := byKey[f.key]; ok {
					ex.typ = unionType([]*jtype{ex.typ, f.typ})
				} else {
					nf := &jfield{key: f.key, typ: f.typ}
					byKey[f.key] = nf
					order = append(order, f.key)
				}
			}
		}
		for _, k := range order {
			f := byKey[k]
			// 任一来源对象缺少该键 → 可选
			for _, o := range objs {
				has := false
				for _, of := range o.fields {
					if of.key == k {
						has = true
						break
					}
				}
				if !has {
					f.optional = true
					break
				}
			}
			res.fields = append(res.fields, f)
		}
	case len(arrs) > 0:
		res.kind = "array"
		var elems []*jtype
		for _, a := range arrs {
			elems = append(elems, a.elem)
		}
		res.elem = unionType(elems)
	default:
		// 标量合并
		kind := scalars[0].kind
		for _, s := range scalars[1:] {
			if s.kind == kind {
				continue
			}
			if (kind == "int" && s.kind == "float") || (kind == "float" && s.kind == "int") {
				kind = "float"
				continue
			}
			kind = "any"
			break
		}
		res.kind = kind
	}
	return res
}

func typesEqual(a, b *jtype) bool {
	if a == nil || b == nil {
		return a == b
	}
	if a.kind != b.kind || a.nullable != b.nullable {
		return false
	}
	switch a.kind {
	case "object":
		if len(a.fields) != len(b.fields) {
			return false
		}
		for i := range a.fields {
			if a.fields[i].key != b.fields[i].key || a.fields[i].optional != b.fields[i].optional {
				return false
			}
			if !typesEqual(a.fields[i].typ, b.fields[i].typ) {
				return false
			}
		}
		return true
	case "array":
		return typesEqual(a.elem, b.elem)
	}
	return true
}

// ---------- 命名 ----------

type nameRegistry struct {
	used map[string]*jtype
}

// assign 为对象类型分配名称：同名且结构相同则复用，否则追加序号
func (r *nameRegistry) assign(base string, t *jtype) string {
	if base == "" {
		base = "Value"
	}
	cand := base
	for i := 2; ; i++ {
		if ex, ok := r.used[cand]; ok {
			if typesEqual(ex, t) {
				return cand
			}
			cand = base + strconv.Itoa(i)
			continue
		}
		r.used[cand] = t
		return cand
	}
}

// 常见缩写词在 PascalCase 中整体大写（Go 命名惯例）
var initialisms = map[string]string{
	"id": "ID", "ip": "IP", "uid": "UID", "uuid": "UUID", "url": "URL", "uri": "URI",
	"api": "API", "http": "HTTP", "https": "HTTPS", "json": "JSON", "xml": "XML",
	"html": "HTML", "css": "CSS", "sql": "SQL", "ui": "UI", "os": "OS", "db": "DB",
}

func pascal(s string) string {
	parts := splitIdent(s)
	var b strings.Builder
	for _, p := range parts {
		if up, ok := initialisms[strings.ToLower(p)]; ok {
			b.WriteString(up)
			continue
		}
		r := []rune(p)
		r[0] = unicode.ToUpper(r[0])
		b.WriteString(string(r))
	}
	out := b.String()
	if out == "" {
		return "F"
	}
	if unicode.IsDigit([]rune(out)[0]) {
		out = "F" + out
	}
	return out
}

// camel 首词全小写、后续词首字母大写（user_id → userId）
func camel(s string) string {
	parts := splitIdent(s)
	if len(parts) == 0 {
		return "f"
	}
	var b strings.Builder
	r := []rune(parts[0])
	r[0] = unicode.ToLower(r[0])
	b.WriteString(string(r))
	for _, p := range parts[1:] {
		r := []rune(p)
		r[0] = unicode.ToUpper(r[0])
		b.WriteString(string(r))
	}
	return b.String()
}

func splitIdent(s string) []string {
	var words []string
	var cur []rune
	for _, r := range s {
		if r == '_' || r == '-' || r == ' ' || r == '.' {
			if len(cur) > 0 {
				words = append(words, string(cur))
				cur = nil
			}
			continue
		}
		cur = append(cur, r)
	}
	if len(cur) > 0 {
		words = append(words, string(cur))
	}
	return words
}

// singularize 简易英文去复数（items→item, buses→bus, studies→study）
func singularize(s string) string {
	lower := strings.ToLower(s)
	switch {
	case strings.HasSuffix(lower, "ies") && len(s) > 4:
		return s[:len(s)-3] + "y"
	case strings.HasSuffix(lower, "sses") || strings.HasSuffix(lower, "shes") || strings.HasSuffix(lower, "xes"):
		return s[:len(s)-2]
	case strings.HasSuffix(lower, "s") && !strings.HasSuffix(lower, "ss") && len(s) > 3:
		return s[:len(s)-1]
	}
	return s
}

// collectStructs 深度优先为对象类型命名，返回生成顺序（根在前）
func collectStructs(t *jtype, hint string, reg *nameRegistry, out *[]*jtype, seen map[*jtype]bool) {
	if t == nil {
		return
	}
	switch t.kind {
	case "object":
		if seen[t] {
			return
		}
		if t.name == "" {
			t.name = reg.assign(hint, t)
		}
		seen[t] = true
		*out = append(*out, t)
		for _, f := range t.fields {
			collectStructs(f.typ, pascal(f.key), reg, out, seen)
		}
	case "array":
		eh := singularize(hint)
		if eh == hint { // 无法去复数时用 hint+Item，避免与别名同名冲突
			eh = hint + "Item"
		}
		collectStructs(t.elem, eh, reg, out, seen)
	}
}

// dedupeByName 结构相同的对象会复用同一类型名，按名称去重避免重复定义
func dedupeByName(structs []*jtype) []*jtype {
	seen := map[string]bool{}
	out := make([]*jtype, 0, len(structs))
	for _, s := range structs {
		if seen[s.name] {
			continue
		}
		seen[s.name] = true
		out = append(out, s)
	}
	return out
}

// ---------- 渲染 ----------

var goKeywords = map[string]bool{"break": true, "case": true, "chan": true, "const": true, "continue": true, "default": true, "defer": true, "else": true, "fallthrough": true, "for": true, "func": true, "go": true, "goto": true, "if": true, "import": true, "interface": true, "map": true, "package": true, "range": true, "return": true, "select": true, "struct": true, "switch": true, "type": true, "var": true}

var javaKeywords = map[string]bool{"abstract": true, "assert": true, "boolean": true, "break": true, "byte": true, "case": true, "catch": true, "char": true, "class": true, "const": true, "continue": true, "default": true, "do": true, "double": true, "else": true, "enum": true, "extends": true, "final": true, "finally": true, "float": true, "for": true, "goto": true, "if": true, "implements": true, "import": true, "instanceof": true, "int": true, "interface": true, "long": true, "native": true, "new": true, "package": true, "private": true, "protected": true, "public": true, "return": true, "short": true, "static": true, "strictfp": true, "super": true, "switch": true, "synchronized": true, "this": true, "throw": true, "throws": true, "transient": true, "try": true, "void": true, "volatile": true, "while": true}

func goType(t *jtype) string {
	if t == nil || t.kind == "any" || t.kind == "null" {
		return "interface{}"
	}
	if t.nullable && t.kind != "object" && t.kind != "array" {
		return "interface{}"
	}
	switch t.kind {
	case "object":
		if t.nullable {
			return "*" + t.name
		}
		return t.name
	case "array":
		return "[]" + goType(t.elem)
	case "string":
		return "string"
	case "int":
		return "int64"
	case "float":
		return "float64"
	case "bool":
		return "bool"
	}
	return "interface{}"
}

func tsType(t *jtype) string {
	if t == nil || t.kind == "any" {
		return "any"
	}
	var core string
	switch t.kind {
	case "null":
		core = "null"
	case "object":
		core = t.name
	case "array":
		et := tsType(t.elem)
		if strings.ContainsAny(et, "|&") {
			et = "(" + et + ")"
		}
		core = et + "[]"
	case "string":
		core = "string"
	case "int", "float":
		core = "number"
	case "bool":
		core = "boolean"
	}
	if t.nullable && t.kind != "null" {
		core += " | null"
	}
	return core
}

func javaType(t *jtype) string {
	if t == nil || t.kind == "any" || t.kind == "null" || t.nullable {
		return "Object"
	}
	switch t.kind {
	case "object":
		return t.name
	case "array":
		return "List<" + javaType(t.elem) + ">"
	case "string":
		return "String"
	case "int":
		return "long"
	case "float":
		return "double"
	case "bool":
		return "boolean"
	}
	return "Object"
}

func renderGo(structs []*jtype, root *jtype) string {
	var b strings.Builder
	if root.kind != "object" {
		fmt.Fprintf(&b, "type Root = %s\n", goType(root))
	}
	for _, s := range structs {
		fmt.Fprintf(&b, "type %s struct {\n", s.name)
		for _, f := range s.fields {
			fn := pascal(f.key)
			if goKeywords[strings.ToLower(fn)] {
				fn += "Field"
			}
			fmt.Fprintf(&b, "\t%s %s `json:%q`\n", fn, goType(f.typ), f.key)
		}
		fmt.Fprintf(&b, "}\n\n")
	}
	return strings.TrimRight(b.String(), "\n")
}

func renderTS(structs []*jtype, root *jtype) string {
	var b strings.Builder
	if root.kind != "object" {
		fmt.Fprintf(&b, "type Root = %s;\n", tsType(root))
	}
	for _, s := range structs {
		fmt.Fprintf(&b, "interface %s {\n", s.name)
		for _, f := range s.fields {
			opt := ""
			if f.optional {
				opt = "?"
			}
			fmt.Fprintf(&b, "\t%s%s: %s;\n", camel(f.key), opt, tsType(f.typ))
		}
		fmt.Fprintf(&b, "}\n\n")
	}
	return strings.TrimRight(b.String(), "\n")
}

func renderJava(structs []*jtype) string {
	// 嵌套类：第一个为顶层类，其余作为其静态内部类
	var b strings.Builder
	b.WriteString("import java.util.List;\n\n")
	root := structs[0]
	inner := structs[1:]
	fmt.Fprintf(&b, "public class %s {\n", root.name)
	writeJavaFields(&b, root, "    ")
	if len(inner) > 0 {
		b.WriteString("\n")
		for _, s := range inner {
			fmt.Fprintf(&b, "    public static class %s {\n", s.name)
			writeJavaFields(&b, s, "        ")
			b.WriteString("    }\n\n")
		}
	}
	b.WriteString("}")
	return strings.TrimRight(b.String(), "\n")
}

func writeJavaFields(b *strings.Builder, s *jtype, indent string) {
	for _, f := range s.fields {
		fn := camel(f.key)
		if javaKeywords[fn] {
			fn = "_" + fn
		}
		fmt.Fprintf(b, "%spublic %s %s;\n", indent, javaType(f.typ), fn)
	}
}

// GenerateCode 根据 JSON 样例生成类型定义代码；lang 支持 go / typescript(ts) / java
func (a *App) GenerateCode(jsonStr, lang string) (string, error) {
	if strings.TrimSpace(jsonStr) == "" {
		return "", fmt.Errorf("请输入JSON内容")
	}
	val, err := parseOrdered([]byte(jsonStr))
	if err != nil {
		return "", fmt.Errorf("JSON格式错误: %v", err)
	}
	lang = strings.ToLower(strings.TrimSpace(lang))
	root := inferType(val)
	reg := &nameRegistry{used: map[string]*jtype{}}
	var structs []*jtype
	collectStructs(root, "Root", reg, &structs, map[*jtype]bool{})
	structs = dedupeByName(structs)

	switch lang {
	case "go":
		return renderGo(structs, root), nil
	case "typescript", "ts":
		return renderTS(structs, root), nil
	case "java":
		if root.kind != "object" {
			return "", fmt.Errorf("Java 代码生成要求顶层为 JSON 对象")
		}
		return renderJava(structs), nil
	default:
		return "", fmt.Errorf("不支持的语言: %s（支持 go / typescript / java）", lang)
	}
}
