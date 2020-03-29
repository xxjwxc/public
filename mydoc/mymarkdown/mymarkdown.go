package mymarkdown

import (
	"bytes"
	"html/template"
	"strconv"
	"strings"
	"time"

	"github.com/xxjwxc/public/tools"
)

// DocMarkdown
type DocMarkdown struct {
	n            int
	baseStructMP map[string][]BaseStruct
}

// NewDoc 新建一个markdown文件
func NewDoc() *DocMarkdown {
	doc := &DocMarkdown{}
	if tools.GetLocalSystemLang(true) == "zh" { // en
		doc.n = 1
	}
	doc.baseStructMP = make(map[string][]BaseStruct)
	return doc
}

// AddBaseStruct 添加一个基础类型
func (m *DocMarkdown) AddBaseStruct(pkg, class, context string) {
	for _, v := range m.baseStructMP[pkg] {
		if v.Class == class {
			return
		}
	}
	m.baseStructMP[pkg] = append(m.baseStructMP[pkg], BaseStruct{
		Class:   class,
		Context: context,
	})
}

// GetBoolStr 获取bool类型字符串
func (m *DocMarkdown) GetBoolStr(b bool) string {
	if b {
		return bTrue[m.n]
	}

	return bFalse[m.n]
}

// GetTableInfo 获取table表格
func (m *DocMarkdown) GetTableInfo() string {
	return tableMod[m.n]
}

// GetBodyInfo 获取table表格
func (m *DocMarkdown) GetBodyInfo() string {
	return bodyMod[m.n]
}

// GetTypeList 转换typelist
func (m *DocMarkdown) GetTypeList(k interface{}, isArray bool) interface{} {
	if isArray {
		return []interface{}{k}
	}
	return k
}

// GetValueType 根据类型获取内容
func (m *DocMarkdown) GetValueType(k, v string, isArray bool) interface{} {
	array := strings.Split(v, ",")
	k = strings.ToLower(k)
	switch k {
	case "string":
		if isArray {
			var list []string
			for _, v := range array {
				list = append(list, v)
			}
			return list
		}
		return v
	case "bool":
		if isArray {
			var list []bool
			for _, v := range array {
				list = append(list, (v == "true" || v == "1"))
			}
			return list
		}

		return (v == "true" || v == "1")
	case "int", "uint", "byte", "rune", "int8", "int16", "int32", "int64", "uint8", "uint16", "uint32", "uint64", "uintptr":
		if isArray {
			var list []int
			for _, v := range array {
				i, _ := strconv.Atoi(v)
				list = append(list, i)
			}
			return list
		}

		i, _ := strconv.Atoi(v)
		return i
	case "float32", "float64":
		if isArray {
			var list []float64
			for _, v := range array {
				f, _ := strconv.ParseFloat(v, 64)
				list = append(list, f)
			}
			return list
		}

		f, _ := strconv.ParseFloat(v, 64)
		return f
	case "map":
		return v
	case "Time":
		if isArray {
			var list []tools.Time
			for _, v := range array {
				var t tools.Time
				t.Time = time.Unix(tools.StringTimetoUnix(v), 0)
				list = append(list, t)
			}
			return list
		}

		var t tools.Time
		t.Time = time.Unix(tools.StringTimetoUnix(v), 0)
		return t
	}

	return v
}

// GenMarkdown 生成markdown
func (m *DocMarkdown) GenMarkdown(info TmpInterface) string {
	info.SC = "`"
	info.BMP = m.baseStructMP
	tmpl, err := template.New("struct_markdown").Parse(m.GetBodyInfo())
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	tmpl.Execute(&buf, info)
	return buf.String()
}
