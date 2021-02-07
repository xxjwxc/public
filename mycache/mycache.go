/*
 key/value 内存缓存，支持基于超时的自动无效功能
*/
package mycache

import (
	"time"

	"github.com/muesli/cache2go"
	"github.com/xxjwxc/public/serializing"
)

// MyCache 内存缓存
type MyCache struct {
	cache *cache2go.CacheTable
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
	mc.cache.Add(key, lifeSpan, encodeValue(value))
	return nil
}

// Value 查找一个cache
func (mc *MyCache) Value(key interface{}, value interface{}) error {
	res, err := mc.cache.Value(key)
	if err == nil {
		bt := res.Data().([]byte)
		return decodeValue(bt, value)
	}
	return err
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

// Clear 清空表內容
func (mc *MyCache) Clear() error {
	mc.cache.Flush()
	return nil
}

// Close 清空表內容
func (mc *MyCache) Close() error {
	return nil
}

func encodeValue(value interface{}) []byte {
	data, _ := serializing.Encode(value)
	return data
}

func decodeValue(in []byte, out interface{}) (err error) {
	return serializing.Decode(in, out)
}
