package entrypoint

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handle := NewRateLimiterHandle()
	handle.Handle(w, req)

	assert.Equal(t, w.Result().StatusCode, http.StatusOK)
	assert.Equal(t, w.Body.String(), "Enpoint acessado com sucesso")
}
