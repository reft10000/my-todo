package middleware_test

import (
	"errors"
	"fmt"
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
		{"AppErrorのCodeを使う", &domain.AppError{Code: http.StatusNotFound, Message: "not found"}, 404},
		{"wrapされたAppError", fmt.Errorf("wrap: %w", domain.ErrInvalidTitle), 400},
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

func TestErrorHandler_DoesNotMutateDomainErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.ErrorHandler())
	r.GET("/test", func(c *gin.Context) {
		c.Error(domain.ErrBadRequest)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(httptest.NewRecorder(), req)

	assert.Equal(t, "invalid title", domain.ErrInvalidTitle.Message)
}
