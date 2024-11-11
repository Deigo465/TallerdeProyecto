package entities

import (
	"testing"
	"time"
)

func TestNewSession(t *testing.T) {
	//given
	user := NewFakeUser()
	sessionToken := "ABC123"
	expiresAt := time.Now().Add(24 * time.Hour)
	loggedInAt := time.Now().Format(time.RFC3339)

	// When
	session := NewSession(&user, sessionToken, expiresAt, loggedInAt)

	// Then
	if session.User != &user {
		t.Errorf("Expected User to be %v, got %v", user, session.User)
	}
	if session.Token != sessionToken {
		t.Errorf("Expected Token to be %s, got %s", sessionToken, session.Token)
	}
	if !session.CreatedAt.Equal(expiresAt) {
		t.Errorf("Expected CreatedAt to be %v, got %v", expiresAt, session.CreatedAt)
	}

}
func TestIsExpired(t *testing.T) {
	// Given
	user := NewFakeUser()
	sessionToken := "ABC123"
	createdAt := time.Now().Add(-25 * time.Hour) // Session created 25 hours ago
	loggedInAt := time.Now().Format(time.RFC3339)

	session := NewSession(&user, sessionToken, createdAt, loggedInAt)

	// When
	isExpired := session.IsExpired()

	// Then
	if !isExpired {
		t.Errorf("Expected session to be expired, but it was not")
	}

}
