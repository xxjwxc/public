package myast

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
	"strings"

	"github.com/xxjwxc/public/mydoc"
	"github.com/xxjwxc/public/mylog"
	"github.com/xxjwxc/public/tools"
)

type structAnalys struct {
	ModPkg, ModFile string
}

// NewStructAnalys 新建一个导出结构体类
func NewStructAnalys(modPkg, modFile string) *structAnalys {
	result := &structAnalys{ModPkg: modPkg, ModFile: modFile}
	return result
}

// ParserStruct 解析结构体定义及相关信息
func (a *structAnalys) ParserStruct(astPkg *ast.Package, structName string) (info *mydoc.StructInfo) {
	if astPkg == nil {
		return nil
	}

	if tools.IsInternalType(structName) { // 内部类型
		return &mydoc.StructInfo{
			Name: structName,
		}
	}
	// ast.Print(token.NewFileSet(), astPkg)

	for _, fl := range astPkg.Files {
		for _, d := range fl.Decls {
			switch specDecl := d.(type) {
			case *ast.GenDecl:
				for _, subitem := range specDecl.Specs {
					switch specDecl.Tok {
					case token.TYPE:
						spec := subitem.(*ast.TypeSpec)
						switch st := spec.Type.(type) {
						case *ast.StructType:
							if spec.Name.Name == structName { // find it
								info = new(mydoc.StructInfo)
								info.Pkg = astPkg.Name
								if specDecl.Doc != nil { // 如果有注释
									for _, v := range specDecl.Doc.List { // 结构体注释
										t := strings.TrimSpace(strings.TrimPrefix(v.Text, "//"))
										if strings.HasPrefix(t, structName) { // find note
											t = strings.TrimSpace(strings.TrimPrefix(t, structName))
											info.Note += t
										}
									}
								}

								info.Name = structName
								info.Items = a.structFieldInfo(astPkg, st)
								return
							}
						}
					}
				}
			}
		}
	}
	return nil
}

func (a *structAnalys) structFieldInfo(astPkg *ast.Package, sinfo *ast.StructType) (items []mydoc.ElementInfo) {
	if sinfo == nil || sinfo.Fields == nil {
		return
	}

	importMP := AnalysisImport(astPkg)

	var info mydoc.ElementInfo
	for _, field := range sinfo.Fields.List {
		info = mydoc.ElementInfo{}
		for _, fnames := range field.Names {
			info.Name += fnames.Name
		}
		// 判断是否是导出属性(导出属性才允许)(首字母大写)
		strArry := []rune(info.Name)
		if len(strArry) > 0 && (strArry[0] >= 97 && strArry[0] <= 122) { // 首字母小写
			continue
		}

		if field.Tag != nil {
			info.Tag = strings.Trim(field.Tag.Value, "`")
			// 可以忽略某个字段, 主要是在扫描response返回类型的时候
			// 有时候不一定会为每个请求方法都定义一个返回的struct(大概是因为懒)
			tag := reflect.StructTag(info.Tag)
			tagStr := tag.Get("json")
			if tagStr == "-" { // 忽略的json字段
				continue
			}
		}
		if field.Comment != nil {
			info.Note = strings.TrimSpace(field.Comment.Text())
		}
		if field.Doc != nil {
			info.Note += strings.TrimSpace(field.Doc.List[0].Text)
		}

		switch exp := field.Type.(type) {
		case *ast.SelectorExpr: // 非本文件包
			a.dealSelectorExpr(exp, &info, importMP)
		case *ast.ArrayType:
			info.IsArray = true
			switch x := exp.Elt.(type) {
			case *ast.SelectorExpr: // 非本文件包
				a.dealSelectorExpr(x, &info, importMP)
			case *ast.Ident:
				a.dealIdent(astPkg, x, &info)
			case *ast.StarExpr:
				switch x1 := x.X.(type) {
				case *ast.SelectorExpr: // 非本文件包
					a.dealSelectorExpr(x1, &info, importMP)
				case *ast.Ident:
					a.dealIdent(astPkg, x1, &info)
				}
			case *ast.InterfaceType:
				info.Type = "Interface"
			}
		case *ast.StarExpr:
			switch x := exp.X.(type) {
			case *ast.SelectorExpr: // 非本文件包
				a.dealSelectorExpr(x, &info, importMP)
			case *ast.Ident:
				a.dealIdent(astPkg, x, &info)
			}
		case *ast.Ident: // 本文件
			a.dealIdent(astPkg, exp, &info)
		case *ast.MapType: // map
			key := ""
			value := ""
			switch x := exp.Key.(type) {
			case *ast.Ident:
				key = x.Name
			case *ast.StarExpr:
				switch x1 := x.X.(type) {
				case *ast.SelectorExpr: // 非本文件包
					key = x1.Sel.Name
				case *ast.Ident:
					key = x1.Name
				}
			case *ast.SelectorExpr: // 非本文件包
				key = x.Sel.Name
			}
			switch x := exp.Value.(type) {
			case *ast.Ident:
				value = x.Name
			case *ast.StarExpr:
				switch x1 := x.X.(type) {
				case *ast.SelectorExpr: // 非本文件包
					value = x1.Sel.Name
				case *ast.Ident:
					value = x1.Name
				}
			case *ast.SelectorExpr: // 非本文件包
				value = x.Sel.Name
			}
			info.Type = fmt.Sprintf("map (%v,%v)", key, value)
		}

		if len(info.Type) == 0 {
			mylog.Errorf("can not deal the type : %v", field.Type)
		}

		items = append(items, info)

	}
	return items
}

func (a *structAnalys) dealSelectorExpr(exp *ast.SelectorExpr, info *mydoc.ElementInfo, importMP map[string]string) { // 非本文件包
	info.Type = exp.Sel.Name
	if !tools.IsInternalType(info.Type) { // 非基础类型(time)
		if x, ok := exp.X.(*ast.Ident); ok {
			if v, ok := importMP[x.Name]; ok {
				objFile := EvalSymlinks(a.ModPkg, a.ModFile, v)
				objPkg := GetImportPkg(v)
				astFile, _b := GetAstPkgs(objPkg, objFile)
				if _b {
					info.TypeRef = a.ParserStruct(astFile, info.Type)
				}
			}
		}
	}
}

func (a *structAnalys) dealIdent(astPkg *ast.Package, exp *ast.Ident, info *mydoc.ElementInfo) { // 本文件
	info.Type = exp.Name
	if !tools.IsInternalType(info.Type) { // 非基础类型
		info.TypeRef = a.ParserStruct(astPkg, info.Type)
	}
}
