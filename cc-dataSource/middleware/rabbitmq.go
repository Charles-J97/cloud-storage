package middleware

import (
	"cc-dataSource/config"
	"github.com/streadway/amqp"
	"log"
)

//注：类里面，私有属性用小写（不可被外部访问），公有属性用大写（可以被外部访问）
type RabbitMQ struct {
	conn       *amqp.Connection
	channel    *amqp.Channel
	QueueName  string
	Exchange   string
	RoutingKey string
	MqUrl      string
}

//创建rabbitMQ基础结构体实例
//要传入的都是公有属性，所以conn和channel都不需要传入，最后一个属性是已经定义好的，所以也不用传入
//故而要传入的只有QueueName，Exchange和RoutingKey三个属性
func NewRabbitMQ(exchange, queueName, routingKey string) *RabbitMQ {
	var err error
	rabbitmq := &RabbitMQ{
		QueueName:  queueName,
		Exchange:   exchange,
		RoutingKey: routingKey,
		MqUrl:      config.RabbitURL,
	}
	//创建rabbitmq连接
	//两个私有属性在这里会自己创建，不依赖于外部传值
	rabbitmq.conn, err = amqp.Dial(rabbitmq.MqUrl)
	if err != nil {
		log.Fatalf("Conection create failed:%s", err.Error())
	}
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	if err != nil {
		log.Fatalf("Channel create failed:%s", err.Error())
	}
	return rabbitmq
}

//手动断开rabbitMQ连接
func (rmq *RabbitMQ) Destory() {
	var err error
	err = rmq.channel.Close()
	if err != nil {
		log.Fatalf("Channel close failed:%s", err.Error())
	}
	err = rmq.conn.Close()
	if err != nil {
		log.Fatalf("Connection close failed:%s", err.Error())
	}
}

//发布消息
func (rmq *RabbitMQ) Publish(msg string) {
	//1.尝试建立Exchange
	err := rmq.channel.ExchangeDeclare(rmq.Exchange, "direct", false, false, false, false, nil)
	if err != nil {
		log.Fatalf(err.Error())
	}

	//2.发布消息
	err = rmq.channel.Publish(
		rmq.Exchange,
		rmq.RoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		})
	if err != nil {
		log.Fatalf(err.Error())
	}
}
