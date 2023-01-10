package config

import (
	"time"

	"github.com/spf13/viper"
)

var Instance = &config{}

type config struct {
	Environment          string
	DatabaseURL          string
	TokenSymmetricKey    string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

func init() {
	viper.AutomaticEnv()

	Instance.Environment = viper.GetString("ENVIRONMENT")
	Instance.DatabaseURL = viper.GetString("DATABASE_URL")
	Instance.TokenSymmetricKey = viper.GetString("TOKEN_SYMMETRIC_KEY")
	Instance.AccessTokenDuration = viper.GetDuration("ACCESS_TOKEN_DURATION")
}
