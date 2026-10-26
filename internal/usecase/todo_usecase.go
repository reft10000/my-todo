package usecase

import (
	"context"
	"fmt"
	"my-todo/internal/domain"

	"github.com/google/uuid"
)

type TodoUsecase struct {
	repo domain.TodoRepository
}

func NewTodoUsecase(repo domain.TodoRepository) *TodoUsecase {
	return &TodoUsecase{repo: repo}
}

func (u *TodoUsecase) Create(ctx context.Context, title string) (*domain.Todo, error) {
	todo, err := domain.NewTodo(title)
	if err != nil {
		return nil, fmt.Errorf("failed to create todo: %w", err)
	}
	if err := u.repo.Create(ctx, todo); err != nil {
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

func (u *TodoUsecase) Get(ctx context.Context, id uuid.UUID) (*domain.Todo, error) {
	todo, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get todo: %w", err)
	}
	return todo, nil
}

type UpdateInput struct {
	Title  *string
	Status *string
}

func (u *TodoUsecase) Update(ctx context.Context, id uuid.UUID, input UpdateInput) (*domain.Todo, error) {
	if input.Title == nil && input.Status == nil {
		return nil, domain.ErrBadRequest
	}

	todo, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get todo: %w", err)
	}

	if input.Title != nil {
		if err := todo.ChangeTitle(*input.Title); err != nil {
			return nil, err
		}
	}
	if input.Status != nil {
		if err := todo.ChangeStatus(domain.Status(*input.Status)); err != nil {
			return nil, err
		}
	}

	if err := u.repo.Update(ctx, todo); err != nil {
		return nil, fmt.Errorf("failed to update todo: %w", err)
	}
	return todo, nil
}

func (u *TodoUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	if err := u.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}
	return nil
}
