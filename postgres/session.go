package postgres

import (
	"fmt"
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

func (s *SessionStore) StartSession(u *todos.User) (todos.Session, error) {

	// go routine to check for all tokens associated with email and clean up?
	var sess todos.Session
	uuid := uuid.New()
	// should check for remember me
	expiresAt := time.Now().Add(120 * time.Hour) // 2hrs
	sess = todos.Session{
		Email:  u.Email,
		Expiry: expiresAt,
		Token:  uuid,
	}

	_, err := s.NamedExec(`INSERT INTO sessions (token, email, expiry)VALUES (:token, :email, :expiry) RETURNING *`, &sess)
	if err != nil {
		return todos.Session{}, err
	}
	return sess, nil
}

// go routine to clean up expired tokens?

func (s *SessionStore) IsExpired(sess *todos.Session) bool {

	fmt.Println(sess.Token)
	if err := s.Get(&sess, `SELECT * FROM sessions WHERE token = $1`, sess.Token); err != nil {
		fmt.Println("ERROR EXPIRED: ", err.Error())

		return true
	}

	return sess.Expiry.Before(time.Now())
}

func (s *SessionStore) GetSession(sessionId uuid.UUID) (todos.Session, error) {
	var sess todos.Session

	if err := s.Get(&sess, `SELECT * FROM sessions WHERE token = $1`, sessionId); err != nil {
		return todos.Session{}, fmt.Errorf("Could not find session: %w", err)
	}

	return sess, nil
}
