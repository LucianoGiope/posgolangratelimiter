package getenv

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	if err := NewConfig("../../.env"); err != nil {
		assert.Fail(t, fmt.Sprintf("failure to init env config for test purpose: %v", err.Error()))
		return
	}

	assert.GreaterOrEqual(t, AppConfig.RateLimiterByIP, 1)
	assert.GreaterOrEqual(t, AppConfig.RateLimiterByToken, 1)
	assert.NotEmpty(t, AppConfig.RateLimiterTimeByIP)
	assert.NotEmpty(t, AppConfig.RateLimiterTimeByToken)
	assert.NotEmpty(t, AppConfig.RedisAddr)
}

func TestNewConfig_FilePathError(t *testing.T) {
	err := NewConfig("some-invalid-path")
	assert.Error(t, err)
}

func TestNewAppConfig(t *testing.T) {
	limitIp := 10
	timeIp := "10m"

	limitToken := 20
	timeToken := "5m"

	redisAddr := "localhost:6379"
	ipFakeToTester := "192.187.55.1"

	cfg := NewAppConfig(timeIp, timeToken, redisAddr, ipFakeToTester, limitIp, limitToken)

	assert.Equal(t, limitIp, cfg.RateLimiterByIP)
	assert.Equal(t, limitToken, cfg.RateLimiterByToken)
	assert.Equal(t, timeIp, cfg.RateLimiterTimeByIP)
	assert.Equal(t, timeToken, cfg.RateLimiterTimeByToken)
	assert.Equal(t, redisAddr, cfg.RedisAddr)
}

func TestNewConfig_UnmarshallError(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "invalidenv")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString("RATE_LIMITER_IP_MAX_REQUESTS=notanumber\n")
	assert.NoError(t, err)
	tmpFile.Close()

	err = NewConfig(tmpFile.Name())
	assert.Error(t, err)
}
