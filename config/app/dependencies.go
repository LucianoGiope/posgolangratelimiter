package app

import (
	"net/http"

	"github.com/LucianoGiope/posgolangdesafioFinalRateLimiter/config/getenv"
	"github.com/LucianoGiope/posgolangdesafioFinalRateLimiter/internal/domain"
	"github.com/LucianoGiope/posgolangdesafioFinalRateLimiter/internal/entrypoint"
	"github.com/LucianoGiope/posgolangdesafioFinalRateLimiter/internal/middleware"
	"github.com/LucianoGiope/posgolangdesafioFinalRateLimiter/internal/repository"
	usecase "github.com/LucianoGiope/posgolangdesafioFinalRateLimiter/internal/usecase/limiter"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

func StartDependencies(router chi.Router, rdb *redis.Client) error {
	config := getenv.AppConfig
	rl := domain.NewRateLimiter(
		config.RateLimiterByIP,
		config.RateLimiterTimeByIP,
		config.RateLimiterByToken,
		config.RateLimiterTimeByToken,
		config.IpFakeToTester,
	)

	redisRepo := repository.NewRedisRepository(rdb)
	limiterUseCase := usecase.NewRateLimiterUseCase(redisRepo, rl)
	limiterHandle := entrypoint.NewRateLimiterHandle()

	router.Use(func(next http.Handler) http.Handler {
		return middleware.RateLimiterMiddleware(limiterUseCase, next)
	})
	router.Get("/", limiterHandle.Handle)

	return nil
}
