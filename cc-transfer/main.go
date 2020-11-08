package main

import (
	"cc-transfer/config"
	util "cc-transfer/middleware"
	"cc-transfer/models/dto"
	"cc-transfer/route"
	"log"
	"os"
	"time"
)

var DataSource = make(map[string]time.Time)

func main() {
	go util.ListenHeartbeat(&DataSource)
	go util.AddServerToHashCircle(DataSource, dto.GlobalHashCircle)
	router := route.Router()
	err := router.Run(config.TransferServiceHost)
	if err != nil {
		log.Print(err.Error())
		os.Exit(1)
	}
}

