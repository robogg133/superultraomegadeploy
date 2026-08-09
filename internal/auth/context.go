package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Email     string
	FirstName string
	LastName  string
}

var userKey = struct{}{}

func ContextUser(ctx context.Context) (*User, bool) {
	u, ok := ctx.Value(userKey).(*User)
	return u, ok
}

func UserID(ctx context.Context) (uuid.UUID, bool) {
	u, ok := ContextUser(ctx)
	if !ok {
		return uuid.Nil, false
	}
	return u.ID, true
}

func UserFromContext(ctx context.Context) (User, error) {
	u, ok := ContextUser(ctx)
	if !ok {
		return User{}, errors.New("user not in context")
	}
	return *u, nil
}