package core

import (
	"context"
	"fmt"

	"github.com/benfiola/ai/pkg/db/sqlc"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserOpts struct {
	Email    string
	Password string
}

func (c *Core) CreateUser(ctx context.Context, opts CreateUserOpts) (int, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(opts.Password), 0)
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

type User struct {
	ID    int
	Email string
}

func (c *Core) GetUser(ctx context.Context, id int) (*User, error) {
	authContext := c.GetAuth(ctx)
	if !authContext.Admin || authContext.User != id {
		return nil, fmt.Errorf("unauthorized")
	}

	dbUser, err := c.DB.Queries.GetUserById(ctx, int64(id))
	if err != nil {
		return nil, err
	}

	user := User{
		ID:    int(dbUser.ID),
		Email: dbUser.Email,
	}

	return &user, nil
}
