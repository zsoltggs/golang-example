package database

import (
	"context"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"github.com/zsoltggs/golang-example/pkg/pds"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	collectionPort = "ports"
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
		collectionPort: {
			mongo.IndexModel{
				Keys:    bson.M{"id": 1},
				Options: options.Index().SetName("portid_idx"),
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

func (s *mongoDB) UpsertPort(ctx context.Context, port *pds.Port) error {
	session, err := s.client.StartSession()
	if err != nil {
		return fmt.Errorf("error creating session: %w", err)
	}

	defer session.EndSession(ctx)
	db := session.Client().Database(s.database)

	_, err = db.Collection(collectionPort).
		UpdateOne(ctx,
			bson.M{"id": port.Id},
			bson.M{"$set": port},
			options.Update().SetUpsert(true))
	if err != nil {
		return fmt.Errorf("error updating port: %w", err)
	}

	return nil
}

func (s *mongoDB) GetPortsPaginated(ctx context.Context, pagination *pds.Pagination) ([]*pds.Port, error) {
	//TODO implement me
	panic("implement me")
}

func (s *mongoDB) GetPortByID(ctx context.Context, ID string) (*pds.Port, error) {
	session, err := s.client.StartSession()
	if err != nil {
		return nil, fmt.Errorf("unable to start session: %w", err)
	}

	defer session.EndSession(ctx)
	db := session.Client().Database(s.database)

	filter := bson.M{"id": ID}
	port := &pds.Port{}

	err = db.Collection(collectionPort).FindOne(ctx, filter).Decode(&port)
	switch {
	case err == mongo.ErrNoDocuments:
		return nil, ErrNotFound
	case err != nil:
		return nil, fmt.Errorf("cannot get by id : %w", err)
	}

	return port, nil
}
