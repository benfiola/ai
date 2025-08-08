package core

import (
	"context"

	"github.com/benfiola/ai/pkg/database/sqlc"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type CreateUserOpts struct {
	Email    string
	Password string
}

func (c *Core) CreateUser(ctx context.Context, opts CreateUserOpts) (int, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(opts.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	id, err := c.DB.Queries.CreateUser(ctx, sqlc.CreateUserParams{
		Email: opts.Email,
		Hash:  hash,
	})
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
