package usecase

import (
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
