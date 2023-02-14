package conf

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

var path = "./settings.yml"
var cfgMiscLogsDB, cfgRabbitMq *viper.Viper

func LoadConfig() {
	viper.SetConfigFile(path)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Read config file fail: %s", err.Error())
	}
	//Replace environment variables
	err = viper.ReadConfig(strings.NewReader(os.ExpandEnv(string(content))))
	if err != nil {
		log.Fatalf("Parse config file fail: %s", err.Error())
	}

	cfgMiscLogsDB = viper.Sub("settings.miscLogsDB")
	if cfgMiscLogsDB == nil {
		panic("No found settings.miscLogsDB in the configuration")
	}
	MiscLogsConfig = InitMiscLogsDB(cfgMiscLogsDB)

	cfgRabbitMq = viper.Sub("settings.rabbitmq")
	if cfgRabbitMq == nil {
		panic("No found settings.rabbitmq in the configuration")
	}
	RabbitMqConfig = InitRabbitMq(cfgRabbitMq)
}
