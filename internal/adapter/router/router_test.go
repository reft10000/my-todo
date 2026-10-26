package router

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"my-todo/internal/adapter/handler"
	"my-todo/internal/domain"
	"my-todo/internal/usecase"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockTodoUsecase struct{}

func (m *mockTodoUsecase) Create(title string) (*domain.Todo, error) {
	return domain.NewTodo(title)
}

func (m *mockTodoUsecase) List(ctx context.Context, input usecase.ListInput) (*domain.ListResult[domain.Todo], error) {
	return &domain.ListResult[domain.Todo]{
		Todos:   []domain.Todo{},
		Total:   0,
		Page:    input.Page,
		Limit:   input.Limit,
		HasNext: false,
	}, nil
}

func (m *mockTodoUsecase) Get(ctx context.Context, id uuid.UUID) (*domain.Todo, error) {
	return &domain.Todo{ID: id, Title: "test todo", Status: domain.StatusPending}, nil
}

func (m *mockTodoUsecase) Update(ctx context.Context, id uuid.UUID, input usecase.UpdateInput) (*domain.Todo, error) {
	todo := &domain.Todo{ID: id, Title: "test todo", Status: domain.StatusPending}
	if input.Status != nil {
		todo.Status = domain.Status(*input.Status)
	}
	return todo, nil
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
		{"GET", "/todos?page=1&limit=20", "", "", 200, ""},
		{"GET", "/todos/0195f3a0-0000-7000-8000-000000000000", "", "", 200, ""},
		{"GET", "/todos/abc", "", "", 400, `{"error":"bad request"}`},
		{"PATCH", "/todos/0195f3a0-0000-7000-8000-000000000000", `{"status":"done"}`, "application/json", 200, `{"id":"0195f3a0-0000-7000-8000-000000000000","title":"test todo","status":"done"}`},
		{"PATCH", "/todos/abc", `{"status":"done"}`, "application/json", 400, `{"error":"bad request"}`},
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
