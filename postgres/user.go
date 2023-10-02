package postgres

import (
	"fmt"

	todos "github.com/ctheil/go-htmx"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func NewUserStore(db *sqlx.DB) *UserStore {
	return &UserStore{
		DB: db,
	}
}

type UserStore struct {
	*sqlx.DB
}

func (s *UserStore) GetUserByEmail(email string) (todos.User, error) {
	var u todos.User
	if err := s.Get(&u, `SELECT * FROM users WHERE email = $1`, email); err != nil {
		return todos.User{}, fmt.Errorf("No user found")
	}
	return u, nil
}

func (s *UserStore) GetUserById(id uuid.UUID) (todos.User, error) {
	var u todos.User
	if err := s.Get(&u, `SELECT * FROM users WHERE id = $1`, id); err != nil {
		return todos.User{}, fmt.Errorf("No user found")
	}
	return u, nil
}
func (s *UserStore) CreateUser(u *todos.User) error {
	// hash password
	if err := s.Get(u, `INSERT INTO users VALUES ($1, $2, $3, $4) RETURNING *`, u.ID, u.Email, u.Password, u.Name); err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	return nil
}
