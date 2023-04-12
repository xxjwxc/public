package myrocketmq

import (
	"context"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/admin"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/xxjwxc/public/message"
	"github.com/xxjwxc/public/mylog"
)

type HandlerRocketMqRead func(msg []byte) //ID , MESSAGEID,内容

type MyRocketAdmin struct {
	admin admin.Admin
}

// 创建主题
func NewAdmin(host []string) (*MyRocketAdmin, error) {
	adm, err := admin.NewAdmin(admin.WithResolver(primitive.NewPassthroughResolver(host)))
	if err != nil {
		mylog.Error(err)
		return nil, err
	}

	return &MyRocketAdmin{admin: adm}, nil
}

func (m *MyRocketAdmin) CreateTopic(topic string) error {
	if m.admin == nil {
		return message.GetError(message.StateError)
	}
	return m.admin.CreateTopic(context.Background(), admin.WithTopicCreate(topic))
}

func (m *MyRocketAdmin) Close() error {
	if m.admin == nil {
		return message.GetError(message.StateError)
	}
	return m.admin.Close()
}

type MyRocketProducer struct {
	producer rocketmq.Producer
}

// 创建消費者
func NewProducer(host []string, group string, retry int) (*MyRocketProducer, error) {
	p, err := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver(host)),
		producer.WithRetry(retry),
		producer.WithGroupName(group),
	)
	if err != nil {
		return nil, err
	}
	err = p.Start()
	if err != nil {
		return nil, err
	}

	return &MyRocketProducer{producer: p}, nil
}

// SendMessage 发送消息(level 代表延迟级别)
func (m *MyRocketProducer) SendMessage(topic string, tag string, msg []byte, level int) error {
	if m.producer == nil {
		return message.GetError(message.StateError)
	}

	req := &primitive.Message{
		Topic: topic,
		Body:  msg,
	}

	req.WithTag(tag)
	if level > 0 {
		req.WithDelayTimeLevel(level) // 延迟级别
	}

	_, err := m.producer.SendSync(context.Background(), req)
	return err
}

func (m *MyRocketProducer) Stop() error {
	if m.producer == nil {
		return message.GetError(message.StateError)
	}
	return m.producer.Shutdown()
}

type MyRocketConsumer struct {
	consumer rocketmq.PushConsumer
}

// 消费者
func NewConsumer(host []string, group string) (*MyRocketConsumer, error) {
	c, err := rocketmq.NewPushConsumer(consumer.WithNameServer(host),
		consumer.WithConsumerModel(consumer.Clustering),
		consumer.WithGroupName(group),
	)
	if err != nil {
		return nil, err
	}

	return &MyRocketConsumer{consumer: c}, nil
}

func (m *MyRocketConsumer) Start(topic string, hand HandlerRocketMqRead) error {
	if m.consumer == nil {
		return message.GetError(message.StateError)
	}
	// 订阅topic
	err := m.consumer.Subscribe(topic, consumer.MessageSelector{}, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for _, message := range msgs {
			hand(message.Body)
		}

		return consumer.ConsumeSuccess, nil
	})

	if err != nil {
		return err
	}

	// 启动consumer
	return m.consumer.Start()
}

func (m *MyRocketConsumer) Stop() error {
	if m.consumer == nil {
		return message.GetError(message.StateError)
	}
	return m.consumer.Shutdown()
}
