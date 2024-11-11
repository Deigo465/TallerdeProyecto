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

func initialiseDBProfile(db *sql.DB) {
	setup := "../../db/setup.sql"
	executeFile(db, setup)

	_, err := db.Exec("INSERT INTO profiles(first_name, mother_last_name, father_last_name, document_number, gender, phone, date_of_birth, cmp, specialty, role) VALUES('Staff', 'Doe', 'Smith', '12345678', 'Male', '123456789', '1970-01-15', 'CMP001', 'Cardiology', 'STAFF')")
	if err != nil {
		panic(err)
	}
}

func setupTestProfile() (*sql.DB, interfaces.ProfileRepository) {
	dbName := "../../db/test.sqlite3"
	os.Remove(dbName)
	db, err := sql.Open("sqlite", dbName)
	if err != nil {
		panic("Failed to open db database connection")
	}
	initialiseDBProfile(db)

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := NewProfileRepository(logger, db)
	return db, repo
}
func TestProfileAdd(t *testing.T) {
	// Given
	db, repo := setupTestProfile()
	defer db.Close()

	newProfile := &entities.Profile{
		FirstName:      "Staff",
		MotherLastName: "Doe",
		FatherLastName: "Smith",
		DocumentNumber: "12345678",
		Gender:         "Male",
		Phone:          "123456789",
		DateOfBirth:    "2006-01-02",
		Cmp:            "CMP001",
		Specialty:      "Cardiology",
		Role:           "STAFF",
	}

	// WHEN
	err := repo.Add(newProfile)

	// THEN
	assert.NoError(t, err)
}

func TestProfileUpdate(t *testing.T) {
	// GIven
	db, repo := setupTestProfile()
	defer db.Close()

	updateProfile := &entities.Profile{
		ID:             1,
		FirstName:      "Staff",
		MotherLastName: "Kunimoto",
		FatherLastName: "Smith",
		DocumentNumber: "12345678",
		Gender:         "Male",
		Phone:          "123456789",
		DateOfBirth:    "1970-01-15T00:00:00Z",
		Cmp:            "CMP001",
		Specialty:      "Cardiology",
		Role:           "STAFF",
	}
	// WHEN
	err := repo.Update(updateProfile)

	// THEN
	assert.NoError(t, err)

	updated, _ := repo.GetById(updateProfile.ID)
	assert.Equal(t, updateProfile, updated)

}
func TestProfileGetById(t *testing.T) {
	// GIven
	db, repo := setupTestProfile()
	defer db.Close()

	validProfile := &entities.Profile{
		ID:             1,
		FirstName:      "Staff",
		MotherLastName: "Doe",
		FatherLastName: "Smith",
		DocumentNumber: "12345678",
		Gender:         "Male",
		Phone:          "123456789",
		DateOfBirth:    "1970-01-15T00:00:00Z",
		Cmp:            "CMP001",
		Specialty:      "Cardiology",
		Role:           "STAFF",
	}
	// WHEN
	result, err := repo.GetById(validProfile.ID)

	// THEN
	assert.NoError(t, err)
	assert.Equal(t, validProfile, result)

}
func TestProfileGetAll(t *testing.T) {
	// GIven
	db, repo := setupTestProfile()
	defer db.Close()
	expected := []*entities.Profile{
		{ID: 1,
			FirstName:      "Staff",
			MotherLastName: "Doe",
			FatherLastName: "Smith",
			DocumentNumber: "12345678",
			Gender:         "Male",
			Phone:          "123456789",
			DateOfBirth:    "1970-01-15T00:00:00Z",
			Cmp:            "CMP001",
			Specialty:      "Cardiology",
			Role:           "STAFF"},
	}
	// WHEN
	result, err := repo.GetAll()

	// THEN
	assert.NoError(t, err)
	assert.Equal(t, expected, result)

}

func TestProfileGetByDocumentNumber(t *testing.T) {
	// GIven
	db, repo := setupTestProfile()
	defer db.Close()

	validProfile := &entities.Profile{
		ID:             1,
		FirstName:      "Staff",
		MotherLastName: "Doe",
		FatherLastName: "Smith",
		DocumentNumber: "12345678",
		Gender:         "Male",
		Phone:          "123456789",
		DateOfBirth:    "1970-01-15T00:00:00Z",
		Cmp:            "CMP001",
		Specialty:      "Cardiology",
		Role:           "STAFF",
	}
	// WHEN
	result, err := repo.GetByDocumentNumber(validProfile.DocumentNumber)

	// THEN
	assert.NoError(t, err)
	assert.Equal(t, validProfile, result)

}
