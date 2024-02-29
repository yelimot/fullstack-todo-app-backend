package repository

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/yelimot/fullstack-todo-app-backend/pkg/model"
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

var _ Repository = (*repositoryImpl)(nil)

type repositoryImpl struct {
	mtx   sync.Mutex
	todos []*model.Todo
	db    *os.File
}

func New(db *os.File) (Repository, error) {
	var todos []*model.Todo
	dec := json.NewDecoder(db)
	if err := dec.Decode(&todos); err != nil && err != io.EOF {
		return nil, err
	}

	return &repositoryImpl{
		db:    db,
		todos: todos,
	}, nil
}

// Create a new todo
func (r *repositoryImpl) Create(todo *model.Todo) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	var uniqueId int = int(uuid.New().ID())
	todo.ID = uniqueId

	r.todos = append(r.todos, todo)

	if err := r.updateDb(); err != nil {
		return err
	}

	return nil
}

// Get a todo by id
func (r *repositoryImpl) Get(id int) (*model.Todo, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	for _, todo := range r.todos {
		if todo.ID == id {
			return todo, nil
		}
	}

	return nil, errors.New("todo not found")
}

// Get all todos
func (r *repositoryImpl) GetAll(filter string, sorting model.Sorting, pagination model.Pagination) ([]*model.Todo, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	var todos []*model.Todo

	// step 1: filter todos
	for _, todo := range r.todos {
		if filter == "" || strings.Contains(strings.ToLower(todo.Title), strings.ToLower(filter)) || strings.Contains(strings.ToLower(todo.Description), strings.ToLower(filter)) {
			todos = append(todos, todo)
		}
	}

	// step 2: sort todos
	todos = sortTodos(todos, sorting)

	// step 3: paginate todos
	if pagination.Limit > 0 {
		start := (pagination.Page - 1) * pagination.Limit
		if start > len(todos) {
			start = len(todos)
		}
		end := start + pagination.Limit
		if end > len(todos) {
			end = len(todos)
		}
		todos = todos[start:end]
	}
	return todos, nil
}

func sortTodos(todos []*model.Todo, sorting model.Sorting) []*model.Todo {
	// sort todos
	switch sorting.SortBy {
	case model.SortByID:
		if sorting.SortType == model.SortAscending {
			sort.Slice(todos, func(i, j int) bool {
				return todos[i].ID < todos[j].ID
			})
		} else {
			sort.Slice(todos, func(i, j int) bool {
				return todos[i].ID > todos[j].ID
			})
		}
	case model.SortByTitle:
		if sorting.SortType == model.SortAscending {
			sort.Slice(todos, func(i, j int) bool {
				return strings.ToLower(todos[i].Title) < strings.ToLower(todos[j].Title)
			})
		} else {
			sort.Slice(todos, func(i, j int) bool {
				return strings.ToLower(todos[i].Title) > strings.ToLower(todos[j].Title)
			})
		}
	case model.SortByDescription:
		if sorting.SortType == model.SortAscending {
			sort.Slice(todos, func(i, j int) bool {
				return strings.ToLower(todos[i].Description) < strings.ToLower(todos[j].Description)
			})
		} else {
			sort.Slice(todos, func(i, j int) bool {
				return strings.ToLower(todos[i].Description) > strings.ToLower(todos[j].Description)
			})
		}
	case model.SortByDueDate:
		if sorting.SortType == model.SortAscending {
			sort.Slice(todos, func(i, j int) bool {
				return todos[i].DueDate < todos[j].DueDate
			})
		} else {
			sort.Slice(todos, func(i, j int) bool {
				return todos[i].DueDate > todos[j].DueDate
			})
		}
	}

	return todos
}

// Update a todo
func (r *repositoryImpl) Update(todo *model.Todo) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	for i, t := range r.todos {
		if t.ID == todo.ID {
			r.todos[i] = todo
			if err := r.updateDb(); err != nil {
				return err
			}
			return nil
		}
	}

	return errors.New("todo not found")
}

// Delete a todo
func (r *repositoryImpl) Delete(id int) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	for i, t := range r.todos {
		if t.ID == id {
			r.todos = append(r.todos[:i], r.todos[i+1:]...)
			if err := r.updateDb(); err != nil {
				return err
			}
			return nil
		}
	}

	return errors.New("todo not found")
}

func (r *repositoryImpl) updateDb() error {

	err := r.db.Truncate(0)
	if err != nil {
		logrus.Infof("%v - 1", err)
		return err
	}

	_, err = r.db.Seek(0, 0)
	if err != nil {
		logrus.Infof("%v - 2", err)
		return err
	}

	enc := json.NewEncoder(r.db)
	if err := enc.Encode(r.todos); err != nil {
		return err
	}

	return nil
}

func (r *repositoryImpl) Shutdown() error {
	return r.db.Close()
}
