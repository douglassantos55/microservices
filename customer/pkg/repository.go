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
	List(curPage, perPage int64) ([]*Customer, int64, error)
	Get(id string) (*Customer, error)
	Create(Customer) (*Customer, error)
	Update(string, Customer) (*Customer, error)
	Delete(string) error
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

func (r *mongoRepository) List(curPage, perPage int64) ([]*Customer, int64, error) {
	collection := r.client.Database("customer").Collection("customers")

	opts := options.Find()
	opts.SetLimit(perPage)
	opts.SetSkip(curPage * perPage)

	ctx := context.Background()
	total, err := collection.EstimatedDocumentCount(ctx)
	if err != nil {
		return nil, 0, err
	}

	cursor, err := collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, 0, err
	}

	var customers []*Customer
	if err := cursor.All(ctx, &customers); err != nil {
		return nil, 0, err
	}

	return customers, total, nil
}

func (r *mongoRepository) Create(data Customer) (*Customer, error) {
	ctx := context.Background()
	collection := r.client.Database("customer").Collection("customers")

	data.ID = primitive.NewObjectID().Hex()
	result, err := collection.InsertOne(ctx, data)

	if err != nil {
		return nil, err
	}

	var customer *Customer
	filter := bson.M{"_id": result.InsertedID.(string)}
	collection.FindOne(ctx, filter).Decode(&customer)

	return customer, nil
}

func (r *mongoRepository) Update(id string, customer Customer) (*Customer, error) {
	ctx := context.Background()
	collection := r.client.Database("customer").Collection("customers")

	_, err := collection.ReplaceOne(ctx, bson.M{"_id": id}, customer)
	if err != nil {
		return nil, err
	}

	return r.Get(id)
}

func (r *mongoRepository) Get(id string) (*Customer, error) {
	var customer *Customer

	ctx := context.Background()
	collection := r.client.Database("customer").Collection("customers")
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&customer)

	return customer, err
}

func (r *mongoRepository) Delete(id string) error {
	collection := r.client.Database("customer").Collection("customers")
	_, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}
