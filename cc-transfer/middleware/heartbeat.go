package middleware

import (
	"cc-transfer/config"
	"cc-transfer/models/dto"
	"log"
	"sync"
	"time"
)

var Mutex sync.Mutex

//这个函数用于监听rabbitMQ并对从其内获得的消息进行处理
func ListenHeartbeat(dataSource *map[string]time.Time) {
	go removeExpiredDataSource(dataSource)
	rabbitmq := NewRabbitMQ(config.HeartBeatExchangeName, config.HeartBeatQueueName, config.HeartBeatRoutingKey)
	rabbitmq.Consume(dataSource)
}

//每隔5秒检测一次dataSource列表，如果发现有休眠时间超过10s的，就将其删除
func removeExpiredDataSource(dataSource *map[string]time.Time) {
	for {
		time.Sleep(5 * time.Second)
		Mutex.Lock()
		for k, v := range *dataSource {
			if v.Add(10 * time.Second).Before(time.Now()) {
				delete(*dataSource, k)
			}
		}
		Mutex.Unlock()
	}
}

//每3秒检查一次DataSource列表，将其中未加入Hash环上的服务器地址加入
func AddServerToHashCircle(serverMap map[string]time.Time, hashCircle *dto.Consistent) {
	for {
		time.Sleep(3 * time.Second)
		for k := range serverMap {
			if _, ok := hashCircle.Circle[hashCircle.GetHashKey(k, 0)]; !ok {
				hashCircle.Add(k)
				log.Println(k + " has been added to hash circle")
			}
		}
	}
}