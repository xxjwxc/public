package myrocketmq

import (
	"fmt"
	"testing"
	"time"
)

func Test_NewAdmin(t *testing.T) {
	topic := "xxjtest"
	host := []string{"192.155.1.151:9876"}
	group := "nlp_cmd_train"
	adm, err := NewAdmin(host) //
	fmt.Println(err)
	adm.CreateTopic(topic)

	pwd, err := NewProducer(host, group, 2) // 生产者
	fmt.Println(err)
	go func() {
		for {
			pwd.SendMessage(topic, []byte("this is xxj test"), 0)
			time.Sleep(1 * time.Second)
		}

	}()

	cs, err := NewConsumer(host, group) // 消費者
	fmt.Println(err)
	cs.Start(topic, func(msg []byte) {
		fmt.Println(string(msg))
	})

	time.Sleep(10 * time.Second)
	cs.Stop()
	pwd.Stop()
	adm.Close()
	time.Sleep(10 * time.Second)
}
