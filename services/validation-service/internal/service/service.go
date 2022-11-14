package service

import (
	"context"
	"fmt"
	"github.com/zsoltggs/golang-example/pkg/users"
	"github.com/zsoltggs/golang-example/services/validation-service/internal/database"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Service struct {
	users.UnimplementedServiceServer
	db database.Database
}

func New(db database.Database) *Service {
	return &Service{
		db: db,
	}
}

func (s Service) CreateUser(ctx context.Context, request *users.CreateUserRequest) (*emptypb.Empty, error) {
	if request.GetUser().GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	err := s.db.Create(ctx, request.GetUser())
	if err != nil {
		return nil, fmt.Errorf("unable to create user: %w", err)
	}
	return &emptypb.Empty{}, nil
}

func (s Service) PatchUser(ctx context.Context, request *users.PatchUserRequest) (*emptypb.Empty, error) {
	if request.GetUser().GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	err := s.db.Patch(ctx, request.GetUser())
	if err != nil {
		return nil, fmt.Errorf("unable to patch user: %w", err)
	}
	return &emptypb.Empty{}, nil
}

func (s Service) DeleteUser(ctx context.Context, request *users.DeleteUserRequest) (*emptypb.Empty, error) {
	if request.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	err := s.db.Delete(ctx, request.GetId())
	if err != nil {
		return nil, fmt.Errorf("unable to delete user: %w", err)
	}
	return &emptypb.Empty{}, nil
}

func (s Service) ListUsers(ctx context.Context, request *users.ListUsersRequest) (*users.ListUsersResponse, error) {
	// Make sure pagination is always present
	request.Pagination = defaultPagination(request.Pagination)
	results, err := s.db.List(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("unable to list users: %w", err)
	}
	return &users.ListUsersResponse{
		Users: results,
	}, nil
}
func defaultPagination(pagination *users.Pagination) *users.Pagination {
	if pagination == nil {
		return &users.Pagination{
			Offset: 0,
			Limit:  10,
		}
	}
	newPagination := users.Pagination{
		Offset: pagination.GetOffset(),
		Limit:  pagination.GetLimit(),
	}

	if newPagination.Limit == 0 {
		newPagination.Limit = 10
	}

	return &newPagination
}

func (s Service) GetByID(ctx context.Context, request *users.GetUserRequest) (*users.GetUserResponse, error) {
	result, err := s.db.GetByID(ctx, request.GetId())
	if err != nil {
		return nil, fmt.Errorf("unable to get by id: %w", err)
	}
	return &users.GetUserResponse{
		User: result,
	}, nil
}
