package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/LucianoGiope/posgolangdesafioFinalRateLimiter/config/app"
	"github.com/LucianoGiope/posgolangdesafioFinalRateLimiter/config/getenv"
	"github.com/LucianoGiope/posgolangdesafioFinalRateLimiter/config/redisdb"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

// Constantes de teste
const (
	IpDeTeste1           = "192.168.1.55:45896"
	IpDeTeste2           = "4.4.4.4:5555"
	IpDeTeste1Isolado    = "192.168.13.66:7251"
	TestTokenIntegration = "token-987-integration"
	TestTokenIsolado1    = "token-1-isolado"
	TestTokenIsolado2    = "token-2-isolado"
)

// Limpa as chaves de teste do Redis
func cleanupTestKeys(rdb *redis.Client) {
	ctx := context.Background()
	ipKeys := []string{
		"192.168.1.55:25666",
		"4.4.4.4:8947",
		"192.168.13.66:6584",
	}
	tokenKeys := []string{
		"token-987-integration",
		"token-1-isolado",
		"token-2-isolado",
	}
	for _, key := range ipKeys {
		rdb.Del(ctx, key)
	}
	for _, key := range tokenKeys {
		rdb.Del(ctx, key)
	}
}

func setupTestServer(t *testing.T) (*httptest.Server, *redis.Client) {
	err := getenv.NewConfig("../../.env")
	if err != nil {
		assert.Fail(t, fmt.Sprintf("Falha ao ler o arquivo de configuração: %v", err.Error()))
	}
	if getenv.AppConfig == nil {
		assert.Fail(t, "Erro ao inicializar dados padrão.")
	}

	router := chi.NewRouter()
	rdb := redisdb.NewRedisClient(getenv.AppConfig.RedisAddr)
	err = app.StartDependencies(router, rdb)
	if err != nil {
		t.Fatalf("Falha ao iniciar o sistema: %v", err)
	}
	return httptest.NewServer(router), rdb
}

func TestRateLimiter_IP(t *testing.T) {
	ts, rdb := setupTestServer(t)
	defer ts.Close()
	defer rdb.Close()
	defer cleanupTestKeys(rdb)

	client := &http.Client{}
	maxRequests := getenv.AppConfig.RateLimiterByIP
	blockTime, _ := time.ParseDuration(getenv.AppConfig.RateLimiterTimeByIP)

	for i := 0; i < maxRequests; i++ {
		req, _ := http.NewRequest("GET", ts.URL+"/", nil)
		req.RemoteAddr = IpDeTeste1
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()
	}

	// Exceeding the limit
	req, _ := http.NewRequest("GET", ts.URL+"/", nil)
	req.RemoteAddr = IpDeTeste1
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(body), "número máximo de tentativas")
	resp.Body.Close()

	// After block time
	t.Logf("Aguarando tempo de %v para o bloqueio...", blockTime)
	time.Sleep(blockTime + time.Second)

	// should allow again
	req, _ = http.NewRequest("GET", ts.URL+"/", nil)
	req.RemoteAddr = IpDeTeste1
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()
}

func TestRateLimiter_Token(t *testing.T) {
	ts, rdb := setupTestServer(t)
	defer ts.Close()
	defer rdb.Close()
	defer cleanupTestKeys(rdb)

	client := &http.Client{}
	maxRequests := getenv.AppConfig.RateLimiterByToken
	blockTime, _ := time.ParseDuration(getenv.AppConfig.RateLimiterTimeByToken)
	token := TestTokenIntegration

	for range maxRequests {
		req, _ := http.NewRequest("GET", ts.URL+"/", nil)
		req.Header.Set("API_KEY", token)
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()
	}

	// Exceeding the limit
	req, _ := http.NewRequest("GET", ts.URL+"/", nil)
	req.Header.Set("API_KEY", token)
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(body), "número máximo de tentativas")
	resp.Body.Close()

	// After block time
	t.Logf("Aguarando tempo de %v para o bloqueio...", blockTime)
	time.Sleep(blockTime + time.Second)

	// should allow again
	req, _ = http.NewRequest("GET", ts.URL+"/", nil)
	req.Header.Set("API_KEY", token)
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()
}

func TestRateLimiter_Isolado(t *testing.T) {
	ts, rdb := setupTestServer(t)
	defer ts.Close()
	defer rdb.Close()
	defer cleanupTestKeys(rdb)

	client := &http.Client{}

	// IP1 hits the limit
	maxRequests := getenv.AppConfig.RateLimiterByIP
	for range maxRequests {
		req, _ := http.NewRequest("GET", ts.URL+"/", nil)
		req.Header.Set("X-Forwarded-For", IpDeTeste1Isolado)
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()
	}

	// IP1 blocked
	req, _ := http.NewRequest("GET", ts.URL+"/", nil)
	req.Header.Set("X-Forwarded-For", IpDeTeste1Isolado)
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	resp.Body.Close()

	// IP2 should pass normally
	req, _ = http.NewRequest("GET", ts.URL+"/", nil)
	req.Header.Set("X-Forwarded-For", IpDeTeste2)
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// Token1 hits the limit
	maxToken := getenv.AppConfig.RateLimiterByToken
	token1 := TestTokenIsolado1
	token2 := TestTokenIsolado2

	for range maxToken {
		req, _ := http.NewRequest("GET", ts.URL+"/", nil)
		req.Header.Set("API_KEY", token1)
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()
	}

	// Token1 blocked
	req, _ = http.NewRequest("GET", ts.URL+"/", nil)
	req.Header.Set("API_KEY", token1)
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	resp.Body.Close()

	// Token2 should pass normalmente
	req, _ = http.NewRequest("GET", ts.URL+"/", nil)
	req.Header.Set("API_KEY", token2)
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()
}
