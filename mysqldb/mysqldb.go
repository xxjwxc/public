package mysqldb

import (
	"fmt"
	"public/mylog"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type MySqlDB struct {
	*gorm.DB
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
			mylog.Print(mylog.Log_Error, fmt.Sprintf("Got error when connect database, the error is '%v'", err))
		}
	}

	i.DB.SingularTable(true) //全局禁用表名复数
	orm = i.DB
	if isDev {
		i.DB.LogMode(true)
		//beedb.OnDebug = true
	} else {
		i.DB.SetLogger(DbLog{})
	}
	return
}

var isDev bool = false

//是否调试
func (i *MySqlDB) SetIsDev(b bool) {
	isDev = b
}

func (i *MySqlDB) OnDestoryDB() {
	if i.DB != nil {
		i.DB.Close()
		i.DB = nil
	}
}
