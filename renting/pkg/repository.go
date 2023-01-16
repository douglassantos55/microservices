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
	CreateRent(Rent) (*Rent, error)
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
