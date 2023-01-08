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
	Get(string) (*Supplier, error)
	Create(Supplier) (*Supplier, error)
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

func (r *mongoRepository) Create(data Supplier) (*Supplier, error) {
	ctx := context.Background()
	collection := r.client.Database("customer").Collection("customers")

	data.ID = primitive.NewObjectID().Hex()
	result, err := collection.InsertOne(ctx, data)

	if err != nil {
		return nil, err
	}

	return r.Get(result.InsertedID.(string))
}

func (r *mongoRepository) Get(id string) (*Supplier, error) {
	var customer *Supplier

	ctx := context.Background()
	collection := r.client.Database("customer").Collection("customers")
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&customer)

	return customer, err
}
