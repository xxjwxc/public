/*
	orm := db.OnCreatDB()
	var sum int64 = 0
	for {
		sum++
		var user User_account_tbl
		user.Id = sum

		orm.SetTable("user_account_tbls")
		err := orm.Where("id=?", sum).Find(&user)
		if err != nil {
			log.Println("-----------:", err)
		} else {
			log.Println(user)
		}

		time.Sleep(time.Second * 2)
	}
*/
package mysqldb

import (
	"fmt"
	"public/mylog"

	"data/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type MySqlDB struct {
	DB *gorm.DB
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

	if config.OnIsDev() {
		orm.LogMode(true)
		//beedb.OnDebug = true
	} else {
		orm.SetLogger(DbLog{})
	}
	return
}

func (i *MySqlDB) OnDestoryDB() {
	if i.DB != nil {
		i.DB.Close()
		i.DB = nil
	}
}

func init() {

}
