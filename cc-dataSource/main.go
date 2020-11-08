package main

import (
	"cc-dataSource/config"
	util "cc-dataSource/middleware"
	"cc-dataSource/route"
	_ "github.com/gin-gonic/gin"
	"log"
	"os"
	"time"
)

func sendHeartBeat() {
	rabbitmq := util.NewRabbitMQ(config.HeartBeatExchangeName, config.HeartBeatQueueName, config.HeartBeatRoutingKey)
	for {
		rabbitmq.Publish(config.DataSourceServiceHost)
		time.Sleep(3 * time.Second)
	}
}

func main() {
	go sendHeartBeat()
	router := route.Router()
	err := router.Run(config.DataSourceServiceHost)
	if err != nil {
		log.Print(err.Error())
		os.Exit(1)
	}
}
