package weixin

import (
	"time"

	"github.com/xxjwxc/public/mycache"
)

// Gocache Memcache struct contains *memcache.Client
type Gocache struct {
	mc *mycache.MyCache
}

//NewGocache create new cache2go
func NewGocache(server string) *Gocache {
	mc := mycache.OnGetCache(server)
	return &Gocache{&mc}
}

//Get return cached value
func (mem *Gocache) Get(key string) interface{} {
	v, _ := mem.mc.Value(key)
	return v
}

// IsExist check value exists in memcache.
func (mem *Gocache) IsExist(key string) bool {
	return mem.mc.IsExist(key)
}

//Set cached value with key and expire time.
func (mem *Gocache) Set(key string, val interface{}, timeout time.Duration) (err error) {
	mem.mc.Add(key, val, timeout)
	return nil
}

//Delete  value in memcache.
func (mem *Gocache) Delete(key string) error {
	return mem.mc.Delete(key)
}
