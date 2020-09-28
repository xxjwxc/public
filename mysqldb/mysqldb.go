package mysqldb

import (
	"errors"

	"github.com/xxjwxc/public/dev"
	myerrors "github.com/xxjwxc/public/errors"
	"gorm.io/gorm/logger"

	"github.com/xxjwxc/public/mylog"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// MySqlDB ...
type MySqlDB struct {
	*gorm.DB
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
		i.DB, err = gorm.Open(mysql.Open(dataSourceName), &gorm.Config{PrepareStmt: false,
			Logger: logger.Default})
		if err != nil {
			mylog.Error(myerrors.Wrap(err, "Got error when connect database:"+dataSourceName))
			return nil
		}
	}

	// i.DB.SingularTable(true) //全局禁用表名复数
	if dev.IsDev() {
		i.DB = i.DB.Debug()
	} else {
		i.DB.Logger = GetDBlog()
	}

	return i.DB
}

// OnDestoryDB destorydb
func (i *MySqlDB) OnDestoryDB() {
	if i.DB != nil {
		sqldb, _ := i.DB.DB()
		sqldb.Close()
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

// RecordNotFound check if returning ErrRecordNotFound error
func (i *MySqlDB) RecordNotFound() bool {
	return !errors.Is(i.Error, gorm.ErrRecordNotFound)
}

// Commit 自动提交(如果有错，Rollback)
func (i *MySqlDB) Commit(db *gorm.DB) {
	if db.Error != nil {
		db.Rollback() // 回滚
	} else {
		db.Commit()
	}
}
