package myleveldb

import (
	"reflect"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/xxjwxc/public/mylog"
	"github.com/xxjwxc/public/tools"
)

var lock sync.Mutex
var locks = map[string]*sync.Mutex{}

// Param kv
type Param struct {
	Key   string
	Value interface{}
}

// OnInitDB 初始化
func OnInitDB(dataSourceName string) MyLevelDB {
	if _, ok := locks[dataSourceName]; !ok {
		lock.Lock()
		if _, ok := locks[dataSourceName]; !ok {
			locks[dataSourceName] = &sync.Mutex{}
		}
		lock.Unlock()
	}

	locks[dataSourceName].Lock()
	var L MyLevelDB
	L.dataSourceName = dataSourceName
	L.DB, L.E = leveldb.OpenFile(dataSourceName, nil)
	if L.E != nil {
		locks[dataSourceName].Unlock()
		mylog.Error(L.E)
	}
	//	L.op = &opt.ReadOptions{
	//		false,
	//		opt.NoStrict,
	//	}
	return L
}

// MyLevelDB ...
type MyLevelDB struct {
	DB             *leveldb.DB
	E              error
	dataSourceName string
	//op    *opt.ReadOptions
	Value interface{}
}

// OnDestoryDB 关闭
func (L *MyLevelDB) OnDestoryDB() {
	L.Close()
}

// Close 关闭
func (L *MyLevelDB) Close() {
	if L.DB != nil {
		L.DB.Close()
		L.DB = nil
		locks[L.dataSourceName].Unlock()
	}
}

// Get 获取数据
func (L *MyLevelDB) Get(key string, value interface{}) (b bool) {
	if L.DB != nil {
		var err error
		var by []byte
		if by, err = L.DB.Get([]byte(key), nil /*L.op*/); err != nil {
			//mylog.Error(err)
		} else {
			if err := tools.DecodeByte(by, value); err != nil {
				//错误处理
				mylog.Error(err)
			} else {
				return true
			}
		}
	}

	return false
}

// Model ...
func (L *MyLevelDB) Model(refs interface{}) *MyLevelDB {
	if reflect.ValueOf(refs).Type().Kind() == reflect.Ptr {
		mylog.ErrorString("Model: attempt to Model into a non-pointer")
		panic(0)
	}
	L.Value = refs
	return L
}

// Find 模糊查找
/*
	t value的类型
	values 为返回结果
	args 传一个参数:表示模糊搜索
	args 传2个参数:表示范围搜索
*/
func (L *MyLevelDB) Find(values *[]Param, args ...string) (b bool) {
	if L.DB != nil && L.Value != nil {
		n := len(args)
		var it iterator.Iterator

		if n == 1 { //模糊查找
			it = L.DB.NewIterator(util.BytesPrefix([]byte(args[0])), nil)
		} else {
			it = L.DB.NewIterator(&util.Range{Start: []byte(args[0]), Limit: []byte(args[1])}, nil)
		}

		for it.Next() {
			tmp := L.Value
			if err := tools.DecodeByte(it.Value(), tmp); err != nil {
				//错误处理
				mylog.Error(err)
			}
			*values = append(*values, Param{string(it.Key()), tmp})
		}

		it.Release()
		//iter := L.DB.NewIterator(nil, nil)
	} else {
		if L.Value == nil {
			panic("not call Model()")
		}
		mylog.ErrorString("not init.")
	}

	return false
}

// Add 添加数据
//注意：只支持基础类型
func (L *MyLevelDB) Add(key string, value interface{}) bool {
	if L.DB != nil {
		by, err := tools.EncodeByte(value)
		if err != nil {
			//错误处理
			mylog.Error(err)
			return false
		}
		if err = L.DB.Put([]byte(key), by, nil); err != nil {
			mylog.Error(err)
		} else {
			return true
		}
	}
	return false
}

// AddList 添加一组数据(比一个一个添加速度快很多)
//注意：只支持基础类型
func (L *MyLevelDB) AddList(array []Param) bool {
	if L.DB != nil {
		batch := new(leveldb.Batch)
		for _, p := range array {
			by, err := tools.EncodeByte(p.Value)
			if err != nil {
				//错误处理
				mylog.Error(err)
				return false
			}
			batch.Put([]byte(p.Key), by)
		}
		err := L.DB.Write(batch, nil)
		if err != nil {
			//错误处理
			mylog.Error(err)
			return false
		}
		return true
	}
	return false
}

// Delete 删除
func (L *MyLevelDB) Delete(key string) bool {
	if L.DB != nil {
		err := L.DB.Delete([]byte(key), nil)
		if err != nil {
			mylog.Error(err)
			return false
		}

		return true
	}
	return false
}
