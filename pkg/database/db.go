package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/benfiola/ai/pkg/database/sqlc"
	_ "github.com/jackc/pgx/v5"
)

type DB struct {
	Logger  *slog.Logger
	Pool    *sql.DB
	Queries *sqlc.Queries
	URL     *url.URL
}

type Opts struct {
	Logger *slog.Logger
	URL    *url.URL
}

func New(opts Opts) (*DB, error) {
	logger := opts.Logger
	if logger == nil {
		logger = slog.New(slog.DiscardHandler)
	}

	var pool *sql.DB
	var err error
	switch opts.URL.Scheme {
	case "postgres":
		pool, err = sql.Open("postgres", opts.URL.String())
	default:
		err = fmt.Errorf("unrecognized db type %s", opts.URL.Scheme)
	}
	if err != nil {
		return nil, err
	}

	queries := sqlc.New(pool)

	db := DB{
		Logger:  logger,
		Pool:    pool,
		URL:     opts.URL,
		Queries: queries,
	}

	return &db, nil
}
