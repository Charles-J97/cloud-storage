package middleware

import (
	"cc-transfer/config"
	"github.com/streadway/amqp"
	"log"
	"time"
)

//注：类里面，私有属性用小写（不可被外部访问），公有属性用大写（可以被外部访问）
type RabbitMQ struct {
	conn *amqp.Connection
	channel *amqp.Channel
	QueueName string
	Exchange string
	RoutingKey string
	MqUrl string
}

//创建rabbitMQ基础结构体实例
//要传入的都是公有属性，所以conn和channel都不需要传入，最后一个属性是已经定义好的，所以也不用传入
//故而要传入的只有QueueName，Exchange和RoutingKey三个属性
func NewRabbitMQ(exchange, queueName, routingKey string) *RabbitMQ {
	var err error
	rabbitmq := &RabbitMQ{
		QueueName: queueName,
		Exchange: exchange,
		RoutingKey: routingKey,
		MqUrl: config.RabbitURL,
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

//消费消息
func (rmq *RabbitMQ) Consume(dataSource *map[string]time.Time) {
	//1.尝试建立Exchange
	err := rmq.channel.ExchangeDeclare(rmq.Exchange,"direct",false,false,false,false,nil)
	if err != nil {
		log.Fatalf(err.Error())
	}

	//2.尝试建立Queue
	_, err = rmq.channel.QueueDeclare(rmq.QueueName, true, false, false, false, nil)
	if err != nil {
		log.Fatalf(err.Error())
	}

	//3.绑定Queue到Exchange上
	err = rmq.channel.QueueBind(rmq.QueueName, rmq.RoutingKey, rmq.Exchange, false, nil)
	if err != nil {
		log.Fatalf(err.Error())
	}

	//4.接收所有消息
	msgs, err := rmq.channel.Consume(rmq.QueueName, "",false,false,false,false,nil)
	if err != nil {
		log.Fatalf(err.Error())
	}

	//5.处理消息
	//forever用于阻塞进程
	forever := make(chan bool)
	go func() {
		for msg := range msgs {
			serverInfo := string(msg.Body)
			if err != nil {
				log.Fatal(err.Error())
				return
			}
			log.Println(serverInfo + " has been added to DataSource list")
			Mutex.Lock()
			(*dataSource)[serverInfo] = time.Now()
			Mutex.Unlock()
		}
	}()
	log.Printf("[*] Waiting for messages. To exit, press CTRL+C")
	<-forever
}