package conf

import (
	"rabbit-mq-consumer/model"
	"rabbit-mq-consumer/queue"

	"github.com/joho/godotenv"
)

// 初始化配置
func Init() {
	// 从本地读取环境变量
	_ = godotenv.Load()
	LoadConfig()
	model.InitConn(MiscLogsConfig.Url)
	queue.RabbitMq(RabbitMqConfig.Url)
}
