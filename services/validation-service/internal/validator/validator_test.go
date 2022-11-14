package validator_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zsoltggs/golang-example/services/validation-service/internal/validator"
	"io"
	"os"
	"testing"
)

const (
	configDocFileName    = "testdata/config.json"
	configSchemaFileName = "testdata/config-schema.json"
)

func readFileContents(t *testing.T, fileName string) string {
	f, err := os.Open(configDocFileName)
	require.NoError(t, err)
	fContents, err := io.ReadAll(f)
	require.NoError(t, err)
	return string(fContents)
}

func Test_SchemaValidationSuccessful(t *testing.T) {
	doc := readFileContents(t, configDocFileName)
	schema := readFileContents(t, configSchemaFileName)
	svc := validator.New()
	ctx := context.Background()

	err := svc.Validate(ctx, validator.JsonPair{
		Schema: schema,
		Doc:    doc,
	})
	require.NoError(t, err)
}

func Test_InvalidSchemas(t *testing.T) {
	svc := validator.New()
	ctx := context.Background()
	schema := readFileContents(t, configSchemaFileName)

	tests := map[string]struct {
		doc           string
		isErrExpected bool
		expectedErr   string
	}{
		"empty": {
			doc:           "{}",
			isErrExpected: false,
		},
		"provided-example": {
			doc:           readFileContents(t, configDocFileName),
			isErrExpected: false,
		},
		"required-chunk-info-missing": {
			doc: `{
					  "source": "/home/alice/image.iso",
					  "destination": "/mnt/storage",
					  "timeout": "2",
  					  "chunks": {
						number: 2
					  }
				   }
`,
			isErrExpected: true,
			expectedErr:   "asd",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := svc.Validate(ctx, validator.JsonPair{
				Schema: schema,
				Doc:    tc.doc,
			})
			if tc.isErrExpected {
				assert.Error(t, err)
				if tc.expectedErr != "" {
					assert.EqualError(t, err, tc.expectedErr)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
