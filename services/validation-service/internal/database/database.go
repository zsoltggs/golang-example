package database

import (
	"context"
	"errors"
	"github.com/zsoltggs/golang-example/services/validation-service/pkg/validationmodels"
)

var ErrNotFound = errors.New("error not found")

type Database interface {
	UpsertSchema(ctx context.Context, request *validationmodels.UpsertSchemaRequest) error
	GetSchema(ctx context.Context, ID string) (string, error)
}
