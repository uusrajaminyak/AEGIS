package config

import (
	"log"
	"github.com/spf13/viper"
)

type Config struct {
	ServerPort string `mapstructure:"SERVER_PORT"`
	Env		   string `mapstructure:"ENV"`
	DBHost	   string `mapstructure:"DB_HOST"`
	DBPort	   string `mapstructure:"DB_PORT"`
	DBUser	   string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName	   string `mapstructure:"DB_NAME"`
	ClickHouseHost string `mapstructure:"CLICKHOUSE_HOST"`
	ClickHousePort string `mapstructure:"CLICKHOUSE_PORT"`
	ClickHouseDB   string `mapstructure:"CLICKHOUSE_DB"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		log.Printf("Error reading config file, %s", err)
		return
	}

	err = viper.Unmarshal(&config)
	return
}