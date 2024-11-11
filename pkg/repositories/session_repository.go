package repository

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/open-wm/blockehr/pkg/domain/entities"
	"github.com/open-wm/blockehr/pkg/interfaces"
)

type sessionRepository struct {
	db *sql.DB
}

// TODO: Test this with something like this:
// https://medium.easyread.co/unit-test-sql-in-golang-5af19075e68e
func NewSessionRepository(db *sql.DB) interfaces.SessionRepository {
	return &sessionRepository{db}
}

func (r *sessionRepository) GetSession(sessionToken string) *entities.Session {
	session := &entities.Session{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(ctx, "SELECT id, user_id, token, updated_at, created_at FROM sessions WHERE token = ?", sessionToken).
		Scan(&session.ID, &session.UserID, &session.Token, &session.UpdatedAt, &session.CreatedAt)
	if err != nil {
		log.Println(err)
		return nil
	}

	return session
}

func (r *sessionRepository) DeleteSession(sessionToken string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(ctx, "DELETE FROM sessions WHERE token = ?", sessionToken)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (r *sessionRepository) SaveSession(session *entities.Session) (*entities.Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if session.ID != 0 {
		_, err := r.db.ExecContext(ctx, "UPDATE sessions SET user_id = ?, token = ?, updated_at = ? WHERE id = ?",
			session.UserID, session.Token, session.UpdatedAt, session.ID)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return session, nil
	}

	_, err := r.db.ExecContext(ctx, "INSERT INTO sessions (user_id, token, created_at, updated_at) VALUES (?, ?, ?, ?)",
		session.UserID, session.Token, session.CreatedAt, session.UpdatedAt)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return session, nil
}
