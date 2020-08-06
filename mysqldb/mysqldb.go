package mysqldb

import (
	"github.com/xxjwxc/public/dev"
	"github.com/xxjwxc/public/errors"

	"github.com/xxjwxc/public/mylog"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// MySqlDB ...
type MySqlDB struct {
	*gorm.DB
	IsInit bool
}

// OnInitDBOrm init MySqlDB
func OnInitDBOrm(dataSourceName string) (orm *MySqlDB) {
	orm = new(MySqlDB)
	orm.OnGetDBOrm(dataSourceName)
	return
}

// OnGetDBOrm get gorm.db
func (i *MySqlDB) OnGetDBOrm(dataSourceName string) *gorm.DB {
	if i.DB == nil {
		var err error
		i.DB, err = gorm.Open("mysql", dataSourceName)
		if err != nil {
			mylog.Error(errors.Wrap(err, "Got error when connect database:"+dataSourceName))
			return nil
		}
		i.IsInit = true
	}

	i.DB.SingularTable(true) //全局禁用表名复数
	if dev.IsDev() {
		i.DB.LogMode(true)
		//beedb.OnDebug = true
	} else {
		i.DB.SetLogger(DbLog{})
	}

	return i.DB
}

// OnDestoryDB destorydb
func (i *MySqlDB) OnDestoryDB() {
	if i.DB != nil {
		i.DB.Close()
		i.DB = nil
	}
}

// IsNotFound 判断错误是否未找到
func (i *MySqlDB) IsNotFound(errs ...error) bool {
	if len(errs) > 0 {
		for _, err := range errs {
			if err == gorm.ErrRecordNotFound {
				return true
			}
		}
	}
	return i.RecordNotFound()
}

// Commit 自动提交(如果有错，Rollback)
func (i *MySqlDB) Commit(db *gorm.DB) {
	if db.Error != nil {
		db.Rollback() // 回滚
	} else {
		db.Commit()
	}
}
