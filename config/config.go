package config

import (
	"log"

	"github.com/spf13/viper"
)

var AppConfig Config

type Config struct {
	Port        int
	Environment string
	Debug       bool

	DBHost     string
	DBPort     int
	DBDatabase string
	DBUsername string
	DBPassword string

	JWTSecret string
}

func InitializeAppConfig() {
	viper.AutomaticEnv()

	AppConfig.Port = viper.GetInt("PORT")
	AppConfig.Environment = viper.GetString("ENVIRONMENT")
	AppConfig.Debug = viper.GetBool("DEBUG")

	AppConfig.DBHost = viper.GetString("MYSQL_HOST")
	AppConfig.DBPort = viper.GetInt("MYSQL_PORT")
	AppConfig.DBDatabase = viper.GetString("MYSQL_DATABASE")
	AppConfig.DBUsername = viper.GetString("MYSQL_USER")
	AppConfig.DBPassword = viper.GetString("MYSQL_PASSWORD")

	AppConfig.JWTSecret = viper.GetString("JWT_SECRET")

	log.Println("[INIT] configuration loaded")
}
