package mysqldb

import (
	"fmt"
)

// TabColumnInfo 表信息
type TabColumnInfo struct {
	ColumnName string
	ColumnType string // 跟数据库一致
	Len        string
	NotNull    bool
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

// CreateTable 创建表
func (t *TablesTools) CreateTable(columns []*TabColumnInfo) error {
	err := t.orm.Table(t.tabName).Set("gorm:table_options", "ENGINE=InnoDB").Migrator().CreateTable(&TablesModel{})
	if err != nil {
		return err
	}

	for _, v := range columns {
		notnul := ""
		if v.NotNull {
			notnul = "NOT NULL"
		}
		_len := v.Len
		if len(v.Len) > 0 {
			_len = fmt.Sprintf("(%v)", v.Len)
		}

		sql := fmt.Sprintf("ALTER TABLE `%v` ADD COLUMN `%v` %v%v %v ;", t.tabName, v.ColumnName, v.ColumnType, _len, notnul)
		err = t.orm.Exec(sql).Error
		if err != nil {
			return err
		}
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

	sql := fmt.Sprintf("ALTER TABLE `%v` ADD COLUMN `%v` %v%v %v ;", t.tabName, column.ColumnName, column.ColumnType, _len, notnul)
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

//
