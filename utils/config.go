package utils

import (
	"time"

	"github.com/spf13/viper"
)

// Config stores all configs of the application.
// The values are read by viper from a config file or env variables
type Config struct {
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	DBSource            string        `mapstructure:"DB_SOURCE"`
	HTTPServerAddr          string        `mapstructure:"HTTP_SERVER_ADDR"`
	GRPCServerAddr          string        `mapstructure:"GRPC_SERVER_ADDR"`
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenLifeTime time.Duration `mapstructure:"ACCESS_TOKEN_LIFETIME"`
	RefreshTokenLifeTime time.Duration `mapstructure:"REFRESH_TOKEN_LIFETIME"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
