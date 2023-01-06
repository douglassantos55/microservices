package pkg

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	Create(Customer) (*Customer, error)
}

type mongoRepository struct {
	client *mongo.Client
}

func NewMongoRepository(mongoUrl string, user, pass string) (Repository, error) {
	uri := fmt.Sprintf("mongodb://%s:%s@%s/", user, pass, mongoUrl)
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return &mongoRepository{client}, nil
}

func (r *mongoRepository) Create(data Customer) (*Customer, error) {
	ctx := context.Background()
	collection := r.client.Database("customer").Collection("customers")

	data.ID = primitive.NewObjectID()
	result, err := collection.InsertOne(ctx, data)

	if err != nil {
		return nil, err
	}

	var customer *Customer
	filter := bson.M{"_id": result.InsertedID.(primitive.ObjectID)}
	collection.FindOne(ctx, filter).Decode(&customer)

	return customer, nil
}
