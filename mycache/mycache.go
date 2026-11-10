/*
key/value 内存缓存，支持基于超时的自动无效功能
*/
package mycache

import (
	"fmt"
	"sync"
	"time"

	"github.com/xxjwxc/public/mycache/cache2go"
	"github.com/xxjwxc/public/serializing"
)

// CacheIFS 缓存操作接口
type CacheIFS interface {
	Destory()                                                                       // 析构
	Add(key interface{}, value interface{}, lifeSpan time.Duration) error           // 添加一个元素
	Value(key interface{}, value interface{}) error                                 // 获取一个value
	IsExist(key interface{}) bool                                                   // 判断是否存在
	Delete(key interface{}) error                                                   // 删除一个
	Clear() error                                                                   // 清空
	Close() (err error)                                                             // 关闭连接
	TryLock(key interface{}, value interface{}, lifeSpan time.Duration) (err error) //  试着加锁
	Unlock(key interface{}) (err error)                                             // 解锁
	GetKeyS(key interface{}) ([]string, error)                                      // 查询所有key
	Refresh(key interface{}, lifeSpan time.Duration) error                          // 更新时间
}

// MyCache 内存缓存
type MyCache struct {
	cache *cache2go.CacheTable
	// mtx 保留字段以兼容旧的 MyCache 字面量构造方式。
	// 缓存操作本身是线程安全的(cache2go 内部自带读写锁),
	// 且 TryLock 已改为原子操作,不再依赖该互斥量。
	mtx sync.Mutex
}

// NewCache 初始化一个cache,cachename 缓存名字
func NewCache(cachename string) (mc *MyCache) {
	mc = &MyCache{}
	mc.cache = cache2go.Cache(cachename)
	return
}

// Destory 添加一个缓存,lifeSpan:缓存时间，0表示永不超时
func (mc *MyCache) Destory() {

}

// Add 添加一个缓存,lifeSpan:缓存时间，0表示永不超时
func (mc *MyCache) Add(key interface{}, value interface{}, lifeSpan time.Duration) error {
	data, err := encodeValue(value)
	if err != nil {
		return err
	}
	mc.cache.Add(key, lifeSpan, data)
	return nil
}

// Value 查找一个cache
func (mc *MyCache) Value(key interface{}, value interface{}) error {
	res, err := mc.cache.Peek(key)
	if err != nil {
		return err
	}
	// 使用安全断言:通过 GetCache2go 写入非 []byte 数据时不再 panic
	bt, ok := res.Data().([]byte)
	if !ok {
		return fmt.Errorf("mycache: key[%v] 缓存值类型为 %T,不是 []byte", key, res.Data())
	}
	return decodeValue(bt, value)
}

// IsExist 	判断key是否存在
func (mc *MyCache) IsExist(key interface{}) bool {
	return mc.cache.Exists(key)
}

// Delete 删除一个cache
func (mc *MyCache) Delete(key interface{}) error {
	_, err := mc.cache.Delete(key)
	return err
}

// GetCache2go 获取原始cache2go操作类
func (mc *MyCache) GetCache2go() *cache2go.CacheTable {
	return mc.cache
}

// GetKeyS 查询所有key
func (mc *MyCache) GetKeyS(key interface{}) ([]string, error) {
	return []string{}, nil
}

// Clear 清空表內容
func (mc *MyCache) Clear() error {
	mc.cache.Flush()
	return nil
}

// Close 清空表內容
func (mc *MyCache) Close() error {
	return nil
}

// TryLock 试着枷锁
// 使用 NotFoundAdd 完成"不存在才写入"的原子操作,避免先 Exists 再 Add 的竞态
func (mc *MyCache) TryLock(key interface{}, value interface{}, lifeSpan time.Duration) (err error) {
	data, err := encodeValue(value)
	if err != nil {
		return err
	}
	if !mc.cache.NotFoundAdd(key, lifeSpan, data) { // 已存在，枷锁失败
		return fmt.Errorf("lock fail")
	}
	return nil
}

// Unlock 试着解锁
func (mc *MyCache) Unlock(key interface{}) (err error) {
	_, err = mc.cache.Delete(key)
	return err
}

// Refresh 更新时间
func (mc *MyCache) Refresh(key interface{}, lifeSpan time.Duration) error {
	return mc.cache.Refresh(key, lifeSpan)
}

func encodeValue(value interface{}) ([]byte, error) {
	return serializing.Encode(value)
}

func decodeValue(in []byte, out interface{}) (err error) {
	return serializing.Decode(in, out)
}
