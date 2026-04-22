package handler_test

import (
	"net/http/httptest"
	"strings"
	"testing"

	"my-todo/internal/adapter/handler"
	"my-todo/internal/domain"

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
