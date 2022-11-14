package service_test

import (
	"context"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zsoltggs/golang-example/services/validation-service/internal/database"
	"github.com/zsoltggs/golang-example/services/validation-service/internal/service"
	"github.com/zsoltggs/golang-example/services/validation-service/internal/validator"
	"github.com/zsoltggs/golang-example/services/validation-service/pkg/validationmodels"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"os"
	"testing"
)

var (
	connString = "mongodb://localhost:27017"
	dbName     = "schemas"
)

const (
	configDocNoNullFileName = "../../resources/config-no-nulls.json"
	configDocFileName       = "../../resources/config.json"
	configSchemaFileName    = "../../resources/config-schema.json"
	schemaID                = "config-schema"
)

// These are integration tests make sure to run MongoDB for this to work

func readFileContents(t *testing.T, fileName string) string {
	f, err := os.Open(fileName)
	require.NoError(t, err)
	fContents, err := io.ReadAll(f)
	require.NoError(t, err)
	return string(fContents)
}

func truncate(t *testing.T, databaseName string) {
	opts := options.Client().ApplyURI(connString)
	client, err := mongo.NewClient(opts)
	assert.Nil(t, err)
	err = client.Connect(context.Background())
	assert.Nil(t, err)
	err = client.Database(databaseName).Drop(context.Background())
	assert.Nil(t, err)
}

func Test_CreateAndGetSchema_Success(t *testing.T) {
	truncate(t, t.Name())
	db, err := database.NewMongo(connString, t.Name())
	require.NoError(t, err)
	validatorSvc := validator.Validator{}
	svc := service.New(db, &validatorSvc)
	ctx := context.Background()

	inputSchema := readFileContents(t, configSchemaFileName)
	res, err := svc.UpsertSchema(ctx, &validationmodels.UpsertSchemaRequest{
		SchemaID: schemaID,
		Schema:   inputSchema,
	})
	require.NoError(t, err)
	expected := &validationmodels.UpsertSchemaResponse{
		HttpResponse: validationmodels.StatusHttpResponse{
			Action: "uploadSchema",
			ID:     schemaID,
			Status: "success",
		},
	}
	if diff := deep.Equal(expected, res); diff != nil {
		t.Error(diff)
	}
	getRes, err := svc.GetSchemaByID(ctx, &validationmodels.GetSchemaRequest{
		ID: schemaID,
	})
	require.NoError(t, err)
	expectedGet := &validationmodels.GetSchemaResponse{
		ID:     schemaID,
		Schema: inputSchema,
	}
	if diff := deep.Equal(expectedGet, getRes); diff != nil {
		t.Error(diff)
	}
}

func Test_SchemaValidation_Success(t *testing.T) {
	truncate(t, t.Name())
	db, err := database.NewMongo(connString, t.Name())
	require.NoError(t, err)
	validatorSvc := validator.Validator{}
	svc := service.New(db, &validatorSvc)
	ctx := context.Background()

	inputSchema := readFileContents(t, configSchemaFileName)
	_, err = svc.UpsertSchema(ctx, &validationmodels.UpsertSchemaRequest{
		SchemaID: schemaID,
		Schema:   inputSchema,
	})
	require.NoError(t, err)

	res, err := svc.ValidateDocument(ctx, &validationmodels.ValidateRequest{
		SchemaID: schemaID,
		Document: readFileContents(t, configDocFileName),
	})
	require.NoError(t, err)
	expected := &validationmodels.ValidateResponse{
		HttpResponse: validationmodels.StatusHttpResponse{
			Action: "validateDocument",
			ID:     schemaID,
			Status: "success",
		},
	}
	if diff := deep.Equal(expected, res); diff != nil {
		t.Error(diff)
	}
}
