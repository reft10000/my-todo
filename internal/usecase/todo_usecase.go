package usecase

import (
	"context"
	"fmt"
	"my-todo/internal/domain"
)

type TodoUsecase struct {
	repo domain.TodoRepository
}

func NewTodoUsecase(repo domain.TodoRepository) *TodoUsecase {
	return &TodoUsecase{repo: repo}
}

func (u *TodoUsecase) Create(title string) (*domain.Todo, error) {
	todo, err := domain.NewTodo(title)
	if err != nil {
		return nil, fmt.Errorf("failed to create todo: %w", err)
	}
	if err := u.repo.Create(todo); err != nil {
		return nil, fmt.Errorf("failed to save todo: %w", err)
	}
	return todo, nil
}

type ListInput struct {
	Page  int
	Limit int
}

func (u *TodoUsecase) List(ctx context.Context, input ListInput) (*domain.ListResult[domain.Todo], error) {
	params := domain.ListParams{
		Pagination: domain.PaginationParams{
			Page:  input.Page,
			Limit: input.Limit,
		},
	}

	result, err := u.repo.List(ctx, params)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
