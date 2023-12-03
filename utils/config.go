package utils

import (
	"time"

	"github.com/spf13/viper"
)

// Config stores all configs of the application.
// The values are read by viper from a config file or env variables
type Config struct {
	Environment          string        `mapstructure:"ENVIRONMENT"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	MigrationURL         string        `mapstructure:"MIGRATION_URL"`
	HTTPServerAddr       string        `mapstructure:"HTTP_SERVER_ADDR"`
	GRPCServerAddr       string        `mapstructure:"GRPC_SERVER_ADDR"`
	RedisAddr            string        `mapstructure:"REDIS_ADDR"`
	EmailSenderName      string        `mapstructure:"EMAIL_SENDER_NAME"`
	EmailSenderAddr      string        `mapstructure:"EMAIL_SENDER_ADDR"`
	EmailSenderPassword  string        `mapstructure:"EMAIL_SENDER_PASSWORD"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenLifeTime  time.Duration `mapstructure:"ACCESS_TOKEN_LIFETIME"`
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
