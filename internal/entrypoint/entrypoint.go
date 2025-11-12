package entrypoint

import (
	"fmt"
	"net/http"
)

type RateLimiterHandle struct{}

func NewRateLimiterHandle() *RateLimiterHandle {
	return &RateLimiterHandle{}
}

func (handle *RateLimiterHandle) Handle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Enpoint acessado com sucesso")
}
