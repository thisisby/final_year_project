package jwt

import "errors"

var (
	ErrTokenExpired     = errors.New("token is expired")
	ErrTokenInvalid     = errors.New("token is invalid")
	ErrTokenNotValidYet = errors.New("token not valid yet")
	ErrTokenMalformed   = errors.New("token is malformed")
)
