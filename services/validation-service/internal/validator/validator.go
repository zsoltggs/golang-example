package validator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"reflect"
	"strings"
)

type Validator struct {
}

// TODO Interface
func New() *Validator {
	return &Validator{}
}

type InputJson struct {
	Schema string
	Doc    string
}

type ErrorInfo struct {
	Field, Description string
}

func (*Validator) Validate(ctx context.Context, p InputJson) error {
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
			allErrors = append(allErrors, fmt.Sprintf("field: %s, description: %s", err.Field(), err.Description()))
		}
		return errors.New(strings.Join(allErrors, "\n"))
	}
	return nil
}

func (*Validator) RemoveNullValuesFromDoc(ctx context.Context, p InputJson) (string, error) {
	convertedMap := map[string]interface{}{}
	err := json.Unmarshal([]byte(p.Doc), &convertedMap)
	if err != nil {
		return "", err
	}
	removeNulls(convertedMap)
	res, err := json.MarshalIndent(convertedMap, "", "  ")
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func removeNulls(m map[string]interface{}) {
	val := reflect.ValueOf(m)
	for _, e := range val.MapKeys() {
		v := val.MapIndex(e)
		if v.IsNil() {
			delete(m, e.String())
			continue
		}
		switch t := v.Interface().(type) {
		// If key is a JSON object (Go Map), use recursion to go deeper
		case map[string]interface{}:
			removeNulls(t)
		}
	}
}
