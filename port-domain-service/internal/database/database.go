package database

import (
	"context"
	"errors"
	"github.com/zsoltggs/golang-example/port-domain-service/pkg/generated/github.com/zsoltggs/golang-example/pkg/pds"
)

var ErrNotFound = errors.New("error not found")

type Database interface {
	UpsertPort(ctx context.Context, port *pds.Port) error
	GetPortsPaginated(ctx context.Context, pagination *pds.Pagination) ([]*pds.Port, error)
	GetPortByID(ctx context.Context, ID string) (*pds.Port, error)
}
