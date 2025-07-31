package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/url"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/sqlc-dev/sqlc"
)

type DB struct {
	Logger *slog.Logger
	Pool   *sql.DB
	URL    *url.URL
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
	case "sqlite":
		dsn := fmt.Sprintf("file:%s?%s", opts.URL.Path, opts.URL.RawQuery)
		pool, err = sql.Open("sqlite3", dsn)
	case "postgres":
		pool, err = sql.Open("postgres", opts.URL.String())
	default:
		err = fmt.Errorf("unrecognized db type %s", opts.URL.Scheme)
	}
	if err != nil {
		return nil, err
	}

	db := DB{
		Logger: logger,
		Pool:   pool,
		URL:    opts.URL,
	}

	return &db, nil
}
