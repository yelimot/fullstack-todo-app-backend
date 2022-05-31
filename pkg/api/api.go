package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/yelimot/fullstack-todo-app-backend/pkg/app"
)

// API configuration
type Config struct {
	Addr string `yaml:"addr"`
}

type API struct {
	Router *mux.Router

	config *Config

	app *app.App

	httpServer *http.Server
}

// New returns the api settings
func New(config *Config, app *app.App) (*API, error) {

	router := mux.NewRouter()
	api := &API{
		config: config,
		app:    app,
		Router: router,
	}

	// Endpoint for browser preflight requests
	api.Router.Methods("OPTIONS").HandlerFunc(api.corsMiddleware(api.preflightHandler))

	// Get All
	api.Router.HandleFunc("/api/v1/todos", api.corsMiddleware(api.logMiddleware(api.GetTodos))).Methods("GET")
	// Get By Id
	api.Router.HandleFunc("/api/v1/todos/{id}", api.corsMiddleware(api.logMiddleware(api.GetTodo))).Methods("GET")
	// Create
	api.Router.HandleFunc("/api/v1/todos", api.corsMiddleware(api.logMiddleware(api.AddTodo))).Methods("POST")
	// Update
	api.Router.HandleFunc("/api/v1/todos", api.corsMiddleware(api.logMiddleware(api.UpdateTodo))).Methods("PUT")
	// Delete
	api.Router.HandleFunc("/api/v1/todos/{id}", api.corsMiddleware(api.logMiddleware(api.DeleteTodo))).Methods("DELETE")

	return api, nil

}

func (a *API) preflightHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (a *API) Start() error {
	a.httpServer = &http.Server{
		Addr:    a.config.Addr,
		Handler: a.Router,
	}

	err := a.httpServer.ListenAndServe()
	if err != http.ErrServerClosed {
		return err
	} else {
		return err
	}
}

// Shutdown stops the server
func (a *API) Shutdown() error {

	a.app.Repository.Shutdown()

	// Shutdown HTTP server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := a.httpServer.Shutdown(ctx)
	if err != nil {
		return err
	}

	logrus.Info("Shutdown HTTP server...")
	return nil
}
