package mydoc

import (
	"bytes"
	"html/template"
	"reflect"
	"strings"

	"github.com/xxjwxc/public/mydoc/mymarkdown"
	"github.com/xxjwxc/public/mydoc/myswagger"
)

func (m *model) analysisStructInfo(info *StructInfo) {
	if info != nil {
		for i := 0; i < len(info.Items); i++ {
			tag := reflect.StructTag(strings.Trim(info.Items[i].Tag, "`"))

			// json
			tagStr := tag.Get("json")
			if tagStr == "-" || tagStr == "" {
				tagStr = tag.Get("url")
			}
			tagStrs := strings.Split(tagStr, ",")
			if len(tagStrs[0]) > 0 {
				info.Items[i].Name = tagStrs[0]
			}
			// -------- end

			// required
			tagStr = tag.Get("binding")
			tagStrs = strings.Split(tagStr, ",")
			for _, v := range tagStrs {
				if strings.EqualFold(v, "required") {
					info.Items[i].Requierd = true
					break
				}
			}
			// ---------------end

			// default
			info.Items[i].Default = tag.Get("default")
			// ---------------end

			if info.Items[i].TypeRef != nil {
				m.analysisStructInfo(info.Items[i].TypeRef)
			}
		}

	}
}

func (m *model) setDefinition(doc *myswagger.DocSwagger, tmp *StructInfo) string {
	if tmp != nil {
		var def myswagger.Definition
		def.Type = "object"
		def.Properties = make(map[string]myswagger.Propertie)
		for _, v2 := range tmp.Items {
			if v2.TypeRef != nil {
				def.Properties[v2.Name] = myswagger.Propertie{
					Ref: m.setDefinition(doc, v2.TypeRef),
				}
			} else {
				def.Properties[v2.Name] = myswagger.Propertie{
					Type:        myswagger.GetKvType(v2.Type, v2.IsArray, true),
					Format:      myswagger.GetKvType(v2.Type, v2.IsArray, false),
					Description: v2.Note,
				}
			}
		}
		doc.AddDefinitions(tmp.Name, def)
		return "#/definitions/" + tmp.Name
	}
	return ""
}

func buildRelativePath(prepath, routerPath string) string {
	if strings.HasSuffix(prepath, "/") {
		if strings.HasPrefix(routerPath, "/") {
			return prepath + strings.TrimPrefix(routerPath, "/")
		}
		return prepath + routerPath
	}

	if strings.HasPrefix(routerPath, "/") {
		return prepath + routerPath
	}

	return prepath + "/" + routerPath
}

func conType(pkg *StructInfo, tp string, isArray bool) string {
	re := tp
	if pkg != nil {
		re = "`" + pkg.Pkg + "." + tp + "`"
	}
	if isArray {
		re = "[]" + re
	}
	return re
}

func (m *model) buildSubStructMD(doc *mymarkdown.DocMarkdown, tmp *StructInfo) (jsonMp map[string]interface{}) {
	jsonMp = make(map[string]interface{})
	if tmp != nil {
		var info mymarkdown.TmpTable
		info.SC = "`"
		info.Pkg = tmp.Pkg
		info.Name = tmp.Name
		info.Note = tmp.Note
		for _, v := range tmp.Items {
			if v.TypeRef != nil {
				mp := m.buildSubStructMD(doc, v.TypeRef)
				jsonMp[v.Name] = doc.GetTypeList(mp, v.IsArray)
			} else {
				jsonMp[v.Name] = doc.GetValueType(v.Type, v.Default, v.IsArray)
			}

			info.Item = append(info.Item, mymarkdown.TmpElement{
				Name:     v.Name,
				Requierd: doc.GetBoolStr(v.Requierd),
				Type:     conType(v.TypeRef, v.Type, v.IsArray),
				Note:     v.Note,
			})
		}

		tmpl, err := template.New("struct_mod").
			Parse(doc.GetTableInfo())
		if err != nil {
			panic(err)
		}
		var buf bytes.Buffer
		tmpl.Execute(&buf, info)
		doc.AddBaseStruct(tmp.Pkg, tmp.Name, buf.String())
	}
	return
}

func (m *model) buildDefinitionMD(doc *mymarkdown.DocMarkdown, tmp *StructInfo) (rest string, jsonMp map[string]interface{}) {
	jsonMp = make(map[string]interface{})
	if tmp != nil {
		var info mymarkdown.TmpTable
		info.SC = "`"
		info.Name = tmp.Name
		info.Note = tmp.Note
		for _, v := range tmp.Items {
			if v.TypeRef != nil {
				mp := m.buildSubStructMD(doc, v.TypeRef)
				jsonMp[v.Name] = doc.GetTypeList(mp, v.IsArray)
			} else {
				jsonMp[v.Name] = doc.GetValueType(v.Type, v.Default, v.IsArray)
			}

			info.Item = append(info.Item, mymarkdown.TmpElement{
				Name:     v.Name,
				Requierd: doc.GetBoolStr(v.Requierd),
				Type:     conType(v.TypeRef, v.Type, v.IsArray),
				Note:     v.Note,
			})
		}

		tmpl, err := template.New("struct_mod").
			Parse(doc.GetTableInfo())
		if err != nil {
			panic(err)
		}
		var buf bytes.Buffer
		tmpl.Execute(&buf, info)
		rest = buf.String()
	}

	return
}
