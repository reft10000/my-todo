package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTodo(t *testing.T) {
	todo, err := NewTodo("買い物")
	assert.NoError(t, err)
	assert.NotEmpty(t, todo.ID)
	assert.Equal(t, "買い物", todo.Title)
	assert.Equal(t, StatusPending, todo.Status)
}

func TestIsValidTitle(t *testing.T) {
	tests := []struct {
		name  string
		title string
		want  bool
	}{
		{"正常系", "買い物", true},
		{"空文字", "", false},
		{"スペースのみ", "   ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, isValidTitle(tt.title))
		})
	}
}
