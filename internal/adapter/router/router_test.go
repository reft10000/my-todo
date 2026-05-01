package router

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"my-todo/internal/adapter/handler"
	"my-todo/internal/domain"

	"github.com/stretchr/testify/assert"
)

type mockTodoUsecase struct{}

func (m *mockTodoUsecase) Create(title string) (*domain.Todo, error) {
	return domain.NewTodo(title)
}

func TestRouter(t *testing.T) {
	h := handler.NewTodoHandler(&mockTodoUsecase{})
	r := SetupRouter(h)

	tests := []struct {
		method      string
		path        string
		body        string
		contentType string
		wantCode    int
		wantBody    string
	}{
		{"GET", "/health", "", "", 200, `{"status":"ok"}`},
		{"POST", "/notfound", "", "", 404, ""},
		{"POST", "/todos", `{"title":"test todo"}`, "application/json", 201, ""},
	}

	for _, tt := range tests {
		var req *http.Request
		if tt.body != "" {
			req = httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", tt.contentType)
		} else {
			req = httptest.NewRequest(tt.method, tt.path, nil)
		}
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, tt.wantCode, w.Code)
		if tt.wantBody != "" {
			assert.JSONEq(t, tt.wantBody, w.Body.String())
		}
	}
}
