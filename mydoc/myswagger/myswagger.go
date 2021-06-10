package myswagger

import (
	"strings"

	"github.com/xxjwxc/public/tools"
)

// DocSwagger ...
type DocSwagger struct {
	client *APIBody
}

// NewDoc 新建一个swagger doc
func NewDoc() *DocSwagger {
	doc := &DocSwagger{}
	doc.client = &APIBody{
		Head:     Head{Swagger: version},
		Info:     info,
		Host:     host,
		BasePath: basePath,
		// Tags
		Schemes: schemes,
		// Patchs
		// SecurityDefinitions
		// Definitions
		ExternalDocs: externalDocs,
	}
	doc.client.Patchs = make(map[string]map[string]Param)
	return doc
}

// AddTag add tag (排他)
func (doc *DocSwagger) AddTag(tag Tag) {
	for _, v := range doc.client.Tags {
		if v.Name == tag.Name { // find it
			return
		}
	}

	doc.client.Tags = append(doc.client.Tags, tag)
}

// AddDefinitions 添加 通用结构体定义
func (doc *DocSwagger) AddDefinitions(key string, def Definition) {
	// for k := range doc.client.Definitions {
	// 	if k == key { // find it
	// 		return
	// 	}
	// }
	if doc.client.Definitions == nil {
		doc.client.Definitions = make(map[string]Definition)
	}

	doc.client.Definitions[key] = def
}

// AddPatch ... API 路径 paths 和操作在 API 规范的全局部分定义
func (doc *DocSwagger) AddPatch(url string, p Param, metheds ...string) {
	if !strings.HasPrefix(url, "/") {
		url = "/" + url
	}

	if doc.client.Patchs[url] == nil {
		doc.client.Patchs[url] = make(map[string]Param)
	}
	if len(p.Consumes) == 0 {
		p.Consumes = reqCtxType
	}
	if len(p.Produces) == 0 {
		p.Produces = respCtxType
	}
	if p.Responses == nil {
		p.Responses = map[string]Resp{
			"400": {Description: "v"},
			"404": {Description: "not found"},
			"405": {Description: "Validation exception"},
		}
	}

	for _, v := range metheds {
		doc.client.Patchs[url][strings.ToLower(v)] = p
	}
}

// GetAPIString 获取返回数据
func (doc *DocSwagger) GetAPIString() string {
	return tools.GetJSONStr(doc.client, true)
}

var kvType = map[string]string{ // array, boolean, integer, number, object, string
	"int":     "integer",
	"uint":    "integer",
	"byte":    "integer",
	"rune":    "integer",
	"int8":    "integer",
	"int16":   "integer",
	"int32":   "integer",
	"int64":   "integer",
	"uint8":   "integer",
	"uint16":  "integer",
	"uint32":  "integer",
	"uint64":  "integer",
	"uintptr": "integer",
	"float32": "integer",
	"float64": "integer",
	"bool":    "boolean",
	"map":     "object",
	"Time":    "string"}

var kvFormat = map[string]string{}

// GetKvType 获取类型转换
func GetKvType(k string, isArray, isType bool) string {
	if isArray {
		if isType {
			return "object"
		}
		return "array"
	}

	if isType {
		if _, ok := kvType[k]; ok {
			return kvType[k]
		}
		return k
	}
	if _, ok := kvFormat[k]; ok {
		return kvFormat[k]
	}
	return k
}
