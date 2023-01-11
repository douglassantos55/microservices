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
	List(page, perPage int) ([]*Equipment, int, error)
	Update(string, Equipment) (*Equipment, error)
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

func (r *mongoRepository) List(page, perPage int) ([]*Equipment, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	collection := r.database.Collection("equipment")

	defer cancel()

	total, err := collection.EstimatedDocumentCount(ctx)
	if err != nil {
		return nil, 0, err
	}

	options := options.Find()
	options.SetLimit(int64(perPage))
	options.SetSkip(int64(page * perPage))

	result, err := collection.Find(ctx, bson.D{}, options)
	if err != nil {
		return nil, 0, err
	}

	var equipment []*Equipment
	if err := result.All(ctx, &equipment); err != nil {
		return nil, 0, err
	}

	return equipment, int(total), nil
}

func (r *mongoRepository) Update(id string, data Equipment) (*Equipment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	collection := r.database.Collection("equipment")

	defer cancel()

	_, err := collection.ReplaceOne(ctx, bson.M{"_id": id}, data)
	if err != nil {
		return nil, err
	}

	return r.Get(id)
}
