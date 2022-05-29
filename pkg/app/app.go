package app

import "github.com/yelimot/fullstack-todo-app/todo-app-backend/pkg/repository"

type App struct {
	Repository repository.Repository
}

func New(repository repository.Repository) *App {
	return &App{
		Repository: repository,
	}
}
