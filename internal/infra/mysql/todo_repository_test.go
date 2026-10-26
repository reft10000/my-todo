package mysql_test

import (
	"context"
	"testing"
	"time"

	"my-todo/internal/domain"
	"my-todo/internal/infra/mysql"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTodoRepository_Create(t *testing.T) {
	client := SetupTestDB(t)
	repo := mysql.NewTodoRepository(client)

	todo, err := domain.NewTodo("買い物")
	assert.NoError(t, err)

	err = repo.Create(todo)
	assert.NoError(t, err)
}

func TestTodoRepository_List(t *testing.T) {
	client := SetupTestDB(t)
	repo := mysql.NewTodoRepository(client)
	ctx := context.Background()

	titles := []string{"買い物", "掃除", "洗濯"}
	for _, title := range titles {
		todo, err := domain.NewTodo(title)
		require.NoError(t, err)
		err = repo.Create(todo)
		require.NoError(t, err)
	}

	result, err := repo.List(ctx, domain.ListParams{
		Pagination: domain.PaginationParams{Page: 1, Limit: 2},
	})

	assert.NoError(t, err)
	assert.Equal(t, 2, len(result.Todos))
	assert.Equal(t, 3, result.Total)
	assert.True(t, result.HasNext)
}

func TestTodoRepository_FindByID(t *testing.T) {
	client := SetupTestDB(t)
	repo := mysql.NewTodoRepository(client)
	ctx := context.Background()

	todo, err := domain.NewTodo("買い物")
	require.NoError(t, err)
	require.NoError(t, repo.Create(todo))

	t.Run("正常系", func(t *testing.T) {
		got, err := repo.FindByID(ctx, todo.ID)
		require.NoError(t, err)
		assert.Equal(t, todo.ID, got.ID)
		assert.Equal(t, todo.Title, got.Title)
		assert.Equal(t, todo.Status, got.Status)
		assert.WithinDuration(t, todo.CreatedAt, got.CreatedAt, time.Second)
	})

	t.Run("not found", func(t *testing.T) {
		got, err := repo.FindByID(ctx, uuid.New())
		assert.Nil(t, got)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}
