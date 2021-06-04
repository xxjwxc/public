package myast

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"runtime"
	"strings"

	"github.com/xxjwxc/public/errors"
	"github.com/xxjwxc/public/tools"
)

var importFile map[string]string // 自定义包文件

func init() {
	importFile = make(map[string]string)
}

// AddImportFile 添加自定义import文件列表
func AddImportFile(k, v string) {
	importFile[k] = v
}

// GetModuleInfo find and get module info , return module [ name ,path ]
// 通过model信息获取[model name] [和 根目录绝对地址]
func GetModuleInfo(n int) (string, string, bool) {
	index := n
	// This is used to support third-party package encapsulation
	// 这样做用于支持第三方包封装,(主要找到main调用者)
	for { // find main file
		_, filename, _, ok := runtime.Caller(index)
		if ok {
			if strings.HasSuffix(filename, "runtime/asm_amd64.s") {
				index = index - 2
				break
			}
			if strings.HasSuffix(filename, "runtime/asm_arm64.s") {
				index = index - 2
				break
			}
			index++
		} else {
			panic(errors.New("package parsing failed:can not find main files"))
		}
	}

	_, filename, _, _ := runtime.Caller(index)
	filename = strings.Replace(filename, "\\", "/", -1) // offset
	for {
		n := strings.LastIndex(filename, "/")
		if n > 0 {
			filename = filename[0:n]
			if tools.CheckFileIsExist(filename + "/go.mod") {
				list := tools.ReadFile(filename + "/go.mod")
				if len(list) > 0 {
					line := strings.TrimSpace(list[0])
					if len(line) > 0 && strings.HasPrefix(line, "module") { // find it
						return strings.TrimSpace(strings.TrimPrefix(line, "module")), filename, true
					}
				}
			}
		} else {
			break
			// panic(errors.New("package parsing failed:can not find module file[go.mod] , golang version must up 1.11"))
		}
	}

	// never reach
	return "", "", false
}

// EvalSymlinks  Return to relative path . 通过module 游标返回包相对路径
func EvalSymlinks(modPkg, modFile, objPkg string) string {
	if strings.EqualFold(objPkg, "main") { // if main return default path
		return modFile
	}

	if strings.HasPrefix(objPkg, modPkg) {
		return modFile + strings.Replace(objPkg[len(modPkg):], ".", "/", -1)
	}

	// 自定义文件中查找
	tmp := importFile[objPkg]
	if len(tmp) > 0 {
		return tmp
	}

	// get the error space
	panic(errors.Errorf("can not eval pkg:[%v] must include [%v]", objPkg, modPkg))
}

// Re
// GetAstPkgs Parsing source file ast structure (with main restriction).解析源文件ast结构(带 main 限制)
func GetAstPkgs(objPkg, objFile string) (*ast.Package, bool) {
	fileSet := token.NewFileSet()
	astPkgs, err := parser.ParseDir(fileSet, objFile, func(info os.FileInfo) bool {
		name := info.Name()
		return !info.IsDir() && !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".go")
	}, parser.ParseComments)
	if err != nil {
		return nil, false
	}

	// check the package is same.判断 package 是否一致
	for _, pkg := range astPkgs {
		if objPkg == pkg.Name || strings.HasSuffix(objPkg, "/"+pkg.Name) { // find it
			return pkg, true
		}
	}

	// not find . maybe is main pakge and find main package
	if objPkg == "main" {
		dirs := tools.GetPathDirs(objFile) // get all of dir
		for _, dir := range dirs {
			if !strings.HasPrefix(dir, ".") {
				pkg, b := GetAstPkgs(objPkg, objFile+"/"+dir)
				if b {
					return pkg, true
				}
			}
		}
	}

	// ast.Print(fileSet, astPkgs)

	return nil, false
}

// GetObjFunMp find all exported func of sturct objName
// GetObjFunMp 类中的所有导出函数
func GetObjFunMp(astPkg *ast.Package, objName string) map[string]*ast.FuncDecl {
	funMp := make(map[string]*ast.FuncDecl)
	// find all exported func of sturct objName
	for _, fl := range astPkg.Files {
		for _, d := range fl.Decls {
			switch specDecl := d.(type) {
			case *ast.FuncDecl:
				if specDecl.Recv != nil {
					if exp, ok := specDecl.Recv.List[0].Type.(*ast.StarExpr); ok { // Check that the type is correct first beforing throwing to parser
						if strings.Compare(fmt.Sprint(exp.X), objName) == 0 { // is the same struct
							funMp[specDecl.Name.String()] = specDecl // catch
						}
					}
				}
			}
		}
	}

	return funMp
}

// AnalysisImport 分析整合import相关信息
func AnalysisImport(astPkgs *ast.Package) map[string]string {
	imports := make(map[string]string)
	for _, f := range astPkgs.Files {
		for _, p := range f.Imports {
			k := ""
			if p.Name != nil {
				k = p.Name.Name
			}
			v := strings.Trim(p.Path.Value, `"`)
			if len(k) == 0 {
				n := strings.LastIndex(v, "/")
				if n > 0 {
					k = v[n+1:]
				} else {
					k = v
				}
			}
			imports[k] = v
		}
	}

	return imports
}

// GetImportPkg 分析得出 pkg
func GetImportPkg(i string) string {
	n := strings.LastIndex(i, "/")
	if n > 0 {
		return i[n+1:]
	}
	return i
}
