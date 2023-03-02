package Repository

import (
	"context"
	"errors"
	"fmt"

	Model "myapp/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	collection  *mongo.Collection
	mongoclient *mongo.Client
	err         error
)

func OpenTable() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoconn := options.Client().ApplyURI("mongodb://localhost:27017")
	mongoclient, err = mongo.Connect(ctx, mongoconn)
	if err != nil {
		return fmt.Errorf("fail with connecting to mongo: %w", err)
	}

	err = mongoclient.Ping(ctx, readpref.Primary())
	if err != nil {
		return fmt.Errorf("error while trying to ping mongo: %w", err)
	}

	collection = mongoclient.Database("person").Collection("persons")

	return nil
}

func Disconnect() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = mongoclient.Disconnect(ctx); err != nil {
		panic(err)
	}
}

func Create(ctx context.Context, p Model.Person) (*interface{}, error) {
	options := bson.D{bson.E{Key: "firstName", Value: p.FirstName}, bson.E{Key: "lastName", Value: p.LastName}, bson.E{Key: "phone", Value: p.Phone}, bson.E{Key: "email", Value: p.Email}}
	res, err := collection.InsertOne(ctx, options)
	if err != nil {
		return nil, err
	}
	id := res.InsertedID

	return &id, nil
}

func ReadOne(ctx context.Context, id string) (*Model.Person, error) {
	var person Model.Person

	pid, _ := primitive.ObjectIDFromHex(id)

	err = collection.FindOne(ctx, bson.D{{Key: "_id", Value: pid}}).Decode(&person)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("record does not exist: %w", err)
	} else if err != nil {
		return nil, err
	}
	return &person, nil
}

func Read(ctx context.Context) ([]*Model.Person, error) {
	var persons []*Model.Person
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	for cursor.Next(ctx) {
		var person Model.Person
		err := cursor.Decode(&person)
		if err != nil {
			return nil, err
		}
		persons = append(persons, &person)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	cursor.Close(ctx)

	return persons, nil
}

func Update(ctx context.Context, id string, p Model.Person) error {

	pid, _ := primitive.ObjectIDFromHex(id)

	filter := bson.D{bson.E{Key: "_id", Value: pid}}
	update := bson.D{bson.E{Key: "$set", Value: bson.D{bson.E{Key: "firstName", Value: p.FirstName}, bson.E{Key: "lastName", Value: p.LastName}, bson.E{Key: "phone", Value: p.Phone}, bson.E{Key: "email", Value: p.Email}}}}
	result, _ := collection.UpdateOne(ctx, filter, update)
	if result.MatchedCount != 1 {
		return errors.New("no matched document found for update")
	}
	return nil
}

func Delete(ctx context.Context, id string) error {

	pid, _ := primitive.ObjectIDFromHex(id)

	filter := bson.D{bson.E{Key: "_id", Value: pid}}
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount != 1 {
		return errors.New("no matched document found for delete")
	}
	return nil
}
