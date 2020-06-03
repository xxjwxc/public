package mynsq

import (
	"data/config"
	"log"
	"public/mylog"

	"github.com/nsqio/go-nsq"
)

var producer *nsq.Producer = nil
var consumerMap map[int]*nsq.Consumer = nil
var cosumerTopics map[string][]int = nil

func init() {
	consumerMap = make(map[int]*nsq.Consumer)
	cosumerTopics = make(map[string][]int)
	cnf := nsq.NewConfig()
	var err error
	producer, err = nsq.NewProducer(config.GetNsqAddr(), cnf)
	if err != nil {
		mylog.Error(err)
		panic(err)
	}

}

//发消息
func ProduceMsg(topic string, message []byte) bool {

	if producer == nil {
		//channel 锁住
		cnf := nsq.NewConfig()
		var err error
		producer, err = nsq.NewProducer(config.GetNsqAddr(), cnf)
		if err != nil {
			mylog.Error(err)
			return false
		}
	}

	if producer != nil {
		err := producer.Publish(topic, message)
		if err != nil {
			mylog.Error(err)
			return false
		} else {
			return true
		}
	}

	return false
}

//单个nsqd处理消息
//index 表示consumer 索引(用于开关使用)
func StartConsumeMsg(index int, topic, channel, nsqd string, handler nsq.Handler) bool {
	StopConsumeMsgByIndex(index)
	//第一次初始化 进入
	if consumerMap[index] == nil {
		conf := nsq.NewConfig()
		//最大允许向两台NSQD服务器接受消息，默认是1
		//config.MaxInFlight = 2
		var err error
		consumerMap[index], err = nsq.NewConsumer(topic, channel, conf)
		if nil != err {
			log.Println(err)
			mylog.Error(err)
			return false
		}

		//开始正式启动(后台，非阻塞方式)
		consumerMap[index].AddHandler(handler)
		err = consumerMap[index].ConnectToNSQD(nsqd)
		if nil != err {
			log.Println(err)
			mylog.Error(err)
			return false
		}

		cosumerTopics[topic] = append(cosumerTopics[topic], index)
		return true
	}

	return false
}

func GetConsumeSize() int {
	return len(consumerMap)
}

//停止消费
func StopConsumeMsgByIndex(index int) {
	if consumerMap[index] != nil {
		consumerMap[index].Stop()
		consumerMap[index] = nil

		for k, v := range cosumerTopics {
			for i := range v {
				if v[i] == index {
					cosumerTopics[k] = append(v[:i], v[i+1:]...)
					break
				}
			}
		}
	}
}

//停止某个topic下全部消费
func StopConsumeMsgByTopic(topic string) {
	v, ok := cosumerTopics[topic]
	if ok {
		for i := range v {
			if consumerMap[v[i]] != nil {
				consumerMap[v[i]].Stop()
				consumerMap[v[i]] = nil
			}
		}
	}

	var temp []int
	cosumerTopics[topic] = temp
}

//停止所有topic的全部消费
func StopAllConsumeMsg() {
	cosumerTopics = make(map[string][]int)
	size := GetConsumeSize()
	for i := 0; i < size; i++ {
		StopConsumeMsgByIndex(i)
	}

	consumerMap = make(map[int]*nsq.Consumer)
}

//func TestNSQ() {
//	for k, v := range cosumerTopics {
//		log.Println("topic is", k)
//		for _, vv := range v {
//			log.Println("consumer  index is", vv)
//			log.Println("consumer is ", consumerMap[vv])
//		}
//	}
//}
