package middleware

import (
	"net/http"

	usecase "github.com/LucianoGiope/posgolangdesafioFinalRateLimiter/internal/usecase/limiter"
)

func RateLimiterMiddleware(useCase usecase.RateLimiterExecutor, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allow, err := useCase.Execute(w, r)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !allow {
			http.Error(w, "Número máximo por tempo atingido.", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
