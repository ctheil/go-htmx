package api

import (
	"github.com/google/uuid"

)

type Todo struct {
	Name string
	Complete bool
	Id uuid.UUID
}

type User struct {
	Name: string
	Password: string
}

type TodosResponse struct {
	Code int

	Todos []Todo
}

type Error struct {
	Code int

	Message string
}
