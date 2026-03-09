package config

import (
	"fmt"

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

func DatabaseConnectionString() string {
	host := viper.GetString("database.host")
	port := viper.GetInt("database.port")
	user := viper.GetString("database.user")
	password := viper.GetString("database.password")
	databaseName := viper.GetString("database.name")

	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", user, password, host, port, databaseName)
}
