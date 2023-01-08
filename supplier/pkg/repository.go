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
	Get(string) (*Supplier, error)
	Create(Supplier) (*Supplier, error)
	List(page, perPage int64) ([]*Supplier, int64, error)
	Update(string, Supplier) (*Supplier, error)
	Delete(string) error
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

func (r *mongoRepository) List(page, perPage int64) ([]*Supplier, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	collection := r.database.Collection("suppliers")

	defer cancel()

	options := options.Find()
	options.SetLimit(perPage)
	options.SetSkip(page * perPage)

	total, err := collection.EstimatedDocumentCount(ctx)
	if err != nil {
		return nil, 0, err
	}

	result, err := collection.Find(ctx, bson.D{}, options)
	if err != nil {
		return nil, 0, err
	}

	var suppliers []*Supplier
	if err := result.All(ctx, &suppliers); err != nil {
		return nil, 0, err
	}

	return suppliers, total, nil
}

func (r *mongoRepository) Update(id string, data Supplier) (*Supplier, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	collection := r.database.Collection("suppliers")

	defer cancel()

	_, err := collection.ReplaceOne(ctx, bson.M{"_id": id}, data)
	if err != nil {
		return nil, err
	}

	return r.Get(id)
}

func (r *mongoRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	collection := r.database.Collection("suppliers")

	defer cancel()

	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
