package todo

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
)

// Router handles HTTP request routing for the Todo REST API server
type Router struct {
	// Prefix is prepended to each route path
	Prefix string

	// Authentication validates that requests are from permitted users
	Authentication Authentication

	// Persistence is used to access Todos
	Persistence Persistence
}

// ServeHTTP implements the http.Handler interface
func (r Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	// end with an error for all not authenticated requests
	userId, err := r.Authentication.Authenticate(req)
	if err != nil {
		r.handleError(rw, req, err)
		return
	}
	log.Printf("Request from [%s]: %s %s", userId, req.Method, req.RequestURI)

	// handle
	// - POST and GET for /todo
	// - DELETE and GET for a path looking like /todo/<id>
	path := req.URL.Path
	todoPath := r.Prefix + "/todo"
	if path == todoPath {
		switch req.Method {
		case http.MethodPost:
			r.create(rw, req, userId)
			return
		case http.MethodGet:
			r.list(rw, req)
			return
		}
	} else if strings.HasPrefix(path, todoPath+"/") {
		id := path[len(todoPath)+1:]
		switch req.Method {
		case http.MethodDelete:
			r.delete(rw, req, id)
			return
		case http.MethodGet:
			r.get(rw, req, id)
			return
		}
	}

	// anything else, we don't now
	rw.WriteHeader(http.StatusNotFound)
	rw.Write([]byte("not found"))
}

func (r Router) create(rw http.ResponseWriter, req *http.Request, userId string) {

	// read Todo from JSON body of HTTP request
	var todo Todo
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&todo); err != nil {
		r.handleError(rw, req, err)
		return
	}

	// create Todo in Persistence
	todo.UserID = userId
	todoID, err := r.Persistence.Create(todo)
	if err != nil {
		r.handleError(rw, req, err)
		return
	}
	r.json(rw, req, map[string]string{"id": todoID})
}

func (r Router) list(rw http.ResponseWriter, req *http.Request) {
	todos, err := r.Persistence.List()
	if err != nil {
		r.handleError(rw, req, err)
		return
	}
	r.json(rw, req, todos)
}

func (r Router) delete(rw http.ResponseWriter, req *http.Request, todoID string) {
	err := r.Persistence.Delete(todoID)
	if err != nil {
		r.handleError(rw, req, err)
		return
	}
	r.json(rw, req, map[string]string{"id": todoID})
}

func (r Router) get(rw http.ResponseWriter, req *http.Request, todoID string) {
	todo, err := r.Persistence.Get(todoID)
	if err != nil {
		r.handleError(rw, req, err)
		return
	}
	r.json(rw, req, todo)
}

// json prints out a JSON HTTP response
func (r Router) json(rw http.ResponseWriter, req *http.Request, data interface{}) {
	rw.Header().Set("content-type", "application/json")
	err := json.NewEncoder(rw).Encode(data)
	if err != nil {
		r.handleError(rw, req, err)
		return
	}
}

// handleError prints out errors in the logs and lets the request fail
func (r Router) handleError(rw http.ResponseWriter, req *http.Request, err error) {
	log.Printf("Error in %s %s: %s", req.Method, req.URL, err)
	rw.Header().Set("content-type", "application/json")
	if errors.Is(NotAllowedError, err) {
		rw.WriteHeader(http.StatusForbidden)
		rw.Write([]byte(`{"error":"forbidden"}`))
	} else {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(`{"error":"internal server error"}`))
	}
}
