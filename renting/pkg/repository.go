package pkg

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	GetRent(id string) (*Rent, error)
	CreateRent(Rent) (*Rent, error)
	ListRents(page, perPage int64) ([]*Rent, int64, error)
	UpdateRent(id string, data Rent) (*Rent, error)
}

type mongoRepository struct {
	database *mongo.Database
}

func NewMongoRepository(url, user, pass, db string) (*mongoRepository, error) {
	options := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s/", user, pass, url))
	conn, err := mongo.Connect(context.Background(), options)

	if err != nil {
		return nil, err
	}

	database := conn.Database(db)
	return &mongoRepository{database}, nil
}

func (r *mongoRepository) CreateRent(data Rent) (*Rent, error) {
	collection := r.database.Collection("rents")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	data.ID = primitive.NewObjectID().Hex()
	result, err := collection.InsertOne(ctx, data)

	if err != nil {
		return nil, err
	}

	return r.GetRent(result.InsertedID.(string))
}

func (r *mongoRepository) ListRents(page, perPage int64) ([]*Rent, int64, error) {
	collection := r.database.Collection("rents")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	defer cancel()

	total, err := collection.EstimatedDocumentCount(ctx)
	if err != nil {
		return nil, 0, err
	}

	options := options.Find()
	options.SetLimit(perPage)
	options.SetSkip(page * perPage)

	result, err := collection.Find(ctx, bson.D{}, options)
	if err != nil {
		return nil, 0, err
	}

	if result.Err() != nil {
		return nil, 0, result.Err()
	}

	rents := make([]*Rent, 0)
	return rents, total, result.All(ctx, &rents)
}

func (r *mongoRepository) GetRent(id string) (*Rent, error) {
	collection := r.database.Collection("rents")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()
	result := collection.FindOne(ctx, bson.M{"_id": id})

	if result.Err() != nil {
		return nil, result.Err()
	}

	var rent *Rent
	if err := result.Decode(&rent); err != nil {
		return nil, err
	}

	return rent, nil
}

func (r *mongoRepository) UpdateRent(id string, data Rent) (*Rent, error) {
	collection := r.database.Collection("rents")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	defer cancel()

	if _, err := collection.ReplaceOne(ctx, bson.M{"_id": id}, data); err != nil {
		return nil, err
	}

	return r.GetRent(id)
}
