package sign

import (
	"strings"
)

func init() {
	OnInit()
}

func OnInit() {
	// str_db := config.GetDbUrl()
	// if len(str_db) > 0 {
	// 	var db mysqldb.MySqlDB
	// 	defer db.OnDestoryDB()
	// 	orm := db.OnGetDBOrm(str_db)
	// 	if orm.HasTable(&Sign_client_tbl{}) { //有这个表
	// 		now := time.Now()
	// 		var list []Sign_client_tbl
	// 		err := orm.Where("expire_time > ?", now).Find(&list).Error
	// 		if err != nil {
	// 			mylog.Error(err)
	// 			return
	// 		}
	// 		cache := mycache.NewCache(_sign_data)
	// 		for _, v := range list { //保存数据到缓存
	// 			cache.Add(v.App_key, v, v.Expire_time.Sub(now))
	// 		}
	// 	}
	// }
}

func getOne(appKey string) (sign Sign_client_tbl) {
	// str_db := config.GetDbUrl()
	// if len(str_db) > 0 {
	// 	var db mysqldb.MySqlDB
	// 	defer db.OnDestoryDB()
	// 	orm := db.OnGetDBOrm(str_db)
	// 	if orm.HasTable(&Sign_client_tbl{}) { //有这个表
	// 		now := time.Now()
	// 		err := orm.Where("app_key = ? and expire_time > ?", appKey, now).Find(&sign).Error
	// 		if err != nil {
	// 			mylog.Error(err)
	// 			return
	// 		}
	// 		cache := mycache.NewCache(_sign_data)
	// 		cache.Add(sign.App_key, sign, sign.Expire_time.Sub(now))
	// 	}
	// }

	return
}

/*
 生成验签
*/
func OnGetSign(appkey string, parm ...interface{}) string {
	// var sign Sign_client_tbl
	// if len(appkey) > 0 {
	// 	cache := mycache.NewCache(_sign_data)
	// 	tp, b := cache.Value(appkey)
	// 	if b {
	// 		sign = tp.(Sign_client_tbl)
	// 	} else {
	// 		sign = getOne(appkey)
	// 	}
	// }

	// if sign.Id == 0 {
	// 	return ""
	// }

	// //开始验签
	// var strKey string
	// for _, v := range parm {
	// 	strKey += tools.AsString(v)
	// }

	// return tools.Md5Encoder(strKey)
	return ""
}

/*
 开始验签
*/
func OnCheckSign(appkey, signature string, parm ...interface{}) bool {
	return strings.EqualFold(signature, OnGetSign(appkey, parm))
}
