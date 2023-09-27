package tools

import (
	"time"

	"github.com/ctheil/go-htmx/api"
	"github.com/google/uuid"
)

type mockDB struct{}

var mockLoginDetails = map[string]LoginDetails{
	"alex": {
		AuthToken: "123ABC",
		Username:  "alex",
	},
	"jason": {
		AuthToken: "456DEF",
		Username:  "jason",
	},
	"marie": {
		AuthToken: "789GHI",
		Username:  "marie",
	},
}

var mockTodoDetails = map[string]TodoDetails{
	"alex": {
		Username: "alex",
		Todos: []api.Todo{
			api.Todo{Name: "todo1",
				Complete: false,
				Id:       uuid.MustParse("bcf6c2e3-8855-4458-99a8-cbc4a1a84886"),
			}},
	},
}

func (d *mockDB) GetUserLoginDetails(username string) *LoginDetails {
	time.Sleep(time.Second * 1)

	var clientData = LoginDetails{}
	clientData, ok := mockLoginDetails[username]
	if !ok {
		return nil
	}

	return &clientData

}
func (d *mockDB) GetUserTodos(username string) *TodoDetails {

	time.Sleep(time.Second * 1)

	var clientData = TodoDetails{}
	clientData, ok := mockTodoDetails[username]
	if !ok {
		return nil
	}

	return &clientData
}

func (d *mockDB) SetupDatabase() error {
	return nil
}
