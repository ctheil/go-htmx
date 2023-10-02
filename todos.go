package todos

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Email    string    `db:"email"`
	Password string    `db:"password"`
	ID       uuid.UUID `db:"id"`
	Name     string    `db:"name"`
	Todos    []Todo
}
type Session struct {
	Email  string    `db:"email"`
	Token  uuid.UUID `db:"token"`
	Expiry time.Time `db:"expiry"`
}
type Todo struct {
	ID          uuid.UUID `db:"id"`
	UserID      uuid.UUID `db:"user_id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	Complete    bool      `db:"complete"`
	Due         time.Time `db:"due"`
}

type UserStore interface {
	GetUserById(id uuid.UUID) (User, error)
	GetUserByEmail(email string) (User, error)
	CreateUser(u *User) error
}
type SessionStore interface {
	StartSession(u *User) (Session, error)
	IsExpired(s *Session) bool
	GetSession(sessionID uuid.UUID) (Session, error)
}
type TodoStore interface {
	GetTodos(userId uuid.UUID) ([]Todo, error)
	CreateTodo(t *Todo) error
}

type Store interface {
	UserStore
	SessionStore
	TodoStore
}
