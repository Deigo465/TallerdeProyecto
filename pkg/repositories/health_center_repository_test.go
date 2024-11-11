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

func initialiseDBHealthCenter(db *sql.DB) {
	// read seed.sql file and execute it
	setup := "../../db/setup.sql"
	executeFile(db, setup)

	_, err := db.Exec("INSERT INTO health_centers (name, district, address) VALUES	('Health Center 1', 'District 1', '123 Main St')")
	if err != nil {
		panic(err)
	}
}
func setupTestHealtCenter() (*sql.DB, interfaces.HealthCenterRepository) {
	dbName := "../../db/test.sqlite3"
	os.Remove(dbName)
	db, err := sql.Open("sqlite", dbName)
	if err != nil {
		panic("Failed to open db database connection")
	}
	initialiseDBHealthCenter(db)

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := NewHealthCenterRepository(logger, db)
	return db, repo
}

func TestHealthCenterGetAll(t *testing.T) {
	// GIVEN
	db, repo := setupTestHealtCenter()
	defer db.Close()

	expected := []*entities.HealthCenter{
		{
			ID:       1,
			Name:     "Health Center 1",
			District: "District 1",
			Address:  "123 Main St",
		},
	}

	// WHEN
	result, err := repo.GetAll()

	// THEN
	assert.NoError(t, err)
	assert.Equal(t, expected, result)

}

func TestHealthCenterGetById(t *testing.T) {
	// GIVEN
	db, repo := setupTestHealtCenter()
	defer db.Close()

	expected := &entities.HealthCenter{

		ID:       1,
		Name:     "Health Center 1",
		District: "District 1",
		Address:  "123 Main St",
	}

	// WHEN
	result, err := repo.GetByID(expected.ID)

	// THEN
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

/*
func TestHealthCenterAdd(t *testing.T) {
	// GIVEN
	// WHEN
	// THEN
}

func TestHealthCenterUpdate(t *testing.T) {
	// GIVEN
	// WHEN
	// THEN
}
*/
