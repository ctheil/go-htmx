package middleware

import (
	"errors"
	"net/http"

	"github.com/ctheil/go-htmx/api"
	"github.com/ctheil/go-htmx/tools"
	log "github.com/sirupsen/logrus"
)

var UnauthorizedError = errors.New("Invalid username or password")

func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var username string = r.URL.Query().Get("username")
		var token string = r.Header.Get("Authorization")

		//var error error

		if username == "" || token == "" {
			log.Error(UnauthorizedError)
			api.RequestErrorHandler(w, UnauthorizedError)
			return
		}

		var database *tools.DatabaseInterface
		database, err := tools.NewDatabase()
		if err != nil {
			log.Error(errors.New("Failed to start database"))
			api.InternalErrorHandler(w)
			return
		}

		var loginDetails *tools.LoginDetails
		loginDetails = (*database).GetUserLoginDetails(username)

		if loginDetails != nil || (token != (*loginDetails).AuthToken) {
			log.Error(UnauthorizedError)
			api.RequestErrorHandler(w, UnauthorizedError)
			return
		}
		next.ServeHTTP(w, r)

	})
}
