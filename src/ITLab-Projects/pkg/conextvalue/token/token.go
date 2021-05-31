package token

import (
	"context"
	"errors"
)

var (
	ErrTokenNotFound = errors.New("Token not found in ctx")
)

type TokenKey struct{}

type TokenContext struct {
	context.Context
}

func New(
	ctx 	context.Context,
	token	string,	
) *TokenContext {
	return &TokenContext{
		Context: context.WithValue(
			ctx,
			TokenKey{},
			token,
		),
	}
}

func GetTokenFromContext(
	ctx		context.Context,
) (string, error) {
	val := ctx.Value(TokenKey{})
	if token, ok := val.(string); !ok {
		return "", ErrTokenNotFound
	} else {
		return token, nil
	}
}