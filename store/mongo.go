package store

import (
	"context"
	"time"

	"github.com/sing3demons/todoapi/todo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoStore struct {
	*mongo.Collection
}

func NewMongoStore(col *mongo.Collection) *MongoStore {
	return &MongoStore{Collection: col}
}

func (tx *MongoStore) New(todo *todo.Todo) error {
	_, err := tx.Collection.InsertOne(context.Background(), todo)
	return err
}

func (tx *MongoStore) List(todos *[]todo.Todo) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := tx.Collection.Find(ctx, bson.M{})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	return cursor.All(ctx, todos)
}

func (tx *MongoStore) Delete(todo *todo.Todo, id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return tx.Collection.FindOneAndDelete(ctx, id).Err()
}
