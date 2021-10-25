package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type dbCollection map[string]*options.CreateCollectionOptions

var (
	validationAction = "error"
	validationLevel  = "strict"
)

// DB encapsulates the mongoDB connection properties.
type DB struct {
	Client      *mongo.Client
	DB          *mongo.Database
	collections []dbCollection
}

// New returns a configured *DB instance.
func New(client *mongo.Client, dbName string) *DB {
	database := client.Database(dbName)

	return &DB{
		Client:      client,
		DB:          database,
		collections: []dbCollection{userCollectionSchema},
	}
}

func (db *DB) addCollection(collection dbCollection) {
	db.collections = append(db.collections, collection)
}

// RegisterNewCollections adds new collections if any to the existing database instance.
func (db *DB) RegisterNewCollections() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	existingCollections, err := db.DB.ListCollectionNames(ctx, bson.D{}, options.ListCollections().SetNameOnly(true))
	if err != nil {
		return err
	}

first:
	for _, col := range db.collections {
		for name, opts := range col {
			for _, colName := range existingCollections {
				if colName == name {
					continue first
				}
			}

			err := db.DB.CreateCollection(ctx, name, opts)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// OpenConnection opens a mongoDB connection with the provided dsn.
func OpenConnection(dsn string) (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(dsn))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return client, nil
}
