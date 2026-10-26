package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"my-todo/internal/adapter/handler"
	"my-todo/internal/domain"
	"my-todo/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
			name:       "limitが上限超過",
			query:      "?page=1&limit=101",
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

func (m *mockTodoUsecase) Get(ctx context.Context, id uuid.UUID) (*domain.Todo, error) {
	args := m.Called(ctx, id)

	var todo *domain.Todo
	if args.Get(0) != nil {
		todo = args.Get(0).(*domain.Todo)
	}
	return todo, args.Error(1)
}

func TestTodoHandler_Get(t *testing.T) {
	id := uuid.New()
	found := &domain.Todo{ID: id, Title: "買い物", Status: domain.StatusPending}

	tests := []struct {
		name       string
		paramID    string
		callUc     bool
		mockResult *domain.Todo
		mockErr    error
		wantErr    error
		wantStatus int
	}{
		{
			name:       "正常系",
			paramID:    id.String(),
			callUc:     true,
			mockResult: found,
			mockErr:    nil,
			wantErr:    nil,
			wantStatus: http.StatusOK,
		},
		{
			name:       "idが不正",
			paramID:    "abc",
			callUc:     false,
			mockResult: nil,
			mockErr:    nil,
			wantErr:    domain.ErrBadRequest,
			wantStatus: 0,
		},
		{
			name:       "not found",
			paramID:    id.String(),
			callUc:     true,
			mockResult: nil,
			mockErr:    domain.ErrNotFound,
			wantErr:    domain.ErrNotFound,
			wantStatus: 0,
		},
		{
			name:       "DB error",
			paramID:    id.String(),
			callUc:     true,
			mockResult: nil,
			mockErr:    errors.New("db error"),
			wantErr:    errors.New("db error"),
			wantStatus: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Request = httptest.NewRequest("GET", "/todos/"+tt.paramID, nil)
			c.Params = gin.Params{{Key: "id", Value: tt.paramID}}

			mockUc := new(mockTodoUsecase)
			if tt.callUc {
				mockUc.On("Get", mock.Anything, id).Return(tt.mockResult, tt.mockErr)
			}

			h := handler.NewTodoHandler(mockUc)
			h.Get(c)

			mockUc.AssertExpectations(t)
			if tt.wantErr != nil {
				assert.True(t, c.Errors.Last() != nil)
				if errors.Is(tt.wantErr, domain.ErrBadRequest) || errors.Is(tt.wantErr, domain.ErrNotFound) {
					assert.ErrorIs(t, c.Errors.Last().Err, tt.wantErr)
				}
				return
			}

			assert.Nil(t, c.Errors.Last())
			assert.Equal(t, tt.wantStatus, rec.Code)
			var body map[string]string
			assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
			assert.Equal(t, id.String(), body["id"])
			assert.Equal(t, "買い物", body["title"])
			assert.Equal(t, "pending", body["status"])
		})
	}
}

func (m *mockTodoUsecase) Update(ctx context.Context, id uuid.UUID, input usecase.UpdateInput) (*domain.Todo, error) {
	args := m.Called(ctx, id, input)

	var todo *domain.Todo
	if args.Get(0) != nil {
		todo = args.Get(0).(*domain.Todo)
	}
	return todo, args.Error(1)
}

func TestTodoHandler_Update(t *testing.T) {
	id := uuid.New()
	title := "掃除"
	status := "done"
	updated := &domain.Todo{ID: id, Title: title, Status: domain.StatusDone}

	tests := []struct {
		name       string
		paramID    string
		body       string
		callUc     bool
		input      usecase.UpdateInput
		mockResult *domain.Todo
		mockErr    error
		wantErr    error
		wantStatus int
	}{
		{
			name:       "正常系",
			paramID:    id.String(),
			body:       `{"title":"掃除","status":"done"}`,
			callUc:     true,
			input:      usecase.UpdateInput{Title: &title, Status: &status},
			mockResult: updated,
			wantStatus: http.StatusOK,
		},
		{
			name:    "idが不正",
			paramID: "abc",
			body:    `{"title":"掃除"}`,
			callUc:  false,
			wantErr: domain.ErrBadRequest,
		},
		{
			name:    "JSONが不正",
			paramID: id.String(),
			body:    `{"title":`,
			callUc:  false,
			wantErr: domain.ErrBadRequest,
		},
		{
			name:    "not found",
			paramID: id.String(),
			body:    `{"status":"done"}`,
			callUc:  true,
			input:   usecase.UpdateInput{Status: &status},
			mockErr: domain.ErrNotFound,
			wantErr: domain.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Request = httptest.NewRequest("PATCH", "/todos/"+tt.paramID, strings.NewReader(tt.body))
			c.Request.Header.Set("Content-Type", "application/json")
			c.Params = gin.Params{{Key: "id", Value: tt.paramID}}

			mockUc := new(mockTodoUsecase)
			if tt.callUc {
				mockUc.On("Update", mock.Anything, id, tt.input).Return(tt.mockResult, tt.mockErr)
			}

			h := handler.NewTodoHandler(mockUc)
			h.Update(c)

			mockUc.AssertExpectations(t)
			if tt.wantErr != nil {
				assert.True(t, c.Errors.Last() != nil)
				assert.ErrorIs(t, c.Errors.Last().Err, tt.wantErr)
				return
			}

			assert.Nil(t, c.Errors.Last())
			assert.Equal(t, tt.wantStatus, rec.Code)
			var body map[string]string
			assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
			assert.Equal(t, id.String(), body["id"])
			assert.Equal(t, "掃除", body["title"])
			assert.Equal(t, "done", body["status"])
		})
	}
}

func (m *mockTodoUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestTodoHandler_Delete(t *testing.T) {
	id := uuid.New()

	tests := []struct {
		name       string
		paramID    string
		callUc     bool
		mockErr    error
		wantErr    error
		wantStatus int
	}{
		{
			name:       "正常系",
			paramID:    id.String(),
			callUc:     true,
			wantStatus: http.StatusNoContent,
		},
		{
			name:    "idが不正",
			paramID: "abc",
			callUc:  false,
			wantErr: domain.ErrBadRequest,
		},
		{
			name:    "not found",
			paramID: id.String(),
			callUc:  true,
			mockErr: domain.ErrNotFound,
			wantErr: domain.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Request = httptest.NewRequest("DELETE", "/todos/"+tt.paramID, nil)
			c.Params = gin.Params{{Key: "id", Value: tt.paramID}}

			mockUc := new(mockTodoUsecase)
			if tt.callUc {
				mockUc.On("Delete", mock.Anything, id).Return(tt.mockErr)
			}

			h := handler.NewTodoHandler(mockUc)
			h.Delete(c)

			mockUc.AssertExpectations(t)
			if tt.wantErr != nil {
				assert.True(t, c.Errors.Last() != nil)
				assert.ErrorIs(t, c.Errors.Last().Err, tt.wantErr)
				return
			}

			assert.Nil(t, c.Errors.Last())
			assert.Equal(t, tt.wantStatus, c.Writer.Status())
			assert.Empty(t, rec.Body.String())
		})
	}
}
