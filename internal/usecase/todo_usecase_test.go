package usecase_test

import (
	"context"
	"errors"
	"testing"

	"my-todo/internal/domain"
	"my-todo/internal/usecase"

	"github.com/stretchr/testify/assert"
)

type mockTodoRepository struct {
	CreateFunc func(todo *domain.Todo) error
	ListFunc   func(ctx context.Context, params domain.ListParams) (domain.ListResult[domain.Todo], error)
}

func (m *mockTodoRepository) Create(todo *domain.Todo) error {
	return m.CreateFunc(todo)
}

func (m *mockTodoRepository) List(ctx context.Context, params domain.ListParams) (domain.ListResult[domain.Todo], error) {
	return m.ListFunc(ctx, params)
}

func TestTodoUsecase_Create(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		repoErr error
		wantErr bool
	}{
		{"正常系", "買い物", nil, false},
		{"タイトル空", "", nil, true},
		{"DB error", "買い物", errors.New("db error"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockTodoRepository{
				CreateFunc: func(todo *domain.Todo) error {
					return tt.repoErr
				},
			}
			uc := usecase.NewTodoUsecase(repo)
			todo, err := uc.Create(tt.title)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, todo)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, todo)
			}
		})
	}
}

func TestTodoUsecase_List(t *testing.T) {
	tests := []struct {
		name    string
		input   usecase.ListInput
		mock    domain.ListResult[domain.Todo]
		repoErr error
		wantErr bool
	}{
		{
			name:  "正常系",
			input: usecase.ListInput{Page: 1, Limit: 20},
			mock: domain.ListResult[domain.Todo]{
				Todos:   []domain.Todo{{Title: "買い物"}, {Title: "掃除"}},
				Total:   2,
				Page:    1,
				Limit:   20,
				HasNext: false,
			},
			repoErr: nil,
			wantErr: false,
		},
		{
			name:    "DB error",
			input:   usecase.ListInput{Page: 1, Limit: 20},
			mock:    domain.ListResult[domain.Todo]{},
			repoErr: errors.New("db error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockTodoRepository{
				ListFunc: func(ctx context.Context, params domain.ListParams) (domain.ListResult[domain.Todo], error) {
					return tt.mock, tt.repoErr
				},
			}
			ctx := context.Background()
			uc := usecase.NewTodoUsecase(repo)
			result, err := uc.List(ctx, tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mock.Total, result.Total)
				assert.Equal(t, len(tt.mock.Todos), len(result.Todos))
			}
		})
	}
}
