package postgres

import (
	"fmt"

	todos "github.com/ctheil/go-htmx"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func NewTodoStore(db *sqlx.DB) *TodoStore {
	return &TodoStore{
		DB: db,
	}
}

type TodoStore struct {
	*sqlx.DB
}

func (s *TodoStore) GetTodos(userId uuid.UUID) ([]todos.Todo, error) {
	// go routine to check for all tokens associated with email and clean up?
	var tt []todos.Todo

	if err := s.Select(&tt, `SELECT * FROM todos WHERE user_id = $1`, userId); err != nil {
		return []todos.Todo{}, fmt.Errorf("error getting user todos: %w", err)
	}
	return tt, nil

}

func (s *TodoStore) CreateTodo(t *todos.Todo) error {
	// hash password
	_, err := s.NamedExec(`INSERT INTO todos (id, title, description, complete, due, user_id) VALUES (:id, :title, :description, :complete, :due, :user_id)`, t)
	if err != nil {
		return fmt.Errorf("error creating todo: %w", err)
	}
	// if err := s.Get(t, `INSERT INTO todos VALUES ($1, $2, $3, $4, $5, $6) RETURNING *`, t.ID, t.Title, t.Description, t.Complete, t.UserID, t.Due); err != nil {
	// 	return fmt.Errorf("error creating todo: %w", err)
	// }
	return nil
}
func (s *TodoStore) EditTodo(userId uuid.UUID) (todos.Todo, error) {

	return todos.Todo{}, nil
}
func (s *TodoStore) DeleteTodo(userId uuid.UUID) error {
	return nil
}
