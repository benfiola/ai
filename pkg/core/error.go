package core

type ErrInvalidCredentials struct{}

func (e ErrInvalidCredentials) Error() string {
	return "invalid credentials"
}

type ErrUnauthorized struct{}

func (e ErrUnauthorized) Error() string {
	return "unauthorized"
}
