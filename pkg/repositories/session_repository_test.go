package repository

import (
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/open-wm/blockehr/pkg/domain/entities"
	"github.com/open-wm/blockehr/pkg/interfaces"
	"github.com/stretchr/testify/assert"
)

func initialiseDbSession(db *sql.DB) {
	setup := "../../db/setup.sql"
	executeFile(db, setup)
	_, err := db.Exec("INSERT INTO sessions (user_id, token, created_at, updated_at) VALUES (1, 'STAFF', '2025-01-01 10:00:00', '2023-01-01 10:00:00')")

	if err != nil {
		panic(err)
	}
}

func setupTestSession() (*sql.DB, interfaces.SessionRepository) {
	dbName := "../../db/test.sqlite3"
	os.Remove(dbName)
	db, err := sql.Open("sqlite", dbName)
	if err != nil {
		panic("Failed to open db database connection")
	}
	initialiseDbSession(db)

	repo := NewSessionRepository(db)
	return db, repo
}

func TestGetSessionByToken(t *testing.T) {
	// GIven
	db, repo := setupTestSession()
	defer db.Close()

	validSession := &entities.Session{
		ID:        1,
		UserID:    1,
		Token:     "STAFF",
		CreatedAt: time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
	}
	// WHEN
	result := repo.GetSession(validSession.Token)

	// THEN
	assert.Equal(t, validSession, result)

}
func TestSaveSession(t *testing.T) {
	// GIven
	db, repo := setupTestSession()
	defer db.Close()

	validSession := &entities.Session{
		ID:        1,
		UserID:    1,
		Token:     "STAFF",
		CreatedAt: time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2027, 1, 1, 10, 0, 0, 0, time.UTC),
	}
	// WHEN
	result, err := repo.SaveSession(validSession)

	// THEN
	assert.NoError(t, err)
	assert.Equal(t, validSession, result)

}
func TestDeleteSession(t *testing.T) {
	// Given
	db, repo := setupTestSession()
	defer db.Close()

	validSession := &entities.Session{
		ID:        1,
		UserID:    1,
		Token:     "STAFF",
		CreatedAt: time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2027, 1, 1, 10, 0, 0, 0, time.UTC),
	}
	_, err := repo.SaveSession(validSession)
	assert.NoError(t, err)

	// WHEN
	err = repo.DeleteSession(validSession.Token)

	// THEN
	assert.NoError(t, err)

	result := repo.GetSession(validSession.Token)
	assert.Nil(t, result, "Expected session not to be present after deletion")
}
