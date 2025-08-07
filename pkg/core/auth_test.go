package core

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestAuthenticate(t *testing.T) {
	t.Run("invalid credentials", func(t *testing.T) {
		core := TNewCore(t)

		token, err := core.Authenticate(context.Background(), AuthenticateOpts{Email: "invalid@email.com", Password: "invalid"})
		require.ErrorIs(t, err, ErrInvalidCredentials{}, "incorrect error raised")
		require.Zero(t, token, "non-empty token received")
	})

	t.Run("valid credentials", func(t *testing.T) {
		core := TNewCore(t)
		user := TCreateUser(t, core)

		token, err := core.Authenticate(context.Background(), AuthenticateOpts{Email: user.Email, Password: user.Password})
		approxExpiry := time.Now().Add(time.Hour * 24)
		require.NoError(t, err, "authentication failed")
		parsed, err := jwt.Parse(token, func(t *jwt.Token) (any, error) { return []byte(core.SecretKey), nil })
		require.NoError(t, err, "jwt unparseable")

		subject, err := parsed.Claims.GetSubject()
		require.NoError(t, err, "token missing subject claim")
		require.Equal(t, subject, strconv.Itoa(user.ID), "token subject mismatch")

		expiry, err := parsed.Claims.GetExpirationTime()
		require.NoError(t, err, "token missing expiry claim")
		require.True(t, approxExpiry.Sub(expiry.Time) < time.Second*1, "token expiry is unexpected value")
	})
}

func TestEmbedAuthInfo(t *testing.T) {
	t.Run("invalid token", func(t *testing.T) {
		core := TNewCore(t)

		ctx := context.Background()
		ctx, err := core.EmbedAuthInfo(ctx, "invalid")
		require.ErrorIs(t, err, ErrInvalidCredentials{}, "incorrect error raised")
		authInfo := ctx.Value(ctxAuthInfo{})
		require.Nil(t, authInfo, "auth info embedded despite error")
	})

	t.Run("valid token", func(t *testing.T) {
		core := TNewCore(t)
		user := TCreateUser(t, core)

		token, err := core.Authenticate(t.Context(), AuthenticateOpts{Email: user.Email, Password: user.Password})
		require.NoError(t, err, "authentication failed")

		ctx, err := core.EmbedAuthInfo(context.Background(), token)
		require.NoError(t, err, "token invalid")

		authInfo := ctx.Value(ctxAuthInfo{}).(AuthInfo)
		require.Equal(t, authInfo.User, user.ID, "user, auth info mismatch")
	})
}
