package tools

import (
	"time"

	todos "github.com/ctheil/go-htmx"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func NewSessionStore(db *sqlx.DB) *SessionStore {
	return &SessionStore{
		DB: db,
	}
}

type SessionStore struct {
	*sqlx.DB
}

type Session struct {
	Email  string
	Expiry time.Time
	Token  uuid.UUID
}

func (s *SessionStore) StartSession(u todos.User) (Session, error) {
	var sess Session
	uuid := uuid.New()
	// should check for remember me
	expiresAt := time.Now().Add(120 * time.Second) // 2min
	sess = Session{
		Email:  u.Email,
		Expiry: expiresAt,
		Token:  uuid,
	}

	if err := s.Get(sess, `INSERT INTO sessions VALUES ($1, $2, $3)`, sess.Token, sess.Email, sess.Expiry); err != nil {
		return Session{}, err
	}
	return sess, nil

}
