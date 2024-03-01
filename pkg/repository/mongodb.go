package repository

import (
	"context"
	"errors"

	"github.com/yelimot/fullstack-todo-app-backend/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	collection *mongo.Collection
}

var _ Repository = (*MongoRepository)(nil)

func NewMongoRepository(client *mongo.Client, databaseName, collectionName string) (*MongoRepository, error) {
	collection := client.Database(databaseName).Collection(collectionName)
	return &MongoRepository{collection: collection}, nil
}

// Create method using MongoDB
func (r *MongoRepository) Create(todo *model.Todo) error {
	_, err := r.collection.InsertOne(context.Background(), todo)
	return err
}

func (r *MongoRepository) Get(id int) (*model.Todo, error) {
	filter := bson.M{"id": id}
	result := r.collection.FindOne(context.Background(), filter)
	if result.Err() != nil {
		return nil, result.Err()
	}
	var todo model.Todo
	err := result.Decode(&todo)
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (r *MongoRepository) GetAll(filterS string, sorting model.Sorting, pagination model.Pagination) ([]*model.Todo, error) {
	// TODO: Perhaps nice to accept default parameters (or query parameters may be optional)?
	// Define a filter based on the provided string
	// You can customize this filter based on your requirements
	filter := bson.M{}
	if filterS != "" {
		filter["$or"] = []bson.M{
			{"title": bson.M{"$regex": primitive.Regex{Pattern: filterS, Options: "i"}}},
			{"description": bson.M{"$regex": primitive.Regex{Pattern: filterS, Options: "i"}}},
		}
	}

	// Define options for sorting and pagination
	options := options.Find()
	if sorting.SortBy == model.SortByID {
		options.SetSort(bson.D{{Key: "id", Value: sorting.SortType}})
	} else if sorting.SortBy == model.SortByTitle {
		options.SetSort(bson.D{{Key: "title", Value: sorting.SortType}})
	} else if sorting.SortBy == model.SortByDescription {
		options.SetSort(bson.D{{Key: "description", Value: sorting.SortType}})
	} else if sorting.SortBy == model.SortByDueDate {
		options.SetSort(bson.D{{Key: "dueDate", Value: sorting.SortType}})
	}

	if pagination.Limit > 0 {
		options.SetLimit(int64(pagination.Limit))
		options.SetSkip(int64((pagination.Page - 1) * pagination.Limit))
	}

	// Perform the find operation
	cursor, err := r.collection.Find(context.Background(), filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	// Iterate over the cursor and decode documents into []*model.Todo
	var todos []*model.Todo
	for cursor.Next(context.Background()) {
		var todo model.Todo
		if err := cursor.Decode(&todo); err != nil {
			return nil, err
		}
		todos = append(todos, &todo)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func (r *MongoRepository) Update(todo *model.Todo) error {
	filter := bson.M{"id": todo.ID}
	update := bson.M{"$set": bson.M{
		"title":       todo.Title,
		"description": todo.Description,
		"dueDate":     todo.DueDate,
	}}

	result, err := r.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	if result.ModifiedCount == 0 {
		return errors.New("todo not found")
	}
	return nil
}

func (r *MongoRepository) Delete(id int) error {
	filter := bson.M{"id": id}
	_, err := r.collection.DeleteOne(context.Background(), filter)
	return err
}

func (r *MongoRepository) Shutdown() error {
	// Disconnect from the MongoDB client
	err := r.collection.Database().Client().Disconnect(context.Background())
	if err != nil {
		return err
	}
	return nil
}
