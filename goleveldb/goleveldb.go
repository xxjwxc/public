package goleveldb

import (
	"os"

	"../../data/config"
	"../log"
	"github.com/syndtr/goleveldb/leveldb"
)

var m_db *leveldb.DB = nil

func init() {
	Clear();
	creat()
}

func creat(){
	if m_db == nil{
		var err error
		m_db, err = leveldb.OpenFile(config.GetLevelDbDir(), nil)
		if err != nil {
			log.Print(log.Log_Error, err.Error())
			m_db.Close()
			m_db = nil
		}
	}
}

/*
 清空数据
*/
func Clear() bool {
	Close();
	os.RemoveAll(config.GetLevelDbDir());
	return true;
}


/*
 关闭
*/
func Close(){
	if m_db != nil{
		m_db.Close()
		m_db = nil;
	}
}

/*
*获取
 */
func Get(key []byte) (data []byte, err error) {
	data = nil
	data, err = m_db.Get(key, nil)
	return
}

/*
设置
*/
func Set(key, value []byte) error {
	err := m_db.Put(key, value, nil)
	return err
}

/*
 删除
*/
func Delete(key []byte) error {
	return m_db.Delete(key, nil)
}

/*
*获取
 */
func Getkv(key string) (data []byte, err error) {
	return Get([]byte(key))
}

/*
*获取
 */
func Setkv(key string, value []byte) error {
	return Set([]byte(key), value)
}

/*
 删除
*/
func Deletekv(key string) error {
	return Delete([]byte(key))
}


