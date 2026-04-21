package router

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouter(t *testing.T) {
	r := setupRouter()

	tests := []struct {
		method   string
		path     string
		wantCode int
		wantBody string
	}{
		{"GET", "/health", 200, `{"status":"ok"}`},
		{"POST", "/notfound", 404, ""},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(tt.method, tt.path, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, tt.wantCode, w.Code)
		if tt.wantBody != "" {
			assert.JSONEq(t, tt.wantBody, w.Body.String())
		}
	}
}
