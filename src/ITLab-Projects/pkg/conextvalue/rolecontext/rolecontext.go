package rolecontext

import (
	"context"
	"errors"
)

var (
	ErrRoleNotFound = errors.New("Role not found in context")
)

type RoleKey struct{}

type RoleContext struct {
	context.Context
}

func New(ctx context.Context, role string) *RoleContext {
	return &RoleContext{
		Context: context.WithValue(
			ctx,
			RoleKey{},
			role,
		),
	}
}

func GetRoleFromContext(ctx context.Context) (string, error) {
	val := ctx.Value(RoleKey{})
	if role, ok := val.(string); !ok {
		return "", ErrRoleNotFound
	} else {
		return role, nil
	}
}