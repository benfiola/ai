package core

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/benfiola/ai/pkg/db/sqlc"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthenticateOpts struct {
	Email    string
	Password string
}

func (c *Core) Authenticate(ctx context.Context, opts AuthenticateOpts) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(opts.Password), 0)
	if err != nil {
		return "", err
	}

	id, err := c.DB.Queries.GetUserIdByCredentials(ctx, sqlc.GetUserIdByCredentialsParams{
		Email: opts.Email,
		Hash:  hash,
	})
	if err != nil {
		return "", err
	}
	if id == 0 {
		return "", fmt.Errorf("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": strconv.Itoa(int(id)),
			"exp": time.Now().Add(time.Hour * 24).Unix(),
		},
	)

	tokenString, err := token.SignedString(c.SecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

type authCtxKey struct{}

type AuthContext struct {
	Admin bool
	User  int
}

func (c *Core) WithAuth(ctx context.Context, tokenString string) context.Context {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) { return c.SecretKey, nil })
	if err != nil {
		return ctx
	}

	if !token.Valid {
		return ctx
	}

	idString, err := token.Claims.GetSubject()
	if err != nil {
		return ctx
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		return ctx
	}

	authContext := AuthContext{
		Admin: false,
		User:  id,
	}

	return context.WithValue(ctx, authCtxKey{}, authContext)
}

func (c *Core) GetAuth(ctx context.Context) AuthContext {
	authContext, ok := ctx.Value(authCtxKey{}).(AuthContext)
	if !ok {
		return AuthContext{}
	}
	return authContext
}
