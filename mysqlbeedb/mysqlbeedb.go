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
package mysqlbeedb

import (
	"database/sql"
	"fmt"
	"public/mylog"

	"data/config"

	"github.com/astaxie/beedb"
	_ "github.com/go-sql-driver/mysql"
)

type MySqlDB struct {
	DB *sql.DB
}

func (i *MySqlDB) OnGetDBOrm(dataSourceName string) (orm beedb.Model) {
	if i.DB == nil {
		var err error
		i.DB, err = sql.Open("mysql", dataSourceName)
		if err != nil {
			mylog.Print(mylog.Log_Error, fmt.Sprintf("Got error when connect database, the error is '%v'", err))
		}
	}

	orm = beedb.New(i.DB)

	if config.OnIsDev() {
		beedb.OnDebug = true
	}
	return
}

func (i *MySqlDB) OnDestoryDB() {
	if i.DB != nil {
		i.DB.Close()
		i.DB = nil
	}
}
