package myswagger

// Head Swagger 版本
type Head struct {
	Swagger string `json:"swagger"`
}

// Info 指定 API 的 info-title
type Info struct {
	Description string `json:"description"`
	Version     string `json:"version"`
	Title       string `json:"title"`
}

// ExternalDocs tags of group
type ExternalDocs struct {
	Description string `json:"description,omitempty"` // 描述
	URL         string `json:"url,omitempty"`         // url addr
}

// Tag group of tags
type Tag struct {
	Name         string        `json:"name"`                   // tags name
	Description  string        `json:"description"`            // 描述
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty"` // doc group of tags
}

// Schema 引用
type Schema struct {
	Ref string `json:"$ref"` // 主体模式和响应主体模式中引用
}

// Element 元素定义
type Element struct {
	In          string `json:"in"`             // 入参
	Name        string `json:"name"`           // 参数名字
	Description string `json:"description"`    // 描述
	Required    bool   `json:"required"`       // 是否必须
	Type        string `json:"type,omitempty"` // 类型
	Schema      Schema `json:"schema"`         // 引用
}

// Param API 路径 paths 和操作在 API 规范的全局部分定义
type Param struct {
	Tags        []string                     `json:"tags"`                  // 分组标记
	Summary     string                       `json:"summary"`               // 摘要
	Description string                       `json:"description"`           // 描述
	OperationID string                       `json:"operationId,omitempty"` // 操作id
	Consumes    []string                     `json:"consumes"`              // Parameter content type
	Produces    []string                     `json:"produces"`              // Response content type
	Parameters  []Element                    `json:"parameters"`            // 请求参数
	Responses   map[string]Resp 			 `json:"responses"`             // 返回参数
	Security    interface{}                 `json:"security,omitempty"`    // 认证信息
}
// 为返回的类型的文档加入描述
type Resp struct {
	Description string                 `json:"description"`
	Schema      map[string]interface{} `json:"schema,omitempty"`
}

// SecurityDefinitions 安全验证
type SecurityDefinitions struct {
	PetstoreAuth interface{} `json:"petstore_auth,omitempty"` // 安全验证定义
	AIPKey       interface{} `json:"api_key,omitempty"`       // api key
}

// Propertie 属性
type Propertie struct {
	Type        string      `json:"type,omitempty"` // 类型
	Format      string      `json:"format"`         // format 类型
	Description string      `json:"description"`    // 描述
	Enum        interface{} `json:"enum,omitempty"` // enum
	Ref         string      `json:"$ref,omitempty"` // 主体模式和响应主体模式中引用
}

// XML xml
type XML struct {
	Name    string `json:"name"`
	Wrapped bool   `json:"wrapped"`
}

// Definition 通用结构体定义
type Definition struct {
	Type       string               `json:"type"`       // 类型 object
	Properties map[string]Propertie `json:"properties"` // 属性列表
	XML        XML                  `json:"xml"`
}

// APIBody swagger api body info
type APIBody struct {
	Head
	Info                Info                        `json:"info"`
	Host                string                      `json:"host"`     // http host
	BasePath            string                      `json:"basePath"` // 根级别
	Tags                []Tag                       `json:"tags"`
	Schemes             []string                    `json:"schemes"`                       // http/https
	Patchs              map[string]map[string]Param `json:"paths"`                         // API 路径
	SecurityDefinitions SecurityDefinitions         `json:"securityDefinitions,omitempty"` // 安全验证
	Definitions         map[string]Definition       `json:"definitions"`                   // 通用结构体定义
	ExternalDocs        ExternalDocs                `json:"externalDocs"`                  // 外部链接
}
