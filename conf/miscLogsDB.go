package conf

import "github.com/spf13/viper"

type MiscLogsDB struct {
	Url string
}

func InitMiscLogsDB(cfg *viper.Viper) *MiscLogsDB {
	return &MiscLogsDB{
		Url: cfg.GetString("url"),
	}
}

var MiscLogsConfig = new(MiscLogsDB)

type VuemobileDB struct {
	Url string
}

func InitVuemobileDB(cfg *viper.Viper) *VuemobileDB {
	return &VuemobileDB{
		Url: cfg.GetString("url"),
	}
}

var VuemobileConfig = new(VuemobileDB)
