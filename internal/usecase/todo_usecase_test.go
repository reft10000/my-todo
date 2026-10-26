package usecase_test

import (
	"context"
	"errors"
	"testing"

	"my-todo/internal/domain"
	"my-todo/internal/usecase"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockTodoRepository struct {
	CreateFunc   func(todo *domain.Todo) error
	ListFunc     func(ctx context.Context, params domain.ListParams) (domain.ListResult[domain.Todo], error)
	FindByIDFunc func(ctx context.Context, id uuid.UUID) (*domain.Todo, error)
	UpdateFunc   func(ctx context.Context, todo *domain.Todo) error
}

func (m *mockTodoRepository) Update(ctx context.Context, todo *domain.Todo) error {
	return m.UpdateFunc(ctx, todo)
}

func (m *mockTodoRepository) Create(todo *domain.Todo) error {
	return m.CreateFunc(todo)
}

func (m *mockTodoRepository) List(ctx context.Context, params domain.ListParams) (domain.ListResult[domain.Todo], error) {
	return m.ListFunc(ctx, params)
}

func (m *mockTodoRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Todo, error) {
	return m.FindByIDFunc(ctx, id)
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

func TestTodoUsecase_Get(t *testing.T) {
	id := uuid.New()
	found := &domain.Todo{ID: id, Title: "買い物", Status: domain.StatusPending}

	tests := []struct {
		name     string
		mock     *domain.Todo
		repoErr  error
		wantErr  error
		wantTodo *domain.Todo
	}{
		{"正常系", found, nil, nil, found},
		{"not found", nil, domain.ErrNotFound, domain.ErrNotFound, nil},
		{"DB error", nil, errors.New("db error"), errors.New("db error"), nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockTodoRepository{
				FindByIDFunc: func(ctx context.Context, gotID uuid.UUID) (*domain.Todo, error) {
					assert.Equal(t, id, gotID)
					return tt.mock, tt.repoErr
				},
			}
			uc := usecase.NewTodoUsecase(repo)
			todo, err := uc.Get(context.Background(), id)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Nil(t, todo)
				if errors.Is(tt.wantErr, domain.ErrNotFound) {
					assert.ErrorIs(t, err, domain.ErrNotFound)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantTodo, todo)
			}
		})
	}
}

func TestTodoUsecase_Update(t *testing.T) {
	id := uuid.New()
	ptr := func(s string) *string { return &s }

	tests := []struct {
		name       string
		input      usecase.UpdateInput
		findErr    error
		updateErr  error
		wantErr    error
		wantUpdate bool
		wantTitle  string
		wantStatus domain.Status
	}{
		{
			name:       "タイトルのみ",
			input:      usecase.UpdateInput{Title: ptr("掃除")},
			wantUpdate: true,
			wantTitle:  "掃除",
			wantStatus: domain.StatusPending,
		},
		{
			name:       "ステータスのみ",
			input:      usecase.UpdateInput{Status: ptr("done")},
			wantUpdate: true,
			wantTitle:  "買い物",
			wantStatus: domain.StatusDone,
		},
		{
			name:       "両方",
			input:      usecase.UpdateInput{Title: ptr("掃除"), Status: ptr("in_progress")},
			wantUpdate: true,
			wantTitle:  "掃除",
			wantStatus: domain.StatusInProgress,
		},
		{
			name:    "どちらも未指定",
			input:   usecase.UpdateInput{},
			wantErr: domain.ErrBadRequest,
		},
		{
			name:    "タイトル不正",
			input:   usecase.UpdateInput{Title: ptr("   ")},
			wantErr: domain.ErrInvalidTitle,
		},
		{
			name:    "ステータス不正",
			input:   usecase.UpdateInput{Status: ptr("foo")},
			wantErr: domain.ErrInvalidStatus,
		},
		{
			name:    "not found",
			input:   usecase.UpdateInput{Title: ptr("掃除")},
			findErr: domain.ErrNotFound,
			wantErr: domain.ErrNotFound,
		},
		{
			name:       "DB error",
			input:      usecase.UpdateInput{Title: ptr("掃除")},
			updateErr:  errors.New("db error"),
			wantErr:    errors.New("db error"),
			wantUpdate: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateCalled := false
			repo := &mockTodoRepository{
				FindByIDFunc: func(ctx context.Context, gotID uuid.UUID) (*domain.Todo, error) {
					assert.Equal(t, id, gotID)
					if tt.findErr != nil {
						return nil, tt.findErr
					}
					return &domain.Todo{ID: id, Title: "買い物", Status: domain.StatusPending}, nil
				},
				UpdateFunc: func(ctx context.Context, todo *domain.Todo) error {
					updateCalled = true
					return tt.updateErr
				},
			}
			uc := usecase.NewTodoUsecase(repo)
			todo, err := uc.Update(context.Background(), id, tt.input)

			assert.Equal(t, tt.wantUpdate, updateCalled)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Nil(t, todo)
				var appErr *domain.AppError
				if errors.As(tt.wantErr, &appErr) {
					assert.ErrorIs(t, err, tt.wantErr)
				}
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, id, todo.ID)
			assert.Equal(t, tt.wantTitle, todo.Title)
			assert.Equal(t, tt.wantStatus, todo.Status)
		})
	}
}
