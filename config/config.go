package config

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var values *Values

type Values struct {
	// database
	DbHost     string
	DbPort     uint
	DbUser     string
	DbName     string
	DbPassword string

	Keystore   string
	Passphrase string
	Port       int
	GinMode    string
	RPC        string
	Contract   string
	MaxGas     string
}

func InitValues() error {
	viper.SetDefault("port", 8888)
	viper.SetDefault("gin_mode", gin.ReleaseMode)
	viper.SetDefault("MAX_GAS", "1000000000000000000")

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	_ = viper.BindEnv("DB_HOST")
	_ = viper.BindEnv("DB_PORT")
	_ = viper.BindEnv("DB_USER")
	_ = viper.BindEnv("DB_NAME")
	_ = viper.BindEnv("DB_PASSWORD")
	_ = viper.BindEnv("PORT")
	_ = viper.BindEnv("GIN_MODE")
	_ = viper.BindEnv("KEYSTORE")
	_ = viper.BindEnv("PASSPHARSE")
	_ = viper.BindEnv("RPC")
	_ = viper.BindEnv("CONTRACT")
	_ = viper.BindEnv("MAX_GAS")

	values = &Values{
		DbHost:     viper.GetString("DB_HOST"),
		DbPort:     viper.GetUint("DB_PORT"),
		DbUser:     viper.GetString("DB_USER"),
		DbName:     viper.GetString("DB_NAME"),
		DbPassword: viper.GetString("DB_PASSWORD"),
		Keystore:   viper.GetString("KEYSTORE"),
		Passphrase: viper.GetString("PASSPHARSE"),
		Port:       viper.GetInt("PORT"),
		GinMode:    viper.GetString("GIN_MODE"),
		RPC:        viper.GetString("RPC"),
		Contract:   viper.GetString("CONTRACT"),
		MaxGas:     viper.GetString("MAX_GAS"),
	}
	return nil
}

func Config() *Values {
	if values == nil {
		log.Fatal("config not initial")
	}
	return values
}
