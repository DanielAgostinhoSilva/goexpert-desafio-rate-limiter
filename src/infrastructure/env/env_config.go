package env

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

type EnvConfig struct {
	WebServerPort        string `mapstructure:"WEB_SERVER_PORT"`
	MaxReqPerSecondToken int    `mapstructure:"MAX_REQUEST_PER_SECOND_BY_TOKEN"`
	MaxReqPerSecondIp    int    `mapstructure:"MAX_REQUEST_PER_SECOND_BY_IP"`
	BlockedTimePerSecond int    `mapstructure:"BLOCKED_TIME_PER_SECOND"`
	RedisAddress         string `mapstructure:"REDIS_ADDR"`
	RedisPassword        string `mapstructure:"REDIS_PASSWORD"`
	RedisDB              int    `mapstructure:"REDIS_DB"`
}

func LoadConfig(filePath string) *EnvConfig {
	var cfg *EnvConfig
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.SetConfigFile(filePath)
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	log.Println(fmt.Sprintf("Env initialized: %+v", cfg))
	return cfg
}
