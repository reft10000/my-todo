package mysql

import (
	"context"
	"fmt"
	"my-todo/internal/infra/ent"

	_ "github.com/go-sql-driver/mysql"
)

func NewClient(dsn string) (*ent.Client, error) {
	client, err := ent.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open ent client: %w", err)
	}
	if err := client.Schema.Create(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}
	return client, nil
}
