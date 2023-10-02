package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

func NewStore(dataSourceName string) (*Store, error) {
	db, err := sqlx.Open("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("Error opening database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Error pinging database: %w", err)
	}
	fmt.Println(db)

	return &Store{
		UserStore:    NewUserStore(db),
		SessionStore: NewSessionStore(db),
		TodoStore:    NewTodoStore(db),
	}, nil

}

type Store struct {
	// todos.UserStore
	// todos.SessionStore
	// todos.TodoStore
	*UserStore
	*SessionStore
	*TodoStore
}
