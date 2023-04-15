package config

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var values *Values

type Values struct {
	Keystore   string
	Passphrase string
	Port       int
	GinMode    string
	RPC        string
	Contract   string
}

func InitValues() error {
	viper.SetDefault("port", 8888)
	viper.SetDefault("gin_mode", gin.ReleaseMode)

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	_ = viper.BindEnv("PORT")
	_ = viper.BindEnv("GIN_MODE")
	_ = viper.BindEnv("KEYSTORE")
	_ = viper.BindEnv("PASSPHARSE")

	values = &Values{
		Keystore:   viper.GetString("KEYSTORE"),
		Passphrase: viper.GetString("PASSPHARSE"),
		Port:       viper.GetInt("PORT"),
		GinMode:    viper.GetString("GIN_MODE"),
		RPC:        viper.GetString("RPC"),
		Contract:   viper.GetString("CONTRACT"),
	}
	return nil
}

func Config() *Values {
	if values == nil {
		log.Fatal("config not initial")
	}
	return values
}
