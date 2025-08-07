package core

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthenticateOpts struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *Core) Authenticate(ctx context.Context, opts AuthenticateOpts) (string, error) {
	user, err := c.DB.Queries.GetUserByEmail(ctx, opts.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ErrInvalidCredentials{}
		}
		return "", err
	}

	err = bcrypt.CompareHashAndPassword(user.Hash, []byte(opts.Password))
	if err != nil {
		return "", ErrInvalidCredentials{}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub": strconv.Itoa(int(user.ID)),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	signedToken, err := token.SignedString([]byte(c.SecretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

type ctxAuthInfo struct{}

type AuthInfo struct {
	User int
}

func (c *Core) EmbedAuthInfo(ctx context.Context, signedToken string) (context.Context, error) {
	token, err := jwt.Parse(signedToken, func(t *jwt.Token) (any, error) { return []byte(c.SecretKey), nil })
	if err != nil {
		return ctx, ErrInvalidCredentials{}
	}

	idString, err := token.Claims.GetSubject()
	if err != nil {
		return ctx, err
	}
	id, err := strconv.Atoi(idString)
	if err != nil {
		return ctx, err
	}

	user, err := c.DB.Queries.GetUserById(ctx, int32(id))
	if err != nil {
		return ctx, err
	}

	authInfo := AuthInfo{
		User: int(user.ID),
	}

	newCtx := context.WithValue(ctx, ctxAuthInfo{}, authInfo)

	return newCtx, nil
}

func (c *Core) getAuthInfo(ctx context.Context) AuthInfo {
	authInfo := ctx.Value(ctxAuthInfo{})
	if authInfo == nil {
		authInfo = AuthInfo{}
	}
	return authInfo.(AuthInfo)
}


