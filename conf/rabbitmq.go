package conf

import "github.com/spf13/viper"

type RabbitMq struct {
	Url string
}

func InitRabbitMq(cfg *viper.Viper) *RabbitMq {
	return &RabbitMq{
		Url: cfg.GetString("url"),
	}
}

var RabbitMqConfig = new(RabbitMq)
