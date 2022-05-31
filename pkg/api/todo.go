package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/yelimot/fullstack-todo-app-backend/pkg/api/response"
	"github.com/yelimot/fullstack-todo-app-backend/pkg/model"
)

func (a *API) GetTodo(w http.ResponseWriter, r *http.Request) {

	id, ok := mux.Vars(r)["id"]
	if !ok {
		err := errors.New("id is required")
		response.Errorf(w, r, err, http.StatusBadRequest, err.Error())
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		response.Errorf(w, r, err, http.StatusBadRequest, err.Error())
		return
	}
	todo, err := a.app.Repository.Get(idInt)
	if err != nil {
		response.Errorf(w, r, err, http.StatusInternalServerError, err.Error())
		return
	}

	response.Write(w, r, todo)
}

func (a *API) AddTodo(w http.ResponseWriter, r *http.Request) {
	todo := model.Todo{}

	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(&todo); err != nil {
		response.Errorf(w, r, err, http.StatusBadRequest, err.Error())
		return
	}

	if err := a.app.Repository.Create(&todo); err != nil {
		response.Errorf(w, r, err, http.StatusInternalServerError, err.Error())
		return
	}

	response.Write(w, r, "OK")
}

func (a *API) GetTodos(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	// filter
	filter := params.Get("filter")

	// sort by
	sortBy := params.Get("sortBy")
	var sortByEnum model.SortBy
	switch sortBy {
	case "id":
		sortByEnum = model.SortByID
	case "title":
		sortByEnum = model.SortByTitle
	case "description":
		sortByEnum = model.SortByDescription
	case "dueDate":
		sortByEnum = model.SortByDueDate
	default:
		sortByEnum = model.SortByID
	}

	// sort type
	sortType := params.Get("sortType")
	var sortTypeEnum model.SortType
	switch sortType {
	case "asc":
		sortTypeEnum = model.SortAscending
	case "desc":
		sortTypeEnum = model.SortDescending
	default:
		sortTypeEnum = model.SortAscending
	}

	// pagination
	limit := params.Get("limit")
	if limit == "" {
		limit = "10"
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		response.Errorf(w, r, err, http.StatusBadRequest, err.Error())
		return
	}

	page := params.Get("page")
	if page == "" {
		page = "1"
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		response.Errorf(w, r, err, http.StatusBadRequest, err.Error())
		return
	}

	sorting := model.Sorting{
		SortBy:   sortByEnum,
		SortType: sortTypeEnum,
	}
	pagination := model.Pagination{
		Page:  pageInt,
		Limit: limitInt,
	}

	todos, err := a.app.Repository.GetAll(filter, sorting, pagination)
	if err != nil {
		response.Errorf(w, r, err, http.StatusInternalServerError, err.Error())
		return
	}

	response.Write(w, r, todos)
}

func (a *API) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	todo := model.Todo{}

	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(&todo); err != nil {
		response.Errorf(w, r, err, http.StatusBadRequest, err.Error())
		return
	}

	_, err := a.app.Repository.Get(todo.ID)
	if err != nil {
		response.Errorf(w, r, err, http.StatusBadRequest, err.Error())
		return
	}

	if err := a.app.Repository.Update(&todo); err != nil {
		response.Errorf(w, r, err, http.StatusInternalServerError, err.Error())
		return
	}

	response.Write(w, r, "OK")
}

func (a *API) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	id, ok := mux.Vars(r)["id"]
	if !ok {
		err := errors.New("id is required")
		response.Errorf(w, r, err, http.StatusBadRequest, err.Error())
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		response.Errorf(w, r, err, http.StatusBadRequest, err.Error())
		return
	}

	_, err = a.app.Repository.Get(idInt)
	if err != nil {
		response.Errorf(w, r, err, http.StatusBadRequest, err.Error())
		return
	}

	if err := a.app.Repository.Delete(idInt); err != nil {
		response.Errorf(w, r, err, http.StatusInternalServerError, err.Error())
		return
	}

	response.Write(w, r, "OK")
}
