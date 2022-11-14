package service_test

import (
	"context"
	"github.com/go-test/deep"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zsoltggs/golang-example/pkg/users"
	"github.com/zsoltggs/golang-example/services/users/internal/database"
	"github.com/zsoltggs/golang-example/services/users/internal/notifier"
	"github.com/zsoltggs/golang-example/services/users/internal/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

var (
	connString = "mongodb://localhost:27017"
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

func Test_Create_NoID(t *testing.T) {
	truncate(t, t.Name())
	db, err := database.NewMongo(connString, t.Name())
	require.NoError(t, err)
	svc := service.New(db, notifier.NewLogNotifier())
	ctx := context.Background()

	_, err = svc.CreateUser(ctx, &users.CreateUserRequest{})
	require.Error(t, err)
}

func Test_Create_Success(t *testing.T) {
	truncate(t, t.Name())
	db, err := database.NewMongo(connString, t.Name())
	require.NoError(t, err)
	svc := service.New(db, notifier.NewLogNotifier())

	ctx := context.Background()
	userID := uuid.New().String()
	expected := getTestUser(userID)
	_, err = svc.CreateUser(ctx, &users.CreateUserRequest{
		User: expected,
	})
	require.NoError(t, err)

	res, err := svc.GetByID(ctx, &users.GetUserRequest{
		Id: userID,
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	if diff := deep.Equal(expected, res.GetUser()); diff != nil {
		t.Error(diff)
	}
}

func Test_Patch_Success(t *testing.T) {
	truncate(t, t.Name())
	db, err := database.NewMongo(connString, t.Name())
	require.NoError(t, err)
	svc := service.New(db, notifier.NewLogNotifier())

	ctx := context.Background()
	userID := uuid.New().String()
	expected := getTestUser(userID)
	_, err = svc.CreateUser(ctx, &users.CreateUserRequest{
		User: expected,
	})
	require.NoError(t, err)

	expected.FirstName = "replaced"

	_, err = svc.PatchUser(ctx, &users.PatchUserRequest{
		User: &users.User{
			Id:        userID,
			FirstName: "replaced",
		}})
	require.NoError(t, err)
	res, err := svc.GetByID(ctx, &users.GetUserRequest{
		Id: userID,
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	if diff := deep.Equal(expected, res.GetUser()); diff != nil {
		t.Error(diff)
	}
}

func Test_Delete_Success(t *testing.T) {
	truncate(t, t.Name())
	db, err := database.NewMongo(connString, t.Name())
	require.NoError(t, err)
	svc := service.New(db, notifier.NewLogNotifier())

	ctx := context.Background()
	userID := uuid.New().String()
	input := getTestUser(userID)
	_, err = svc.CreateUser(ctx, &users.CreateUserRequest{
		User: input,
	})
	require.NoError(t, err)

	_, err = svc.DeleteUser(ctx, &users.DeleteUserRequest{
		Id: userID,
	})
	require.NoError(t, err)
	res, err := svc.GetByID(ctx, &users.GetUserRequest{
		Id: userID,
	})
	require.Error(t, err)
	require.Nil(t, res)
	require.ErrorIs(t, err, database.ErrNotFound)
}

func Test_List(t *testing.T) {
	truncate(t, t.Name())
	db, err := database.NewMongo(connString, t.Name())
	require.NoError(t, err)
	svc := service.New(db, notifier.NewLogNotifier())

	ctx := context.Background()
	filter := uuid.New().String()
	userID := uuid.New().String()
	first := getTestUser(userID)
	first.Email = filter
	_, err = svc.CreateUser(ctx, &users.CreateUserRequest{
		User: first,
	})

	require.NoError(t, err)
	userID2 := uuid.New().String()
	second := getTestUser(userID2)
	second.Email = filter
	second.CreatedAt = timestamppb.New(time.Now().Add(2 * time.Second))
	_, err = svc.CreateUser(ctx, &users.CreateUserRequest{
		User: second,
	})
	require.NoError(t, err)

	// with no filters returns all
	res, err := svc.ListUsers(ctx, &users.ListUsersRequest{
		Filters: &users.Filters{
			Email: filter,
		},
	})
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Len(t, res.GetUsers(), 2)

	res, err = svc.ListUsers(ctx, &users.ListUsersRequest{
		Pagination: &users.Pagination{
			Offset: 0,
			Limit:  1,
		},
		Filters: &users.Filters{
			Email: filter,
		},
	})
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Len(t, res.GetUsers(), 1)
	assert.Equal(t, userID, res.GetUsers()[0].GetId(), userID)

	res, err = svc.ListUsers(ctx, &users.ListUsersRequest{
		Pagination: &users.Pagination{
			Offset: 1,
			Limit:  1,
		},
		Filters: &users.Filters{
			Email: filter,
		},
	})
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Len(t, res.GetUsers(), 1)
	assert.Equal(t, userID2, res.GetUsers()[0].GetId())
}

func getTestUser(userID string) *users.User {
	return &users.User{
		Id:        userID,
		FirstName: "name",
		LastName:  "last",
		Nickname:  "nick-name",
		Password:  "pw",
		Email:     "email",
		Country:   "country",
		CreatedAt: timestamppb.Now(),
		UpdatedAt: timestamppb.Now(),
	}
}
