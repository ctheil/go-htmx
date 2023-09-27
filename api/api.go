package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type Todo struct {
	Name     string
	Complete bool
	Id       uuid.UUID
}

type User struct {
	Name     string
	Password string
}

type TodosResponse struct {
	Code int

	Todos []Todo
}

type Error struct {
	Code int

	Message string
}

func writeError(w http.ResponseWriter, message string, code int) {
	resp := Error{
		Code:    code,
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(resp)

}

var (
	RequestErrorHandler = func(w http.ResponseWriter, err error) {
		writeError(w, err.Error(), http.StatusBadRequest)
	}
	InternalErrorHandler = func(w http.ResponseWriter) {
		writeError(w, "An unexpected error occurred...", http.StatusInternalServerError)
	}
)
