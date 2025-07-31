package core

import (
	"fmt"
	"log/slog"

	"github.com/benfiola/ai/pkg/db"
)

type Core struct {
	DB     *db.DB
	Logger *slog.Logger
}

type Opts struct {
	DB     *db.DB
	Logger *slog.Logger
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
		DB:     db,
		Logger: logger,
	}

	return &core, nil
}

func (c *Core) Health() error {
	return nil
}
