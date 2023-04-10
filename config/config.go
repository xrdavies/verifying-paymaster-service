package config

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Values struct {
	Keystore   string
	Passphrase string
	Port       int
	GinMode    string
}

func GetValues() *Values {
	viper.SetDefault("port", 8888)
	viper.SetDefault("gin_mode", gin.ReleaseMode)

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
	}

	_ = viper.BindEnv("PORT")
	_ = viper.BindEnv("GIN_MODE")
	_ = viper.BindEnv("KEYSTORE")
	_ = viper.BindEnv("PASSPHARSE")

	return &Values{
		Keystore:   viper.GetString("KEYSTORE"),
		Passphrase: viper.GetString("PASSPHARSE"),
		Port:       viper.GetInt("PORT"),
		GinMode:    viper.GetString("GIN_MODE"),
	}
}
