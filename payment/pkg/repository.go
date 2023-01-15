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
	GetPaymentMethod(string) (*PaymentMethod, error)
	CreatePaymentMethod(PaymentMethod) (*PaymentMethod, error)
	ListPaymentMethods() ([]*PaymentMethod, error)
	UpdatePaymentMethod(string, PaymentMethod) (*PaymentMethod, error)
	DeletePaymentMethod(string) error

	CreatePaymentType(PaymentType) (*PaymentType, error)
	GetPaymentType(string) (*PaymentType, error)
	ListPaymentTypes() ([]*PaymentType, error)
	UpdatePaymentType(string, PaymentType) (*PaymentType, error)
	DeletePaymentType(string) error

	CreatePaymentCondition(Condition) (*Condition, error)
	GetPaymentCondition(string) (*Condition, error)
	ListPaymentConditions() ([]*Condition, error)
	UpdatePaymentCondition(string, Condition) (*Condition, error)
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

func (r *mongoRepository) CreatePaymentMethod(data PaymentMethod) (*PaymentMethod, error) {
	collection := r.database.Collection("payment_methods")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	data.ID = primitive.NewObjectID().Hex()
	result, err := collection.InsertOne(ctx, data)

	if err != nil {
		return nil, err
	}

	return r.GetPaymentMethod(result.InsertedID.(string))
}

func (r *mongoRepository) GetPaymentMethod(id string) (*PaymentMethod, error) {
	collection := r.database.Collection("payment_methods")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	result := collection.FindOne(ctx, bson.M{"_id": id})
	if result.Err() != nil {
		return nil, result.Err()
	}

	var method *PaymentMethod
	if err := result.Decode(&method); err != nil {
		return nil, err
	}

	return method, nil
}

func (r *mongoRepository) ListPaymentMethods() ([]*PaymentMethod, error) {
	collection := r.database.Collection("payment_methods")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	methods := make([]*PaymentMethod, 0)
	if err := cursor.All(ctx, &methods); err != nil {
		return nil, err
	}

	return methods, nil
}

func (r *mongoRepository) UpdatePaymentMethod(id string, data PaymentMethod) (*PaymentMethod, error) {
	collection := r.database.Collection("payment_methods")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	_, err := collection.ReplaceOne(ctx, bson.M{"_id": id}, data)
	if err != nil {
		return nil, err
	}

	return r.GetPaymentMethod(id)
}

func (r *mongoRepository) DeletePaymentMethod(id string) error {
	collection := r.database.Collection("payment_methods")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})

	return err
}

func (r *mongoRepository) CreatePaymentType(data PaymentType) (*PaymentType, error) {
	collection := r.database.Collection("payment_types")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	data.ID = primitive.NewObjectID().Hex()
	result, err := collection.InsertOne(ctx, data)

	if err != nil {
		return nil, err
	}

	return r.GetPaymentType(result.InsertedID.(string))
}

func (r *mongoRepository) GetPaymentType(id string) (*PaymentType, error) {
	collection := r.database.Collection("payment_types")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()
	result := collection.FindOne(ctx, bson.M{"_id": id})

	var paymentType *PaymentType
	err := result.Decode(&paymentType)

	return paymentType, err
}

func (r *mongoRepository) ListPaymentTypes() ([]*PaymentType, error) {
	collection := r.database.Collection("payment_types")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()
	result, err := collection.Find(ctx, bson.D{})

	if err != nil {
		return nil, err
	}

	if result.Err() != nil {
		return nil, result.Err()
	}

	paymentTypes := make([]*PaymentType, 0)
	if err := result.All(ctx, &paymentTypes); err != nil {
		return nil, err
	}

	return paymentTypes, nil
}

func (r *mongoRepository) UpdatePaymentType(id string, data PaymentType) (*PaymentType, error) {
	collection := r.database.Collection("payment_types")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	_, err := collection.ReplaceOne(ctx, bson.M{"_id": id}, data)
	if err != nil {
		return nil, err
	}

	return r.GetPaymentType(id)
}

func (r *mongoRepository) DeletePaymentType(id string) error {
	collection := r.database.Collection("payment_types")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})

	return err
}

func (r *mongoRepository) CreatePaymentCondition(data Condition) (*Condition, error) {
	collection := r.database.Collection("payment_conditions")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	data.ID = primitive.NewObjectID().Hex()
	result, err := collection.InsertOne(ctx, data)

	if err != nil {
		return nil, err
	}

	return r.GetPaymentCondition(result.InsertedID.(string))
}

func (r *mongoRepository) GetPaymentCondition(id string) (*Condition, error) {
	collection := r.database.Collection("payment_conditions")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()
	result := collection.FindOne(ctx, bson.M{"_id": id})

	if result.Err() != nil {
		return nil, result.Err()
	}

	var condition *Condition
	err := result.Decode(&condition)

	return condition, err
}

func (r *mongoRepository) ListPaymentConditions() ([]*Condition, error) {
	collection := r.database.Collection("payment_conditions")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()
	result, err := collection.Find(ctx, bson.D{})

	if err != nil {
		return nil, err
	}

	if result.Err() != nil {
		return nil, result.Err()
	}

	conditions := make([]*Condition, 0)
	err = result.All(ctx, &conditions)

	return conditions, err
}

func (r *mongoRepository) UpdatePaymentCondition(id string, data Condition) (*Condition, error) {
	collection := r.database.Collection("payment_conditions")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()
	_, err := collection.ReplaceOne(ctx, bson.M{"_id": id}, data)

	if err != nil {
		return nil, err
	}

	return r.GetPaymentCondition(id)
}
