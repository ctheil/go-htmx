package handlers

import (
	"fmt"
	"net/http"

	"github.com/ctheil/go-htmx/api"
	"github.com/gorilla/schema"
	log "github.com/sirupsen/logrus"
)

func GetTodos(w http.ResponseWriter, r *http.Request) {
	var params = api.TodosResponse{}
	var decoder *schema.Decoder = schema.NewDecoder()
	fmt.Println("/todos")

	if err := decoder.Decode(&params, r.URL.Query()); err != nil {
		log.Error(err)
		api.InternalErrorHandler(w)
		return
	}

	// var database *tools.DatabaseInterface
	// database, err = tools.NewDatabase()

	//
	// var todoDetails *tools.TodoDetails
	// todoDetails = (*database).GetUserTodos(params.Username)
	// if todoDetails == nil {
	// 	log.Error(err)
	// 	api.InternalErrorHandler(w)
	// }
	//
	// var response = api.TodosResponse{
	// 	Todos: (*todoDetails).Todos,
	// 	Code:  http.StatusOK,
	// }
	// w.Header().Set("Content-Type", "application/json")
	// err = json.NewEncoder(w).Encode(response)
	// if err != nil {
	// 	log.Error(err)
	// 	api.InternalErrorHandler(w)
	// 	return
	// }
	//
}
