package database

import (
	"context"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"github.com/zsoltggs/golang-example/pkg/users"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	collectionUsers = "users"
)

type mongoDB struct {
	client   *mongo.Client
	database string
}

func (s *mongoDB) Health() error {
	return s.client.Ping(context.Background(), nil)
}

func NewMongo(connString string, database string) (Database, error) {
	opts := options.Client().ApplyURI(connString).
		SetRetryWrites(true)

	err := opts.Validate()
	if err != nil {
		return nil, fmt.Errorf("failed to create mongo connection: %w", err)
	}

	client, err := mongo.NewClient(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create mongo connection: %w", err)
	}

	err = client.Connect(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongo: %w", err)
	}

	svc := &mongoDB{client: client, database: database}

	err = svc.ensureIndicies()
	if err != nil {
		return nil, fmt.Errorf("failed to create collection indicies: %w", err)
	}

	return svc, nil
}

func (s *mongoDB) ensureIndicies() error {
	ctx := context.Background()
	session, err := s.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)
	db := session.Client().Database(s.database)

	indexes := map[string][]mongo.IndexModel{
		collectionUsers: {
			mongo.IndexModel{
				Keys:    bson.M{"id": 1},
				Options: options.Index().SetName("users_idx"),
			},
			// TODO Indexes
		},
	}

	for collection, indices := range indexes {
		_, err := db.
			Collection(collection).
			Indexes().
			CreateMany(ctx, indices)
		if err != nil {
			return fmt.Errorf("unable to create database indices for collection %q: %w", collection, err)
		}
	}

	return nil
}

func (m *mongoDB) Create(ctx context.Context, user *users.User) error {
	session, err := m.client.StartSession()
	if err != nil {
		return fmt.Errorf("error creating session: %w", err)
	}

	defer session.EndSession(ctx)
	db := session.Client().Database(m.database)

	_, err = db.Collection(collectionUsers).
		UpdateOne(ctx,
			bson.M{"id": user.Id},
			bson.M{"$set": user},
			options.Update().SetUpsert(true))

	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	return nil
}

func (m *mongoDB) GetByID(ctx context.Context, id string) (*users.User, error) {
	res := m.client.Database(m.database).
		Collection(collectionUsers).
		FindOne(ctx, bson.M{"id": id})
	switch {
	case res.Err() == mongo.ErrNoDocuments:
		return nil, ErrNotFound
	case res.Err() != nil:
		return nil, fmt.Errorf("could not get service: %w", res.Err())
	}

	var s users.User
	err := res.Decode(&s)
	if err != nil {
		return nil, fmt.Errorf("could not decode service: %w", err)
	}
	return &s, nil
}

func (m *mongoDB) Patch(ctx context.Context, user *users.User) error {
	session, err := m.client.StartSession()
	if err != nil {
		return fmt.Errorf("error creating session: %w", err)
	}

	defer session.EndSession(ctx)
	db := session.Client().Database(m.database)

	updates := bson.M{}
	if user.GetFirstName() != "" {
		updates["firstname"] = user.GetFirstName()
	}
	if user.GetLastName() != "" {
		updates["lastname"] = user.GetLastName()
	}
	if user.GetNickname() != "" {
		updates["nickname"] = user.GetNickname()
	}
	if user.GetPassword() != "" {
		updates["password"] = user.GetPassword()
	}
	if user.GetEmail() != "" {
		updates["email"] = user.GetEmail()
	}
	if user.GetCountry() != "" {
		updates["country"] = user.GetCountry()
	}
	if user.GetCreatedAt() != nil {
		updates["createdat"] = user.GetCreatedAt()
	}
	if user.GetUpdatedAt() != nil {
		updates["updatedat"] = user.GetUpdatedAt()
	}

	_, err = db.Collection(collectionUsers).
		UpdateOne(ctx,
			bson.M{"id": user.Id},
			bson.M{"$set": updates})
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	return nil
}

func (m *mongoDB) Delete(ctx context.Context, id string) error {
	res, err := m.client.Database(m.database).
		Collection(collectionUsers).
		DeleteOne(ctx, bson.M{"id": id})
	if res.DeletedCount == 0 {
		return ErrNotFound
	}

	if err != nil {
		return fmt.Errorf("could not delete user by id: %w", err)
	}
	return nil
}

func (m *mongoDB) List(ctx context.Context, request *users.ListUsersRequest) ([]*users.User, error) {
	session, err := m.client.StartSession()
	if err != nil {
		return nil, err
	}

	defer session.EndSession(ctx)
	db := session.Client().Database(m.database)

	var filters []bson.M
	if request.Filters != nil {
		if request.GetFilters().GetEmail() != "" {
			filters = append(filters, bson.M{"email": request.GetFilters().GetEmail()})
		}
		if request.GetFilters().GetCountry() != "" {
			filters = append(filters, bson.M{"country": request.GetFilters().GetCountry()})
		}
	}

	sorting := bson.M{"createdat": 1}

	query := bson.M{}
	if len(filters) != 0 {
		query = bson.M{"$and": filters}
	}
	find, err := db.Collection(collectionUsers).Find(
		ctx,
		query,
		options.Find().
			SetSkip(int64(request.GetPagination().GetOffset())).
			SetLimit(int64(request.GetPagination().GetLimit())).
			SetSort(sorting),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to list users: %w", err)
	}

	var results []*users.User
	err = find.All(ctx, &results)
	if err != nil {
		return nil, fmt.Errorf("unable to find all results: %w", err)
	}

	return results, nil
}
