package core

import (
	"context"
	"net/url"
	"testing"

	"github.com/benfiola/ai/pkg/database"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TNewCore(t *testing.T) *Core {
	dbContainer, err := postgres.Run(context.Background(), "postgres:14.18-alpine", postgres.BasicWaitStrategies())
	testcontainers.CleanupContainer(t, dbContainer)
	require.NoError(t, err, "db container failed to create")
	dbConnString, err := dbContainer.ConnectionString(context.Background(), "sslmode=disable")
	require.NoError(t, err, "failed to determine db connection string")
	dbUrl, err := url.Parse(dbConnString)
	require.NoError(t, err, "db url unparseable")
	db, err := database.New(database.Opts{
		URL: dbUrl,
	})
	require.NoError(t, err, "db object creation failed")

	err = db.Migrate(database.VersionHead)
	require.NoError(t, err, "db migration failed")

	core := Core{
		DB:        db,
		SecretKey: "secret-key",
	}

	return &core
}
