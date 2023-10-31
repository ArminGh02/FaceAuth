package mongodb

import (
	"context"
	"errors"

	"github.com/ArminGh02/go-auth-system/internal/model"
	"github.com/ArminGh02/go-auth-system/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DB struct {
	client *mongo.Client
	coll   *mongo.Collection
}

const (
	Database   = "auth_system"
	Collection = "users"
)

func New(ctx context.Context, uri string) (*DB, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	{
		err := client.Ping(ctx, readpref.Primary())
		if err != nil {
			return nil, err
		}
	}

	coll := client.Database(Database).Collection(Collection)

	{
		_, err := coll.Indexes().CreateMany(
			ctx,
			[]mongo.IndexModel{
				{
					Keys:    bson.D{{Key: "email", Value: 1}},
					Options: options.Index().SetUnique(true),
				},
				{
					Keys:    bson.D{{Key: "national_id", Value: 1}},
					Options: options.Index().SetUnique(true),
				},
			},
		)
		if err != nil {
			return nil, err
		}
	}

	return &DB{
		client: client,
		coll:   coll,
	}, nil
}

func (db *DB) Insert(ctx context.Context, user *model.User) error {
	_, err := db.coll.InsertOne(ctx, user)
	return err
}

func (db *DB) Update(ctx context.Context, user *model.User) error {
	_, err := db.coll.UpdateOne(ctx, bson.D{{"national_id", user.NationalID}}, bson.D{{"$set", user}})
	if errors.Is(err, mongo.ErrNoDocuments) {
		return repository.ErrNotFound
	}
	return err
}

func (db *DB) GetByNationalID(ctx context.Context, nationalID string) (*model.User, error) {
	var user model.User
	err := db.coll.FindOne(ctx, bson.D{{"national_id", nationalID}}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, repository.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *DB) Close(ctx context.Context) error {
	return db.client.Disconnect(ctx)
}
