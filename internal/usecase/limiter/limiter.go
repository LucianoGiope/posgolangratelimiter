package usecase

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/LucianoGiope/posgolangdesafioFinalRateLimiter/internal/domain"
	"github.com/LucianoGiope/posgolangdesafioFinalRateLimiter/internal/repository"
	"github.com/LucianoGiope/posgolangdesafioFinalRateLimiter/internal/statics"
)

type RateLimiterExecutor interface {
	Execute(w http.ResponseWriter, r *http.Request) (bool, error)
}

type RateLimiterUseCase struct {
	Repo        repository.IRateLimiterRepository
	RateLimiter *domain.RateLimiter
}

func NewRateLimiterUseCase(repo repository.IRateLimiterRepository, rl *domain.RateLimiter) *RateLimiterUseCase {
	return &RateLimiterUseCase{
		Repo:        repo,
		RateLimiter: rl,
	}
}

func (useCase *RateLimiterUseCase) Execute(w http.ResponseWriter, r *http.Request) (bool, error) {
	token := r.Header.Get(statics.API_KEY)
	if token != "" {
		isAllowedToken, err := useCase.Repo.AllowToken(
			useCase.RateLimiter.Context,
			token,
			useCase.RateLimiter.RateLimiterByToken,
			useCase.RateLimiter.RateLimiterTimeByToken)

		if err != nil {
			return false, fmt.Errorf("Falha ao validar o token: %v", err)
		}

		return isAllowedToken, nil
	}

	ip := useCase.RateLimiter.IpFakeToTester
	if strings.TrimSpace(ip) == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}
	if ip == "" {

		var err error
		ip, _, err = net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return false, fmt.Errorf("Falha ao identificar e ao validar o ip: %v", err)
		}
	}

	isAllowedIP, err := useCase.Repo.AllowIP(
		useCase.RateLimiter.Context,
		ip,
		useCase.RateLimiter.RateLimiterByIP,
		useCase.RateLimiter.RateLimiterTimeByIP)

	if err != nil {
		return false, fmt.Errorf("Falha ao validar o ip: %v, erro:%v", ip, err)
	}

	return isAllowedIP, nil
}
