package database

import (
	"context"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"github.com/zsoltggs/golang-example/services/validation-service/pkg/validationmodels"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	collectionSchemas = "schemas"
)

type mongoDB struct {
	client   *mongo.Client
	database string
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
		collectionSchemas: {
			mongo.IndexModel{
				Keys:    bson.M{"id": 1},
				Options: options.Index().SetName("schema_id_idx"),
			},
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

func (m *mongoDB) UpsertSchema(ctx context.Context, req *validationmodels.UpsertSchemaRequest) error {
	session, err := m.client.StartSession()
	if err != nil {
		return fmt.Errorf("error creating session: %w", err)
	}

	defer session.EndSession(ctx)
	db := session.Client().Database(m.database)

	_, err = db.Collection(collectionSchemas).
		UpdateOne(ctx,
			bson.M{"id": req.SchemaID},
			bson.M{"$set": bson.M{"schema": req.Schema}},
			options.Update().SetUpsert(true))

	if err != nil {
		return fmt.Errorf("error upserting schema: %w", err)
	}
	return nil
}

type mongoSchema struct {
	Schema string `json:"schema"`
}

func (m *mongoDB) GetSchema(ctx context.Context, ID string) (string, error) {
	res := m.client.Database(m.database).
		Collection(collectionSchemas).
		FindOne(ctx, bson.M{"id": ID})
	switch {
	case res.Err() == mongo.ErrNoDocuments:
		return "", ErrNotFound
	case res.Err() != nil:
		return "", fmt.Errorf("could not get schema: %w", res.Err())
	}

	var s mongoSchema
	err := res.Decode(&s)
	if err != nil {
		return "", fmt.Errorf("could not decode schema: %w", err)
	}
	return s.Schema, nil
}

func (s *mongoDB) Health() error {
	return s.client.Ping(context.Background(), nil)
}
