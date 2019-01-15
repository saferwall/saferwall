package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Load configration
func Load() {
	viper.SetConfigName("app")    // no need to include file extension
	viper.AddConfigPath("config") // set the path of your config file

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("failed to read config file :", err)
	}
}
