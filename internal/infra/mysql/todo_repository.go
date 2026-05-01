package mysql

import (
	"context"
	"fmt"
	"my-todo/internal/domain"
	"my-todo/internal/infra/ent"
	ent_todo "my-todo/internal/infra/ent/todo"
)

var _ domain.TodoRepository = (*TodoRepository)(nil)

type TodoRepository struct {
	client *ent.Client
}

func NewTodoRepository(client *ent.Client) *TodoRepository {
	return &TodoRepository{client: client}
}

func (r *TodoRepository) Create(todo *domain.Todo) error {
	_, err := r.client.Todo.
		Create().
		SetID(todo.ID).
		SetTitle(todo.Title).
		SetStatus(ent_todo.Status(todo.Status)).
		SetCreatedAt(todo.CreatedAt).
		Save(context.Background())
	if err != nil {
		return fmt.Errorf("failed to create todo: %w", err)
	}
	return nil
}
