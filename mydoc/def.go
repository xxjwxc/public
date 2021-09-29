package mydoc

// ElementInfo 结构信息
type ElementInfo struct {
	Name string // 参数名
	// URL      string      // web 访问参数
	Tag      string      // 标签
	Type     string      // 类型
	TypeRef  *StructInfo // 类型定义
	IsArray  bool        // 是否是数组
	Requierd bool        // 是否必须
	Note     string      // 注释
	Default  string      // 默认值
}

// StructInfo struct define
type StructInfo struct {
	Items []ElementInfo // 结构体元素
	Note  string        // 注释
	Name  string        //结构体名字
	Pkg   string        // 包名
}

// DocModel model
type DocModel struct {
	RouterPath string
	Methods    []string
	Note       string
	MethodName string
	Req, Resp  *StructInfo
}
