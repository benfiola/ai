package core

import (
	"fmt"
	"log/slog"

	"github.com/benfiola/ai/pkg/db"
)

type Core struct {
	DB        *db.DB
	Logger    *slog.Logger
	SecretKey string
}

type Opts struct {
	DB        *db.DB
	Logger    *slog.Logger
	SecretKey string
}

func New(opts Opts) (*Core, error) {
	db := opts.DB
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	logger := opts.Logger
	if logger == nil {
		logger = slog.New(slog.DiscardHandler)
	}

	core := Core{
		DB:        db,
		Logger:    logger,
		SecretKey: "secret-key",
	}

	return &core, nil
}
