package core

import (
	"fmt"
	"log/slog"

	"github.com/benfiola/ai/pkg/database"
)

type Core struct {
	DB        *database.DB
	Logger    *slog.Logger
	SecretKey string
}

type Opts struct {
	DB        *database.DB
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

	secretKey := opts.SecretKey
	if secretKey == "" {
		secretKey = "secret-key"
	}
	if secretKey == "" {
		return nil, fmt.Errorf("secret key is nil")
	}

	core := Core{
		DB:        db,
		Logger:    logger,
		SecretKey: secretKey,
	}

	return &core, nil
}
