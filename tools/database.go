package tools

import (
	"github.com/ctheil/go-htmx/api"
	log "github.com/sirupsen/logrus"
)

// database collections

type LoginDetails struct {
	Username  string
	AuthToken string
}

type TodoDetails struct {
	Username string
	Todos    []api.Todo
}

// user interface becuase allows us to easily swap out databases and type of databases
type DatabaseInterface interface {
	GetUserLoginDetails(username string) *LoginDetails
	GetUserTodos(username string) *TodoDetails
	SetupDatabase() error
}

func NewDatabase() (*DatabaseInterface, error) {
	var database DatabaseInterface = &mockDB{}

	var err error = database.SetupDatabase()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &database, nil

}
