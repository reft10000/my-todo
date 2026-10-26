package domain

import (
	"strings"
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

func TestNewTodo_InvalidTitle(t *testing.T) {
	tests := []struct {
		name  string
		title string
	}{
		{"空文字", ""},
		{"スペースのみ", "   "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			todo, err := NewTodo(tt.title)
			assert.Nil(t, todo)
			assert.ErrorIs(t, err, ErrInvalidTitle)
		})
	}
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
		{"255文字", strings.Repeat("あ", 255), true},
		{"256文字", strings.Repeat("あ", 256), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, isValidTitle(tt.title))
		})
	}
}

func TestTodo_ChangeTitle(t *testing.T) {
	tests := []struct {
		name      string
		title     string
		wantErr   error
		wantTitle string
	}{
		{"正常系", "掃除", nil, "掃除"},
		{"空文字", "", ErrInvalidTitle, "買い物"},
		{"スペースのみ", "   ", ErrInvalidTitle, "買い物"},
		{"256文字", strings.Repeat("あ", 256), ErrInvalidTitle, "買い物"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			todo, err := NewTodo("買い物")
			assert.NoError(t, err)

			err = todo.ChangeTitle(tt.title)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantTitle, todo.Title)
		})
	}
}

func TestTodo_ChangeStatus(t *testing.T) {
	tests := []struct {
		name       string
		status     Status
		wantErr    error
		wantStatus Status
	}{
		{"pending", StatusPending, nil, StatusPending},
		{"in_progress", StatusInProgress, nil, StatusInProgress},
		{"done", StatusDone, nil, StatusDone},
		{"不正な値", Status("foo"), ErrInvalidStatus, StatusPending},
		{"空文字", Status(""), ErrInvalidStatus, StatusPending},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			todo, err := NewTodo("買い物")
			assert.NoError(t, err)

			err = todo.ChangeStatus(tt.status)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantStatus, todo.Status)
		})
	}
}
