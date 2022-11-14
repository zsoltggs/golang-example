package service_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

var (
	connString = "mongodb://localhost:27017"
	dbName     = "ports"
)

func truncate(t *testing.T, databaseName string) {
	opts := options.Client().ApplyURI(connString)
	client, err := mongo.NewClient(opts)
	assert.Nil(t, err)
	err = client.Connect(context.Background())
	assert.Nil(t, err)
	err = client.Database(databaseName).Drop(context.Background())
	assert.Nil(t, err)
}

/*
func Test_Create_NoID(t *testing.T) {
	truncate(t, t.Name())
	db, err := database.NewMongo(connString, t.Name())
	require.NoError(t, err)
	svc := service.New(db)
	ctx := context.Background()

	_, err = svc.CreateUser(ctx, &users.CreateUserRequest{})
	require.Error(t, err)
}
*/
