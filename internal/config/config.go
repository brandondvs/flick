package config

import (
	"github.com/spf13/viper"
)

func Load(configFilePath string) error {
	viper.SetConfigFile(configFilePath)
	return viper.ReadInConfig()
}

func ServerHost() string {
	return viper.GetString("server.host")
}

func ServerPort() int {
	return viper.GetInt("server.port")
}
