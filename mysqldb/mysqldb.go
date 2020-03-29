package mysqldb

import (
	"github.com/xxjwxc/public/dev"
	"github.com/xxjwxc/public/errors"

	"github.com/xxjwxc/public/mylog"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type MySqlDB struct {
	*gorm.DB
	IsInit bool
}

func OnInitDBOrm(dataSourceName string) (orm *MySqlDB) {
	orm = new(MySqlDB)
	orm.OnGetDBOrm(dataSourceName)
	return
}

func (i *MySqlDB) OnGetDBOrm(dataSourceName string) (orm *gorm.DB) {
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
	if dev.OnIsDev() {
		i.DB.LogMode(true)
		//beedb.OnDebug = true
	} else {
		i.DB.SetLogger(DbLog{})
	}
	orm = i.DB
	return
}

func (i *MySqlDB) OnDestoryDB() {
	if i.DB != nil {
		i.DB.Close()
		i.DB = nil
	}
}
