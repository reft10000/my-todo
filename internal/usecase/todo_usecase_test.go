package usecase_test

import (
	"errors"
	"testing"

	"my-todo/internal/domain"
	"my-todo/internal/usecase"

	"github.com/stretchr/testify/assert"
)

type mockTodoRepository struct {
	CreateFunc func(todo *domain.Todo) error
}

func (m *mockTodoRepository) Create(todo *domain.Todo) error {
	return m.CreateFunc(todo)
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
