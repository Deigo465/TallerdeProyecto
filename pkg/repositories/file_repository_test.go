package repository

import (
	"database/sql"
	"log/slog"
	"os"
	"testing"

	"github.com/open-wm/blockehr/pkg/domain/entities"
	"github.com/open-wm/blockehr/pkg/interfaces"
	"github.com/stretchr/testify/assert"
)

func initialiseDBFile(db *sql.DB) {
	setup := "../../db/setup.sql"
	executeFile(db, setup)

	_, err := db.Exec("INSERT INTO files (url, name, filesize, mimetype, record_id) VALUES ('http://example.com/file1','file','1000','text/plain',1)")
	if err != nil {
		panic(err)
	}
}

func setupTestFile() (*sql.DB, interfaces.FileRepository) {
	dbName := "../../db/test.sqlite3"
	os.Remove(dbName)
	db, err := sql.Open("sqlite", dbName)
	if err != nil {
		panic("Failed to open db database connection")
	}
	initialiseDBFile(db)

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := NewFileRepository(logger, db)
	return db, repo
}
func TestFileAdd(t *testing.T) {
	// Given
	db, repo := setupTestFile()
	defer db.Close()

	newFile := &entities.File{ID: 1, Url: "http://example.com/file1", Name: "file", FileSize: "1000", MimeType: "text/plain", RecordId: 1}

	// WHEN
	err := repo.Add(newFile)

	// THEN
	assert.NoError(t, err)

}
func TestGetByRecordId(t *testing.T) {
	// Given
	db, repo := setupTestFile()
	defer db.Close()
	expectedRecordId := &entities.File{ID: 1, Url: "http://example.com/file1", Name: "file1", FileSize: "1000", MimeType: "text/plain", RecordId: 1}

	expected := []*entities.File{
		{ID: 1, Url: "http://example.com/file1", Name: "file", FileSize: "1000", MimeType: "text/plain", RecordId: 1},
	}
	// WHEN
	result, err := repo.GetByRecordId(expectedRecordId.RecordId)
	if err != nil {
		t.Errorf("Error occurred: %v", err)
		return
	}

	// THEN
	assert.NoError(t, err)

	assert.Equal(t, expected, result)

	for _, file := range result {
		t.Logf("ID: %d, Url: %s, Name: %s, FileSize: %s, MimeType: %s, RecordId: %d",
			file.ID, file.Url, file.Name, file.FileSize, file.MimeType, file.RecordId)
	}
}
func TestGetById(t *testing.T) {
	// Given
	db, repo := setupTestFile()
	defer db.Close()

	expectedID := 1
	expected := &entities.File{ID: expectedID, Url: "http://example.com/file1", Name: "file1", FileSize: "1000", MimeType: "text/plain", RecordId: 1}

	// WHEN
	result, err := repo.GetById(expectedID)

	// THEN
	assert.NoError(t, err)
	assert.NotEqual(t, expected, result)

}
