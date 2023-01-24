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
	GetPaymentMethod(string) (*Method, error)
	CreatePaymentMethod(Method) (*Method, error)
	ListPaymentMethods() ([]*Method, error)
	UpdatePaymentMethod(string, Method) (*Method, error)
	DeletePaymentMethod(string) error

	CreatePaymentType(Type) (*Type, error)
	GetPaymentType(string) (*Type, error)
	ListPaymentTypes() ([]*Type, error)
	UpdatePaymentType(string, Type) (*Type, error)
	DeletePaymentType(string) error

	CreatePaymentCondition(Condition) (*Condition, error)
	GetPaymentCondition(string) (*Condition, error)
	ListPaymentConditions() ([]*Condition, error)
	UpdatePaymentCondition(string, Condition) (*Condition, error)
	DeletePaymentCondition(string) error

	CreateInvoice(Invoice) (*Invoice, error)
	ListInvoices(int64, int64) ([]*Invoice, int64, error)
	UpdateInvoice(string, Invoice) (*Invoice, error)
	GetInvoice(string) (*Invoice, error)
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

func (r *mongoRepository) CreatePaymentMethod(data Method) (*Method, error) {
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

func (r *mongoRepository) GetPaymentMethod(id string) (*Method, error) {
	collection := r.database.Collection("payment_methods")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	result := collection.FindOne(ctx, bson.M{"_id": id})
	if result.Err() != nil {
		return nil, result.Err()
	}

	var method *Method
	if err := result.Decode(&method); err != nil {
		return nil, err
	}

	return method, nil
}

func (r *mongoRepository) ListPaymentMethods() ([]*Method, error) {
	collection := r.database.Collection("payment_methods")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	methods := make([]*Method, 0)
	if err := cursor.All(ctx, &methods); err != nil {
		return nil, err
	}

	return methods, nil
}

func (r *mongoRepository) UpdatePaymentMethod(id string, data Method) (*Method, error) {
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

func (r *mongoRepository) CreatePaymentType(data Type) (*Type, error) {
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

func (r *mongoRepository) GetPaymentType(id string) (*Type, error) {
	collection := r.database.Collection("payment_types")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()
	result := collection.FindOne(ctx, bson.M{"_id": id})

	var paymentType *Type
	err := result.Decode(&paymentType)

	return paymentType, err
}

func (r *mongoRepository) ListPaymentTypes() ([]*Type, error) {
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

	paymentTypes := make([]*Type, 0)
	if err := result.All(ctx, &paymentTypes); err != nil {
		return nil, err
	}

	return paymentTypes, nil
}

func (r *mongoRepository) UpdatePaymentType(id string, data Type) (*Type, error) {
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

func (r *mongoRepository) DeletePaymentCondition(id string) error {
	collection := r.database.Collection("payment_conditions")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})

	return err
}

func (r *mongoRepository) CreateInvoice(data Invoice) (*Invoice, error) {
	collection := r.database.Collection("invoices")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	data.ID = primitive.NewObjectID().Hex()
	result, err := collection.InsertOne(ctx, data)

	if err != nil {
		return nil, err
	}

	return r.GetInvoice(result.InsertedID.(string))
}

func (r *mongoRepository) GetInvoice(id string) (*Invoice, error) {
	collection := r.database.Collection("invoices")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	result := collection.FindOne(ctx, bson.M{"_id": id})
	if result.Err() != nil {
		return nil, result.Err()
	}

	var invoice *Invoice
	return invoice, result.Decode(&invoice)
}

func (r *mongoRepository) ListInvoices(page, perPage int64) ([]*Invoice, int64, error) {
	collection := r.database.Collection("invoices")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

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

	invoices := make([]*Invoice, 0)
	return invoices, total, result.All(ctx, &invoices)
}

func (r *mongoRepository) UpdateInvoice(id string, data Invoice) (*Invoice, error) {
	collection := r.database.Collection("invoices")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()
	_, err := collection.ReplaceOne(ctx, bson.M{"_id": id}, data)

	if err != nil {
		return nil, err
	}

	return r.GetInvoice(id)
}
