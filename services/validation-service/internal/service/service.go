package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zsoltggs/golang-example/services/validation-service/internal/database"
	"github.com/zsoltggs/golang-example/services/validation-service/internal/validator"
	"github.com/zsoltggs/golang-example/services/validation-service/pkg/validationmodels"
)

var ErrInvalidJSON = errors.New("invalid JSON")
var ErrValidationError = errors.New("validation error")

type Service interface {
	UpsertSchema(ctx context.Context,
		req *validationmodels.UpsertSchemaRequest) (*validationmodels.UpsertSchemaResponse, error)
	GetSchemaByID(ctx context.Context,
		request *validationmodels.GetSchemaRequest) (*validationmodels.GetSchemaResponse, error)
	ValidateDocument(ctx context.Context,
		req *validationmodels.ValidateRequest) (*validationmodels.ValidateResponse, error)
}

type service struct {
	db           database.Database
	validatorSvc validator.Validator
}

func New(db database.Database, validatorSvc validator.Validator) Service {
	return &service{
		db:           db,
		validatorSvc: validatorSvc,
	}
}

func (s service) UpsertSchema(ctx context.Context,
	req *validationmodels.UpsertSchemaRequest) (*validationmodels.UpsertSchemaResponse, error) {
	if !json.Valid([]byte(req.Schema)) {
		return &validationmodels.UpsertSchemaResponse{
			HttpResponse: validationmodels.StatusHttpResponse{
				Action:  "uploadSchema",
				ID:      req.SchemaID,
				Status:  "error",
				Message: "Invalid JSON",
			},
		}, ErrInvalidJSON
	}

	err := s.db.UpsertSchema(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("unable to create schema: %w", err)
	}
	return &validationmodels.UpsertSchemaResponse{
		HttpResponse: validationmodels.StatusHttpResponse{
			Action: "uploadSchema",
			ID:     req.SchemaID,
			Status: "success",
		},
	}, nil
}

func (s service) GetSchemaByID(ctx context.Context,
	request *validationmodels.GetSchemaRequest) (*validationmodels.GetSchemaResponse, error) {
	result, err := s.db.GetSchema(ctx, request.ID)
	if err != nil {
		return nil, fmt.Errorf("unable to get by id: %w", err)
	}
	return &validationmodels.GetSchemaResponse{
		ID:     request.ID,
		Schema: result,
	}, nil
}

func (s service) ValidateDocument(ctx context.Context,
	req *validationmodels.ValidateRequest) (*validationmodels.ValidateResponse, error) {
	if !json.Valid([]byte(req.Document)) {
		return &validationmodels.ValidateResponse{
			HttpResponse: validationmodels.StatusHttpResponse{
				Action:  "validateDocument",
				ID:      req.SchemaID,
				Status:  "error",
				Message: "Invalid JSON",
			},
		}, ErrInvalidJSON
	}

	schema, err := s.db.GetSchema(ctx, req.SchemaID)
	if err != nil {
		//TODO Better error handling
		return nil, errors.New("not found")
	}
	documentWithoutNull, err := s.validatorSvc.RemoveNullValuesFromDoc(ctx, req.Document)
	if err != nil {
		return nil, fmt.Errorf("unable to remove null values from document: %w", err)
	}
	err = s.validatorSvc.Validate(ctx, validator.InputJson{
		Schema: schema,
		Doc:    documentWithoutNull,
	})
	if err != nil {
		switch {
		case errors.As(err, &validator.SchemaValidationError{}):
			return &validationmodels.ValidateResponse{
				HttpResponse: validationmodels.StatusHttpResponse{
					Action:  "validateDocument",
					ID:      req.SchemaID,
					Status:  "error",
					Message: err.Error(),
				},
			}, ErrValidationError
		default:
			return nil, fmt.Errorf("unable to validate document: %w", err)
		}

	}

	return &validationmodels.ValidateResponse{
		HttpResponse: validationmodels.StatusHttpResponse{
			Action: "validateDocument",
			ID:     req.SchemaID,
			Status: "success",
		},
	}, nil
}
