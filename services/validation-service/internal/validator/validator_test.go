package validator_test

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/zsoltggs/golang-example/services/validation-service/internal/validator"
	"io"
	"os"
	"testing"
)

const (
	configDocFileName    = "../../resources/config.json"
	configSchemaFileName = "../../resources/config-schema.json"
)

func readFileContents(t *testing.T, fileName string) string {
	f, err := os.Open(fileName)
	require.NoError(t, err)
	fContents, err := io.ReadAll(f)
	require.NoError(t, err)
	return string(fContents)
}

func Test_SchemaValidationScenarios(t *testing.T) {
	svc := validator.New()
	ctx := context.Background()
	schema := readFileContents(t, configSchemaFileName)

	tests := map[string]struct {
		doc                   string
		isErrExpected         bool
		expectedErrorMessages []string
	}{

		"empty": {
			doc:           "{}",
			isErrExpected: true,
			expectedErrorMessages: []string{
				"Property (root) source is required",
				"Property (root) destination is required",
			},
		},
		"provided-example": {
			doc:           readFileContents(t, configDocFileName),
			isErrExpected: false,
		},
		"required-chunk-info-missing": {
			doc: `{
					  "source": "/home/alice/image.iso",
					  "destination": "/mnt/storage",
					  "timeout": 6,
  					  "chunks": {
						"number": 2
					  }
				   }
`,
			isErrExpected:         true,
			expectedErrorMessages: []string{"Property chunks size is required"},
		},
		"multiple-err-scenario": {
			doc: `{
					  "source": "/home/alice/image.iso",
					  "destination": "/mnt/storage",
					  "timeout": -10,
  					  "chunks": {
						"number": 2
					  }
				   }
`,
			isErrExpected: true,
			expectedErrorMessages: []string{
				"Property timeout Must be greater than or equal to 0",
				"Property chunks size is required",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			docWithoutNulls, err := svc.RemoveNullValuesFromDoc(ctx, tc.doc)
			require.NoError(t, err)

			err = svc.Validate(ctx, validator.InputJson{
				Schema: schema,
				Doc:    docWithoutNulls,
			})
			if tc.isErrExpected {
				require.Error(t, err)
				var typedErr validator.SchemaValidationError
				require.ErrorAs(t, err, &typedErr)
				if len(tc.expectedErrorMessages) != 0 {
					require.Len(t, tc.expectedErrorMessages, len(typedErr.Errors))
					for _, errMsg := range tc.expectedErrorMessages {
						require.Contains(t, typedErr.Errors, errMsg)
					}
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}
