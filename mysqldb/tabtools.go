package mysqldb

import (
	"fmt"
	"regexp"
	"time"

	"github.com/xxjwxc/public/mylog"
	"gorm.io/gorm"
)

// TabColumnInfo 表信息
type TabColumnInfo struct {
	ColumnName string
	ColumnType string // 跟数据库一致
	Len        string
	NotNull    bool
	Comment    string // 注释
}

type TablesModel struct {
	ID int `gorm:"primaryKey;column:id" json:"-"` // 主键id
}

type TablesTools struct {
	tabName string
	orm     *MySqlDB
}

// Tables
func NewTabTools(orm *MySqlDB, tabName string) (*TablesTools, error) {
	return &TablesTools{
		orm:     orm,
		tabName: tabName,
	}, nil
}

// GetDB 获取db
func (t *TablesTools) GetDB() *gorm.DB {
	return t.orm.Table(t.tabName)
}

// CreateTable 创建表
func (t *TablesTools) CreateTable(columns []*TabColumnInfo) error {
	err := t.orm.Table(t.tabName).Set("gorm:table_options", "ENGINE=InnoDB").Migrator().CreateTable(&TablesModel{})
	if err != nil {
		return err
	}

	for _, v := range columns {
		err := t.AddColumn(v)
		if err != nil {
			return err
		}
		// notnul := ""
		// if v.NotNull {
		// 	notnul = "NOT NULL"
		// }
		// _len := v.Len
		// if len(v.Len) > 0 {
		// 	_len = fmt.Sprintf("(%v)", v.Len)
		// }

		// sql := fmt.Sprintf("ALTER TABLE `%v` ADD COLUMN `%v` %v%v %v COMMENT '%v';", t.tabName, v.ColumnName, v.ColumnType, _len, notnul, v.Comment)
		// err = t.orm.Exec(sql).Error
		// if err != nil {
		// 	return err
		// }
	}

	return err
}

// HasTable 表是否存在
func (t *TablesTools) HasTable() bool {
	return t.orm.Migrator().HasTable(t.tabName)
}

// DropTable 如果存在表则删除（删除时会忽略、删除外键约束)
func (t *TablesTools) DropTable() error {
	return t.orm.Migrator().DropTable(t.tabName)
}

// TruncateTable 截断表
func (t *TablesTools) TruncateTable() error {
	return t.orm.Exec(fmt.Sprintf("TRUNCATE TABLE %v;", t.tabName)).Error
}

// GetTables  获取所有表名
func (t *TablesTools) GetTables() (tableList []string, err error) {
	return t.orm.Migrator().GetTables()
}

// RenameTable 重命名表
func (t *TablesTools) RenameTable(newTabName string) error {
	return t.orm.Migrator().RenameTable(t.tabName, newTabName)
}

// AddColumn 添加列元素
func (t *TablesTools) AddColumn(column *TabColumnInfo) error {
	notnul := ""
	if column.NotNull {
		notnul = "NOT NULL"
	}
	_len := column.Len
	if len(column.Len) > 0 {
		_len = fmt.Sprintf("(%v)", column.Len)
	}

	sql := fmt.Sprintf("ALTER TABLE `%v` ADD COLUMN `%v` %v%v %v COMMENT '%v';", t.tabName, column.ColumnName, column.ColumnType, _len, notnul, column.Comment)
	err := t.orm.Exec(sql).Error
	if err != nil {
		return err
	}

	return nil
}

// HasColumn 是否有列
func (t *TablesTools) HasColumn(column string) bool {
	return t.orm.Table(t.tabName).Migrator().HasColumn(&TablesModel{}, column)
}

// DropColumn 删除列
func (t *TablesTools) DropColumn(column string) error {
	return t.orm.Table(t.tabName).Migrator().DropColumn(&TablesModel{}, column)
}

// RenameColumn 字段重命名
func (t *TablesTools) RenameColumn(oldColumn, newColumn string) error {
	return t.orm.Table(t.tabName).Migrator().RenameColumn(&TablesModel{}, oldColumn, newColumn)
}

// ColumnTypes 获取列属性
func (t *TablesTools) ColumnTypes() ([]gorm.ColumnType, error) {
	return t.orm.Table(t.tabName).Migrator().ColumnTypes(&TablesModel{})
}

// CreateIndex Indexes
func (t *TablesTools) CreateIndex(column string) error {
	return t.orm.Table(t.tabName).Migrator().CreateIndex(&TablesModel{}, column)
}

// DropIndex Indexes
func (t *TablesTools) DropIndex(column string) error {
	return t.orm.Table(t.tabName).Migrator().DropIndex(&TablesModel{}, column)
}

// HasIndex Indexes
func (t *TablesTools) HasIndex(column string) bool {
	return t.orm.Table(t.tabName).Migrator().HasIndex(&TablesModel{}, column)
}

// RenameIndex Indexes
func (t *TablesTools) RenameIndex(oldColumn, newColumn string) error {
	return t.orm.Table(t.tabName).Migrator().RenameIndex(&TablesModel{}, oldColumn, newColumn)
}

type ColumnTypeInfo struct {
	Name string
	Type string
	Desc string
}

// RawInfo 列信息
type RawInfo struct {
	ColumnTypes []ColumnTypeInfo
	Values      [][]interface{}
}

// Select 查询
func (t *TablesTools) Select(DB *gorm.DB) (*RawInfo, error) {
	colTypes, err := t.ColumnTypes()
	if err != nil {
		return nil, err
	}

	out := &RawInfo{}
	for i, v := range colTypes {
		commont, _ := v.Comment()
		out.ColumnTypes = append(out.ColumnTypes, ColumnTypeInfo{
			Name: v.Name(),
			Type: ColumnType(colTypes[i]),
			Desc: commont,
		})
	}

	rows, err := DB.Table(t.tabName).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var values []interface{} //创建一个与列的数量相当的空接口
		for _, v := range out.ColumnTypes {
			deserialize(v.Type, &values)
		}

		err := rows.Scan(values...) //开始读行，Scan函数只接受指针变量
		if err != nil {
			return nil, err
		}
		out.Values = append(out.Values, values) //将单行所有列的键值对附加在总的返回值上（以行为单位）
	}

	return out, err
}

func deserialize(tp string, data *[]interface{}) interface{} {
	switch tp {
	case "int":
		var tmp int
		*data = append(*data, &tmp)
	case "int64":
		var tmp int64
		*data = append(*data, &tmp)
	case "string":
		var tmp string
		*data = append(*data, &tmp)
	case "[]byte":
		var tmp []byte
		*data = append(*data, &tmp)
	case "time":
		var tmp time.Time
		*data = append(*data, &tmp)
	case "bool":
		var tmp bool
		*data = append(*data, &tmp)
	case "float64":
		var tmp float64
		*data = append(*data, &tmp)
	default:
		mylog.Errorf("type (%v) not match in any way.maybe need to add on", tp)
	}

	return nil
}

// "VARCHAR", "TEXT", "NVARCHAR", "DECIMAL", "BOOL",
// "INT", and "BIGINT".
func ColumnType(t gorm.ColumnType) string {
	tp, _ := t.ColumnType()
	// Precise matching first.先精确匹配
	if v, ok := TypeMysqlDicMp[tp]; ok {
		return v
	}

	// Fuzzy Regular Matching.模糊正则匹配
	for _, l := range TypeMysqlMatchList {
		if ok, _ := regexp.MatchString(l.Key, tp); ok {
			return l.Value
		}
	}

	return tp
}

// TypeMysqlDicMp Accurate matching type.精确匹配类型
var TypeMysqlDicMp = map[string]string{
	"int":                "int",
	"int unsigned":       "int",
	"tinyint":            "int",
	"tinyint unsigned":   "int",
	"mediumint":          "int",
	"mediumint unsigned": "int",

	"smallint":          "int64",
	"smallint unsigned": "int64",
	"bigint":            "int64",
	"bigint unsigned":   "int64",
	"timestamp":         "int64",
	"integer":           "int64",

	"varchar":    "string",
	"char":       "string",
	"json":       "string",
	"text":       "string",
	"mediumtext": "string",
	"longtext":   "string",
	"tinytext":   "string",
	"enum":       "string",
	"nvarchar":   "string",

	"bit(1)":     "[]byte",
	"tinyblob":   "[]byte",
	"blob":       "[]byte",
	"mediumblob": "[]byte",
	"longblob":   "[]byte",
	"binary":     "[]byte",

	"date":          "time",
	"datetime":      "time",
	"time":          "time",
	"smalldatetime": "time", //sqlserver

	"tinyint(1)":          "bool", // tinyint(1) 默认设置成bool
	"tinyint(1) unsigned": "bool", // tinyint(1) 默认设置成bool

	"double":          "float64",
	"double unsigned": "float64",
	"float":           "float64",
	"float unsigned":  "float64",
	"real":            "float64",
	"numeric":         "float64",
}

// TypeMysqlMatchList Fuzzy Matching Types.模糊匹配类型
var TypeMysqlMatchList = []struct {
	Key   string
	Value string
}{
	{`^(tinyint)[(]\d+[)] unsigned`, "int"},
	{`^(smallint)[(]\d+[)] unsigned`, "int"},
	{`^(int)[(]\d+[)] unsigned`, "int"},
	{`^(tinyint)[(]\d+[)]`, "int"},
	{`^(smallint)[(]\d+[)]`, "int"},
	{`^(int)[(]\d+[)]`, "int"},
	{`^(mediumint)[(]\d+[)]`, "int"},
	{`^(mediumint)[(]\d+[)] unsigned`, "int"},
	{`^(integer)[(]\d+[)]`, "int"},

	{`^(bigint)[(]\d+[)] unsigned`, "int64"},
	{`^(bigint)[(]\d+[)]`, "int64"},
	{`^(timestamp)[(]\d+[)]`, "int64"},

	{`^(float)[(]\d+,\d+[)] unsigned`, "float64"},
	{`^(double)[(]\d+,\d+[)] unsigned`, "float64"},
	{`^(decimal)[(]\d+,\d+[)]`, "float64"},
	{`^(double)[(]\d+,\d+[)]`, "float64"},
	{`^(float)[(]\d+,\d+[)]`, "float64"},

	{`^(char)[(]\d+[)]`, "string"},
	{`^(enum)[(](.)+[)]`, "string"},
	{`^(varchar)[(]\d+[)]`, "string"},
	{`^(text)[(]\d+[)]`, "string"},
	{`^(set)[(][\s\S]+[)]`, "string"},

	{`^(varbinary)[(]\d+[)]`, "[]byte"},
	{`^(blob)[(]\d+[)]`, "[]byte"},
	{`^(binary)[(]\d+[)]`, "[]byte"},
	{`^(bit)[(]\d+[)]`, "[]byte"},
	{`^(geometry)[(]\d+[)]`, "[]byte"},

	{`^(datetime)[(]\d+[)]`, "time"},
}
