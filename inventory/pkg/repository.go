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
	Get(string) (*Equipment, error)
	Create(Equipment) (*Equipment, error)
}

type mongoRepository struct {
	database *mongo.Database
}

func NewMongoRepository(url, user, pass, db string) (Repository, error) {
	mongoUrl := fmt.Sprintf("mongodb://%s:%s@%s/", user, pass, url)
	options := options.Client().ApplyURI(mongoUrl)

	client, err := mongo.Connect(context.Background(), options)
	if err != nil {
		return nil, err
	}

	return &mongoRepository{client.Database(db)}, nil
}

func (r *mongoRepository) Create(data Equipment) (*Equipment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	collection := r.database.Collection("equipment")

	defer cancel()

	data.ID = primitive.NewObjectID().Hex()
	result, err := collection.InsertOne(ctx, data)

	if err != nil {
		return nil, err
	}
	return r.Get(result.InsertedID.(string))
}

func (r *mongoRepository) Get(id string) (*Equipment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	collection := r.database.Collection("equipment")

	defer cancel()

	result := collection.FindOne(ctx, bson.M{"_id": id})
	if result.Err() != nil {
		return nil, result.Err()
	}

	var equipment *Equipment
	err := result.Decode(&equipment)

	return equipment, err
}
