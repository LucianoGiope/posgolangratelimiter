package getenv

import (
	"github.com/spf13/viper"
)

type appConfig struct {
	RateLimiterByIP        int    `mapstructure:"RATE_LIMITER_IP_MAX_REQUESTS"`
	RateLimiterTimeByIP    string `mapstructure:"RATE_LIMITER_IP_BLOCK_TIME"`
	RateLimiterByToken     int    `mapstructure:"RATE_LIMITER_TOKEN_MAX_REQUESTS"`
	RateLimiterTimeByToken string `mapstructure:"RATE_LIMITER_TOKEN_BLOCK_TIME"`
	IpFakeToTester         string `mapstructure:"IP_FAKE_TO_TESTER"`
	RedisAddr              string `mapstructure:"REDIS_ADDR"`
	ServerPortDefault      string `mapstructure:"SERVER_PORT_DEFAULT"`
}

var AppConfig *appConfig

func NewAppConfig(timeIp, timeToken, redisAddr, ipFakeToTester string, limitIp, limitToken int) *appConfig {
	return &appConfig{
		RateLimiterByIP:        limitIp,
		RateLimiterTimeByIP:    timeIp,
		RateLimiterByToken:     limitToken,
		RateLimiterTimeByToken: timeToken,
		IpFakeToTester:         ipFakeToTester,
		RedisAddr:              redisAddr,
	}
}

func NewConfig(paths ...string) error {
	config := &appConfig{}

	var lastErr error
	for _, path := range paths {
		viper.SetConfigFile(path)
		viper.SetConfigType("env")
		if err := viper.ReadInConfig(); err == nil {
			viper.AutomaticEnv()
			if err := viper.Unmarshal(&config); err != nil {
				return err
			}
			AppConfig = config
			return nil
		} else {
			lastErr = err
		}
	}
	return lastErr
}
