package mydoc

import (
	"html/template"
	"sort"

	"github.com/xxjwxc/public/mylog"

	"github.com/xxjwxc/public/mydoc/mymarkdown"
	"github.com/xxjwxc/public/mydoc/myswagger"
	"github.com/xxjwxc/public/tools"
)

type model struct {
	Group string // group 标记
	MP    map[string]map[string]DocModel
}

// NewDoc 新建一个doc模板
func NewDoc(group string) *model {
	doc := &model{Group: group}
	doc.MP = make(map[string]map[string]DocModel)
	return doc
}

// 添加一个
func (m *model) AddOne(group, methodName, routerPath string, methods []string, note string, req, resp *StructInfo) {
	if m.MP[group] == nil {
		m.MP[group] = make(map[string]DocModel)
	}

	m.analysisStructInfo(req)
	m.analysisStructInfo(resp)
	m.MP[group][routerPath] = DocModel{
		RouterPath: routerPath,
		Methods:    methods,
		Note:       note,
		Req:        req,
		Resp:       resp,
		MethodName: methodName,
	}
}

// GenSwagger 生成swagger文档
func (m *model) GenSwagger(outPath string) {
	doc := myswagger.NewDoc()
	reqRef, _ := "", ""

	var sortStr []string
	// define
	for k, v := range m.MP {
		for _, v1 := range v {
			reqRef = m.setDefinition(doc, v1.Req)
			//respRef = m.setDefinition(doc, v1.Resp)
		}
		sortStr = append(sortStr, k)
	}
	// ------------------end
	sort.Strings(sortStr)

	for _, k := range sortStr {
		v := m.MP[k]
		tag := myswagger.Tag{Name: k}
		doc.AddTag(tag)
		for _, v1 := range v {
			var p myswagger.Param
			p.Tags = []string{k}
			p.Summary = v1.Note
			p.Description = v1.Note
			// p.OperationID = "addPet"
			p.Parameters = []myswagger.Element{myswagger.Element{
				In:          "body", //  body, header, formData, query, path
				Name:        "body", //  body, header, formData, query, path
				Description: v1.Note,
				Required:    true,
				Schema: myswagger.Schema{
					Ref: reqRef,
				},
			}}
			doc.AddPatch(buildRelativePath(m.Group, v1.RouterPath), p, v1.Methods...)
		}
	}

	jsonsrc := doc.GetAPIString()
	mylog.Infof("output swagger doc: %v", outPath+"swagger.json")
	tools.WriteFile(outPath+"swagger.json", []string{jsonsrc}, true)
}

// GenMd 生成 markdown 文档
func (m *model) GenMarkdown(outPath string) {
	var sortStr []string
	// define
	for k := range m.MP {
		sortStr = append(sortStr, k)
	}
	sort.Strings(sortStr)

	for _, k := range sortStr {
		v := m.MP[k]
		doc := mymarkdown.NewDoc()
		var tmp mymarkdown.TmpInterface
		tmp.Class = k
		tmp.Note = "Waiting to write..."
		for _, v1 := range v {
			reqTable, reqMp := m.buildDefinitionMD(doc, v1.Req)
			resTable, respMpon := m.buildDefinitionMD(doc, v1.Resp)
			var sub mymarkdown.TmpSub
			sub.ReqTab = reqTable
			sub.ReqJSON = template.HTML(tools.GetJSONStr(reqMp, true))

			sub.RespTab = resTable
			sub.RespJSON = template.HTML(tools.GetJSONStr(respMpon, true))

			sub.Methods = v1.Methods
			sub.Note = v1.Note
			sub.InterfaceName = v1.MethodName
			sub.RouterPath = buildRelativePath(myswagger.GetHost(), buildRelativePath(m.Group, v1.RouterPath))
			tmp.Item = append(tmp.Item, sub)
		}
		jsonsrc := doc.GenMarkdown(tmp)
		mylog.Infof("output markdown doc: %v", outPath+k+".md")
		tools.WriteFile(outPath+k+".md", []string{jsonsrc}, true)
	}
}
