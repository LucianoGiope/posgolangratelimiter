package app

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"

	"github.com/LucianoGiope/posgolangdesafioFinalRateLimiter/config/getenv"
)

type mockRouter struct {
	chi.Router
	useCalled  bool
	getCalled  bool
	getPattern string
}

func (m *mockRouter) Use(middleware ...func(http.Handler) http.Handler) {
	m.useCalled = true
	if len(middleware) > 0 {
		for _, m := range middleware {
			dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
			_ = m(dummyHandler)
		}
	}
}

func (m *mockRouter) Get(pattern string, handlerFn http.HandlerFunc) {
	m.getCalled = true
	m.getPattern = pattern
}

func TestStartDependencies_RegistersMiddlewareAndRoute(t *testing.T) {
	getenv.AppConfig = getenv.NewAppConfig("1s", "1s", "localhost:6379", "192.187.55.1", 10, 20)

	mockR := &mockRouter{}
	rdb := &redis.Client{}

	err := StartDependencies(mockR, rdb)
	assert.NoError(t, err)
	assert.True(t, mockR.useCalled)
	assert.True(t, mockR.getCalled)
	assert.Equal(t, "/", mockR.getPattern)
}
