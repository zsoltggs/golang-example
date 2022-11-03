package service

import (
	"context"
	"fmt"
	"github.com/zsoltggs/golang-example/port-domain-service/internal/database"
	"github.com/zsoltggs/golang-example/port-domain-service/pkg/generated/github.com/zsoltggs/golang-example/pkg/pds"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Service struct {
	pds.UnimplementedServiceServer
	db database.Database
}

func New(db database.Database) *Service {
	return &Service{
		db: db,
	}
}

func (s Service) UpsertPort(ctx context.Context, request *pds.UpsertPortRequest) (*emptypb.Empty, error) {
	err := s.db.UpsertPort(ctx, request.GetPort())
	if err != nil {
		return nil, fmt.Errorf("unable to upsert port with id %q: %w", request.GetPort().GetId(), err)
	}
	return &emptypb.Empty{}, nil
}

func (s Service) GetPortsPaginated(ctx context.Context, request *pds.GetPortsPaginatedRequest) (*pds.GetPortsPaginatedResponse, error) {
	res, err := s.db.GetPortsPaginated(ctx, request.GetPagination())
	if err != nil {
		return nil, fmt.Errorf("unable to get ports for pagination %v: %w", request.GetPagination(), err)
	}
	return &pds.GetPortsPaginatedResponse{
		Ports: res,
	}, nil
}

func (s Service) GetPortByID(ctx context.Context, request *pds.GetPortByIDRequest) (*pds.GetPortByIDResponse, error) {
	res, err := s.db.GetPortByID(ctx, request.GetId())
	if err != nil {
		return nil, fmt.Errorf("unable to get ports for id %s: %w", request.GetId(), err)
	}
	return &pds.GetPortByIDResponse{
		Port: res,
	}, nil
}
