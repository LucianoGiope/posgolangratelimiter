package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockRateLimiterUseCase struct {
	allow  bool
	err    error
	called bool
}

func (m *mockRateLimiterUseCase) Execute(w http.ResponseWriter, r *http.Request) (bool, error) {
	m.called = true
	return m.allow, m.err
}

func TestRateLimiterMiddleware_Allow(t *testing.T) {
	mockUseCase := &mockRateLimiterUseCase{allow: true, err: nil}
	handlerCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	})

	middleware := RateLimiterMiddleware(mockUseCase, next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	middleware.ServeHTTP(rr, req)

	assert.True(t, handlerCalled)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.True(t, mockUseCase.called)
}

func TestRateLimiterMiddleware_LimitReached(t *testing.T) {
	mockUseCase := &mockRateLimiterUseCase{allow: false, err: nil}
	handlerCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	})

	middleware := RateLimiterMiddleware(mockUseCase, next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	middleware.ServeHTTP(rr, req)

	assert.False(t, handlerCalled)
	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	assert.Contains(t, rr.Body.String(), "Número máximo e tentativas atingido.")
	assert.True(t, mockUseCase.called)
}

func TestRateLimiterMiddleware_Error(t *testing.T) {
	errMsg := "internal error"
	mockUseCase := &mockRateLimiterUseCase{allow: false, err: errors.New(errMsg)}
	handlerCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	})

	middleware := RateLimiterMiddleware(mockUseCase, next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	middleware.ServeHTTP(rr, req)

	assert.False(t, handlerCalled)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), errMsg)
	assert.True(t, mockUseCase.called)
}
