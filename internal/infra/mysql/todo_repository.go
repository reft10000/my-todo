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

func (r *TodoRepository) List(ctx context.Context, params domain.ListParams) (domain.ListResult[domain.Todo], error) {
	query := r.client.Todo.Query()

	total, err := query.Count(ctx)
	if err != nil {
		return domain.ListResult[domain.Todo]{}, err
	}

	offset := (params.Pagination.Page - 1) * params.Pagination.Limit
	todos, err := query.
		Order(ent.Asc("created_at")).
		Limit(params.Pagination.Limit).
		Offset(offset).
		All(ctx)
	if err != nil {
		return domain.ListResult[domain.Todo]{}, err
	}

	items := make([]domain.Todo, len(todos))
	for i, t := range todos {
		items[i] = domain.Todo{
			ID:        t.ID,
			Title:     t.Title,
			Status:    domain.Status(t.Status),
			CreatedAt: t.CreatedAt,
		}
	}

	return domain.ListResult[domain.Todo]{
		Todos:   items,
		Total:   total,
		Page:    params.Pagination.Page,
		Limit:   params.Pagination.Limit,
		HasNext: offset+len(todos) < total,
	}, nil
}
