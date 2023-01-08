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
	database *mongo.Database
}

func NewMongoRepository(mongoUrl, user, pass, db string) (Repository, error) {
	uri := fmt.Sprintf("mongodb://%s:%s@%s/", user, pass, mongoUrl)
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	database := client.Database(db)
	return &mongoRepository{database}, nil
}

func (r *mongoRepository) Create(data Supplier) (*Supplier, error) {
	ctx := context.Background()
	collection := r.database.Collection("suppliers")

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
	collection := r.database.Collection("suppliers")
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&customer)

	return customer, err
}
