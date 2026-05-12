package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"my-todo/internal/adapter/handler"
	"my-todo/internal/domain"
	"my-todo/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockTodoUsecase struct {
	mock.Mock
}

func (m *mockTodoUsecase) Create(title string) (*domain.Todo, error) {
	args := m.Called(title)

	var todo *domain.Todo
	if args.Get(0) != nil {
		todo = args.Get(0).(*domain.Todo)
	}
	return todo, args.Error(1)
}

func TestTodoHandler_Create(t *testing.T) {
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)

	title := "invalid"
	body := `{"title":"` + title + `"}`
	c.Request = httptest.NewRequest("POST", "/todos", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	mockUc := new(mockTodoUsecase)
	mockUc.On("Create", title).Return(nil, domain.ErrInvalidTitle)

	h := handler.NewTodoHandler(mockUc)

	h.Create(c)

	assert.True(t, c.Errors.Last() != nil)
	assert.ErrorIs(t, c.Errors.Last().Err, domain.ErrInvalidTitle)
}

func (m *mockTodoUsecase) List(ctx context.Context, input usecase.ListInput) (*domain.ListResult[domain.Todo], error) {
	args := m.Called(ctx, input)

	var result *domain.ListResult[domain.Todo]
	if args.Get(0) != nil {
		result = args.Get(0).(*domain.ListResult[domain.Todo])
	}
	return result, args.Error(1)
}

func TestTodoHandler_List(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		mockResult *domain.ListResult[domain.Todo]
		mockErr    error
		wantErr    bool
		wantStatus int
	}{
		{
			name:  "正常系",
			query: "?page=1&limit=20",
			mockResult: &domain.ListResult[domain.Todo]{
				Todos:   []domain.Todo{{Title: "買い物"}, {Title: "掃除"}},
				Total:   2,
				Page:    1,
				Limit:   20,
				HasNext: false,
			},
			mockErr:    nil,
			wantErr:    false,
			wantStatus: http.StatusOK,
		},
		{
			name:       "pageが不正",
			query:      "?page=0&limit=20",
			mockResult: nil,
			mockErr:    nil,
			wantErr:    true,
			wantStatus: 0,
		},
		{
			name:       "limitが不正",
			query:      "?page=1&limit=0",
			mockResult: nil,
			mockErr:    nil,
			wantErr:    true,
			wantStatus: 0,
		},
		{
			name:       "DB error",
			query:      "?page=1&limit=20",
			mockResult: nil,
			mockErr:    errors.New("db error"),
			wantErr:    true,
			wantStatus: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Request = httptest.NewRequest("GET", "/todos"+tt.query, nil)

			mockUc := new(mockTodoUsecase)
			input := usecase.ListInput{Page: 1, Limit: 20}
			mockUc.On("List", mock.Anything, input).Return(tt.mockResult, tt.mockErr)

			h := handler.NewTodoHandler(mockUc)
			h.List(c)

			if tt.wantErr {
				assert.True(t, c.Errors.Last() != nil)
			} else {
				assert.Nil(t, c.Errors.Last())
				assert.Equal(t, tt.wantStatus, rec.Code)
			}
		})
	}
}
