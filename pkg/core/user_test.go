package core

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var tUserCount = 0

type TUser struct {
	User
	Password string
}

func TCreateUser(t *testing.T, core *Core) TUser {
	email := fmt.Sprintf("user-%d@test.com", tUserCount)
	password := fmt.Sprintf("password-%d", tUserCount)
	tUserCount += 1

	userId, err := core.CreateUser(t.Context(), CreateUserOpts{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err, "test user creation failed")

	return TUser{
		User: User{
			ID:    userId,
			Email: email,
		},
		Password: password,
	}
}
