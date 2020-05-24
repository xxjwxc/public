package myleveldb

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

func Test_order(t *testing.T) {
	fmt.Println("ssss")
	//初始化db
	ldb := OnInitDB("./database")
	defer ldb.OnDestoryDB()

	//	var www Tweet
	//	www.Location = "xiexiaojun"
	//	www.Created = time.Now()
	//	www.Tags = append(www.Tags, "12334444", "ssss", "wwwww")
	//	b := ldb.Add([]byte("wwww"), www)
	//	fmt.Println(b)
	//	fmt.Println(www.Tags)

	var eee Tweet
	bb := ldb.Get("wwww", &eee)
	fmt.Println(bb)
	fmt.Println(eee.Location)
	fmt.Println(eee.Created)
	fmt.Println(eee.Tags)

	//	var temp []myleveldb.Param
	//	for i := 0; i < 10; i++ {
	//		var www Tweet
	//		www.Location = "xiexiaojun" + strconv.Itoa(i)
	//		www.Created = time.Now()
	//		temp = append(temp, myleveldb.Param{"key" + strconv.Itoa(i), www})
	//	}
	//bbb := ldb.AddList(temp)
	//fmt.Println(bbb)

	bb = ldb.Get("key9", &eee)
	fmt.Println(bb)
	fmt.Println(eee.Location)
	fmt.Println(eee.Created)
	fmt.Println(eee.Tags)

	var tmp []Param
	bb = ldb.Model(&Tweet{}).Find(&tmp, "key3", "key7") //查找

	fmt.Println(tmp)

	return
}
