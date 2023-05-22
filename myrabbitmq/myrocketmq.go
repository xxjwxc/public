package myrabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/xxjwxc/public/message"
)

type HandlerRocketMqRead func(msg []byte) error //ID , MESSAGEID,内容

// MQURL 格式 amqp://账号：密码@rabbitmq服务器地址：端口号/vhost (默认是5672端口)
// 端口可在 /etc/rabbitmq/rabbitmq-env.conf 配置文件设置，也可以启动后通过netstat -tlnp查看
const MQURL = "amqp://admin:huan91uncc@172.21.138.131:5672/"

type MyRabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	// 队列名称
	QueueName string
	// 交换机
	Exchange string
	// routing Key
	RoutingKey string
	//MQ链接字符串
	Mqurl string
}

// NewRabbitMQ 创建结构体实例
func NewRabbitMQ(mqurl string) (*MyRabbitMQ, error) {
	rabbitMQ := MyRabbitMQ{
		Mqurl: mqurl,
	}
	var err error
	//创建rabbitmq连接
	rabbitMQ.Conn, err = amqp.Dial(rabbitMQ.Mqurl)
	if err != nil {
		return nil, err
	}

	//创建Channel
	rabbitMQ.Channel, err = rabbitMQ.Conn.Channel()
	if err != nil {
		return nil, err
	}

	return &rabbitMQ, nil
}

// 释放资源,建议NewRabbitMQ获取实例后 配合defer使用
func (mq *MyRabbitMQ) Close() error {
	if mq.Conn == nil || mq.Channel == nil {
		return message.GetError(message.InValidOp)
	}

	err := mq.Conn.Close()
	if err != nil {
		return err
	}

	return mq.Channel.Close()
}

// NewProducer 创建生产者
func (mq *MyRabbitMQ) NewProducer(queueName, exchange, routingKey string) (*amqp.Channel, error) {
	mq.QueueName = queueName
	mq.Exchange = exchange
	mq.RoutingKey = routingKey
	// 1.声明队列
	/*
	  如果只有一方声明队列，可能会导致下面的情况：
	   a)消费者是无法订阅或者获取不存在的MessageQueue中信息
	   b)消息被Exchange接受以后，如果没有匹配的Queue，则会被丢弃

	  为了避免上面的问题，所以最好选择两方一起声明
	  ps:如果客户端尝试建立一个已经存在的消息队列，Rabbit MQ不会做任何事情，并返回客户端建立成功的
	*/
	_, err := mq.Channel.QueueDeclare( // 返回的队列对象内部记录了队列的一些信息，这里没什么用
		mq.QueueName, // 队列名
		true,         // 是否持久化
		false,        // 是否自动删除(前提是至少有一个消费者连接到这个队列，之后所有与这个队列连接的消费者都断开时，才会自动删除。注意：生产者客户端创建这个队列，或者没有消费者客户端与这个队列连接时，都不会自动删除这个队列)
		false,        // 是否为排他队列（排他的队列仅对“首次”声明的conn可见[一个conn中的其他channel也能访问该队列]，conn结束后队列删除）
		false,        // 是否阻塞
		nil,          // 额外属性（我还不会用）
	)
	if err != nil {
		return nil, err
	}

	// 2.声明交换器
	err = mq.Channel.ExchangeDeclare(
		mq.Exchange, //交换器名
		"topic",     //exchange type：一般用fanout、direct、topic
		true,        // 是否持久化
		false,       // 是否自动删除（自动删除的前提是至少有一个队列或者交换器与这和交换器绑定，之后所有与这个交换器绑定的队列或者交换器都与此解绑）
		false,       // 设置是否内置的。true表示是内置的交换器，客户端程序无法直接发送消息到这个交换器中，只能通过交换器路由到交换器这种方式
		false,       // 是否阻塞
		nil,         // 额外属性
	)
	if err != nil {
		return nil, err
	}

	// 3.建立Binding(可随心所欲建立多个绑定关系)
	err = mq.Channel.QueueBind(
		mq.QueueName,  // 绑定的队列名称
		mq.RoutingKey, // bindkey 用于消息路由分发的key
		mq.Exchange,   // 绑定的exchange名
		false,         // 是否阻塞
		nil,           // 额外属性
	)
	if err != nil {
		return nil, err
	}

	return mq.Channel, nil
}

// SendMessage 发送消息(level 代表延迟级别)
func (mq *MyRabbitMQ) SendMessage(msg []byte) error {
	// 4.发送消息
	return mq.Channel.Publish(
		mq.Exchange,   // 交换器名
		mq.RoutingKey, // routing key
		false,         // 是否返回消息(匹配队列)，如果为true, 会根据binding规则匹配queue，如未匹配queue，则把发送的消息返回给发送者
		false,         // 是否返回消息(匹配消费者)，如果为true, 消息发送到queue后发现没有绑定消费者，则把发送的消息返回给发送者
		amqp.Publishing{ // 发送的消息，固定有消息体和一些额外的消息头，包中提供了封装对象
			ContentType: "text/plain", // 消息内容的类型
			Body:        msg,          // 消息内容
		},
	)
}

// 消费者
func (mq *MyRabbitMQ) NewConsumer(queueName string) (<-chan amqp.Delivery, error) {
	mq.QueueName = queueName
	// 1.声明队列（两端都要声明，原因在生产者处已经说明）
	_, err := mq.Channel.QueueDeclare( // 返回的队列对象内部记录了队列的一些信息，这里没什么用
		mq.QueueName, // 队列名
		true,         // 是否持久化
		false,        // 是否自动删除(前提是至少有一个消费者连接到这个队列，之后所有与这个队列连接的消费者都断开时，才会自动删除。注意：生产者客户端创建这个队列，或者没有消费者客户端与这个队列连接时，都不会自动删除这个队列)
		false,        // 是否为排他队列（排他的队列仅对“首次”声明的conn可见[一个conn中的其他channel也能访问该队列]，conn结束后队列删除）
		false,        // 是否阻塞
		nil,          // 额外属性（我还不会用）
	)
	if err != nil {
		fmt.Println("声明队列失败", err)
		return nil, err
	}

	// 2.从队列获取消息（消费者只关注队列）consume方式会不断的从队列中获取消息
	msgChanl, err := mq.Channel.Consume(
		mq.QueueName, // 队列名
		"",           // 消费者名，用来区分多个消费者，以实现公平分发或均等分发策略
		false,        // 是否自动应答
		false,        // 是否排他
		false,        // 是否接收只同一个连接中的消息，若为true，则只能接收别的conn中发送的消息
		false,        // 队列消费是否阻塞
		nil,          // 额外属性
	)
	if err != nil {
		fmt.Println("获取消息失败", err)
		return nil, err
	}

	return msgChanl, nil
}

func (m *MyRabbitMQ) Start(msgChanl <-chan amqp.Delivery, hand HandlerRocketMqRead) error { // 阻塞模式
	for msg := range msgChanl {
		// 这里写你的处理逻辑
		// 获取到的消息是amqp.Delivery对象，从中可以获取消息信息
		if hand(msg.Body) == nil {
			msg.Ack(true) // 主动应答
		}

	}

	return nil
}
