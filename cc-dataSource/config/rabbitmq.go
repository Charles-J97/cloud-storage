package config

const (
	RabbitURL             = "amqp://guest:guest@127.0.0.1:5672/"
	HeartBeatExchangeName = "heartbeat.ex"
	HeartBeatQueueName    = "heartbeat.q"
	HeartBeatErrQueueName = "heartbeat.errq"
	HeartBeatRoutingKey   = "heartbeat"
)
