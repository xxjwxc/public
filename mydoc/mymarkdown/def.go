package mymarkdown

import "html/template"

var headMod = []string{
	`
	#### Brief description:
	
	- [%s]
	- [Waiting to write...]
	`,
	`
	#### 简要描述：
	
	- [%s]
	- [等待写入]
	`,
}

var bTrue = []string{"`YES`", "`是`"}
var bFalse = []string{"NO", "否"}

var tableMod = []string{
	`{{$obj := .}}
- {{$obj.SC}} {{$obj.Name}} {{$obj.SC}} : {{$obj.Note}}

|Parameter| Requierd | Type | Description|
|:----    |:---|:----- |-----   |{{range $oem := $obj.Item}}
|{{$obj.SC}}{{$oem.Name}}{{$obj.SC}} | {{$oem.Requierd}}|{{$oem.Type}}|{{$oem.Note}}   |{{end}}
	`,
	`{{$obj := .}}
- {{$obj.SC}} {{$obj.Name}} {{$obj.SC}} : {{$obj.Note}}

|参数名|是否必须|类型|说明|
|:----    |:---|:----- |-----   |{{range $oem := $obj.Item}}
|{{$obj.SC}}{{$oem.Name}}{{$obj.SC}} | {{$oem.Requierd}}|{{$oem.Type}}|{{$oem.Note}}   |{{end}}
`,
}

var bodyMod = []string{
	`
{{$obj := .}}
## [Viewing tools](https://www.iminho.me/)

## Overview:
- [{{$obj.Class}}]
- [{{$obj.Note}}]
{{range $oem := $obj.Item}}
--------------------

### {{$oem.InterfaceName}}

#### Brief description:

- [{{$oem.Note}}]

#### Request URL:

- {{$oem.RouterPath}}

#### Methods:
{{range $me := $oem.Methods}}
- {{$me}}{{end}}

#### Parameters:
{{$oem.ReqTab}}

#### Request example:
{{$obj.SC}}{{$obj.SC}}{{$obj.SC}}
{{$oem.ReqJSON}}
{{$obj.SC}}{{$obj.SC}}{{$obj.SC}}

#### Return parameter description:
{{$oem.RespTab}}

#### Return example:
	
{{$obj.SC}}{{$obj.SC}}{{$obj.SC}}
{{$oem.RespJSON}}
{{$obj.SC}}{{$obj.SC}}{{$obj.SC}}

#### Remarks:

- {{$oem.Note}}
{{end}}	

--------------------
--------------------

#### Custom type:
{{range $k,$v := $obj.BMP}}
#### {{$obj.SC}} {{$k}} {{$obj.SC}}
{{range $v1 := $v}}
{{$v1.Context}}
{{end}}
{{end}}	
`,
	`
{{$obj := .}}
## [推荐查看工具](https://www.iminho.me/)

## 总览:
- [{{$obj.Class}}]
- [{{$obj.Note}}]
{{range $oem := $obj.Item}}
--------------------

### {{$oem.InterfaceName}}

#### 简要描述：

- [{{$oem.Note}}]

#### 请求URL:

- {{$oem.RouterPath}}

#### 请求方式：
{{range $me := $oem.Methods}}
- {{$me}}{{end}}

#### 请求参数:
{{$oem.ReqTab}}

#### 请求示例:
{{$obj.SC}}{{$obj.SC}}{{$obj.SC}}
{{$oem.ReqJSON}}
{{$obj.SC}}{{$obj.SC}}{{$obj.SC}}

#### 返回参数说明:
{{$oem.RespTab}}

#### 返回示例:
	
{{$obj.SC}}{{$obj.SC}}{{$obj.SC}}
{{$oem.RespJSON}}
{{$obj.SC}}{{$obj.SC}}{{$obj.SC}}

#### 备注:

- {{$oem.Note}}
{{end}}	

--------------------
--------------------

#### 自定义类型:
{{range $k,$v := $obj.BMP}}
#### {{$obj.SC}} {{$k}} {{$obj.SC}}
{{range $v1 := $v}}
{{$v1.Context}}
{{end}}
{{end}}
`,
}

// TmpElement 元素
type TmpElement struct {
	Name     string
	Requierd string
	Type     string
	Note     string
}

// TmpTable 模板
type TmpTable struct {
	SC   string
	Pkg  string
	Name string
	Note string
	Item []TmpElement
}

// TmpSub 模板
type TmpSub struct {
	ReqTab  string        // 请求参数列表
	ReqJSON template.HTML // 请求示例

	RespTab  string        // 返回参数列表
	RespJSON template.HTML // 返回示例

	Methods []string // 请求方式

	Note          string // 注释
	RouterPath    string // 请求url
	InterfaceName string // 接口名
}

// TmpInterface 模板
type TmpInterface struct {
	SC    string
	Class string
	Note  string
	Item  []TmpSub
	BMP   map[string][]BaseStruct
}

// BaseStruct 模板
type BaseStruct struct { // 基础类型
	Class   string
	Context string
}
