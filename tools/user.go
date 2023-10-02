package tools

import (
	"fmt"

	todos "github.com/ctheil/go-htmx"
	"github.com/jmoiron/sqlx"
)

type UserStore struct {
	*sqlx.DB
}

func (s *UserStore) GetUser(email string) (todos.User, error) {
	var u todos.User
	if err := s.Get(&u, `SELECT * FROM users WHERE email = $1`, email); err != nil {
		return todos.User{}, fmt.Errorf("No user found")
	}
	return u, nil
}
