package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"my-todo/internal/adapter/middleware"
	"my-todo/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestErrorHandler(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantCode int
	}{
		{"AppError 400", domain.ErrBadRequest, 400},
		{"unknown error", errors.New("unknown"), 500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.Use(middleware.ErrorHandler())
			r.GET("/test", func(c *gin.Context) {
				c.Error(tt.err)
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)
		})
	}
}
