package mycache

import (
	"fmt"
	"testing"
	"time"
)

// Tweet is a structure used for serializing/deserializing data in Elasticsearch.
type Tweet struct {
	User     string    `json:"user"`
	Message  string    `json:"message"`
	Retweets int       `json:"retweets"`
	Image    string    `json:"image,omitempty"`
	Created  time.Time `json:"created,omitempty"`
	Tags     []string  `json:"tags,omitempty"`
	Location string    `json:"location,omitempty"`
}

func Test_cache(t *testing.T) {
	//获取
	cache := NewCache("_cache")
	var tp interface{}
	tp, b := cache.Value("key")
	if b {
		tmp := tp.(Tweet)
		fmt.Println(tmp)
	} else {
		var tmp Tweet
		//添加
		cache.Add("key", tmp, 24*time.Hour)
	}

	return
}
