package database

import (
	"context"
	"errors"
	"github.com/zsoltggs/golang-example/pkg/users"
)

var ErrNotFound = errors.New("error not found")

type Database interface {
	Create(ctx context.Context, user *users.User) error
	Patch(ctx context.Context, user *users.User) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*users.User, error)
	List(ctx context.Context, request *users.ListUsersRequest) ([]*users.User, error)
	Health() error
}
