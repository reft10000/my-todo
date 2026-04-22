package mysql_test

import (
	"context"
	"testing"

	"my-todo/internal/infra/ent"
	"my-todo/internal/infra/mysql"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tcmysql "github.com/testcontainers/testcontainers-go/modules/mysql"
)

func SetupTestDB(t *testing.T) *ent.Client {
	ctx := context.Background()

	container, err := tcmysql.Run(ctx,
		"mysql:8",
		tcmysql.WithDatabase("todo"),
		tcmysql.WithUsername("root"),
		tcmysql.WithPassword("password"),
	)
	require.NoError(t, err)
	t.Cleanup(func() { container.Terminate(ctx) })

	dsn, err := container.ConnectionString(ctx, "parseTime=true")
	assert.NoError(t, err)

	client, err := mysql.NewClient(dsn)
	assert.NoError(t, err)
	t.Cleanup(func() { client.Close() })

	return client
}
