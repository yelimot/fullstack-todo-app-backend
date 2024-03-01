package repository

import (
	"errors"
	"os"

	"github.com/yelimot/fullstack-todo-app-backend/pkg/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	// Create a new todo
	Create(todo *model.Todo) error
	// Get a todo by id
	Get(id int) (*model.Todo, error)
	// Get all todos
	GetAll(filter string, sorting model.Sorting, pagination model.Pagination) ([]*model.Todo, error)
	// Update a todo
	Update(todo *model.Todo) error
	// Delete a todo
	Delete(id int) error

	Shutdown() error
}

func New(client interface{}) (Repository, error) {
	switch client := client.(type) {
	case *os.File:
		return NewJSONRepository(client)
	case *mongo.Client:
		return NewMongoRepository(client, "todo", "todos")
	default:
		return nil, errors.New("unsupported client type")
	}
}
