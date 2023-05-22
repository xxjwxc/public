package myrabbitmq

import (
	"fmt"
	"testing"
	"time"

	"github.com/xxjwxc/public/message"
)

func Test_NewAdmin(t *testing.T) {
	topic := "xxjtest"
	tag := "tagtest"
	host := "amqp://admin:admin@192.155.1.151:5672/"
	group := "nlp_cmd_train"
	// 初始化mq
	mq, err := NewRabbitMQ(host)
	if err != nil {
		fmt.Println(err)
	}
	defer mq.Close() // 完成任务释放资源
	_, err = mq.NewProducer(topic, tag, group)
	fmt.Println(err)
	// go func() {
	// 	for i := 0; i < 1000; i++ {
	// 		mq.SendMessage([]byte(fmt.Sprintf("this is xxj test %v", i)))
	// 		time.Sleep(1 * time.Microsecond)
	// 	}
	// }()
	// time.Sleep(10 * time.Minute)

	go func() {
		ch, err := mq.NewConsumer(topic) // 消費者
		fmt.Println(err)
		mq.Start(ch, func(msg []byte) error {
			fmt.Println("消费者1", string(msg))
			time.Sleep(1 * time.Second)
			return nil
		})
		time.Sleep(10 * time.Minute)
	}()

	go func() {
		ch, err := mq.NewConsumer(topic) // 消費者
		fmt.Println(err)
		mq.Start(ch, func(msg []byte) error {
			fmt.Println("消费者2", string(msg))
			time.Sleep(1 * time.Second)
			return message.GetError(message.ActvFailure)
		})
		time.Sleep(10 * time.Minute)

	}()

	time.Sleep(10 * time.Minute)
	mq.Close()
}
