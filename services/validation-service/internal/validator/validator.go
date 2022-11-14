package validator

import (
	"context"
	"errors"
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"strings"
)

type Validator struct {
}

// TODO Interface
func New() *Validator {
	return &Validator{}
}

type JsonPair struct {
	Schema, Doc string
}

func (*Validator) Validate(ctx context.Context, p JsonPair) error {
	schemaLoader := gojsonschema.NewStringLoader(p.Schema)
	gSchema, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		return fmt.Errorf("unable to create schema: %w", err)
	}
	docLoader := gojsonschema.NewStringLoader(p.Doc)
	result, err := gSchema.Validate(docLoader)
	if err != nil {
		return fmt.Errorf("error validating schema: %w", err)
	}
	if !result.Valid() {
		var allErrors []string
		for _, err := range result.Errors() {
			// Err implements the ResultError interface
			allErrors = append(allErrors, err.Description())
		}
		return errors.New(strings.Join(allErrors, ","))
	}
	return nil
}

func (*Validator) RemoveNullValues(ctx context.Context, schema, doc string) error {
	return nil
}
