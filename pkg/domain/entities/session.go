package entities

import "time"

// Stores the user session
type Session struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	User      *User     `json:"user"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewSession(user *User, sessionToken string, expiresAt time.Time, loggedInAt string) Session {
	return Session{
		User:      user,
		Token:     sessionToken,
		CreatedAt: expiresAt,
	}
}

func (s *Session) IsExpired() bool { //verifica si la sesion ha expirado

	const SESSION_DURATION = 24 * time.Hour
	return s.CreatedAt.Add(SESSION_DURATION).Before(time.Now()) // si el tiempo de expiracion de la sesion es anterior al tiempo actual, TRUE
}
