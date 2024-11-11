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

func initialiseDRecord(db *sql.DB) {
	setup := "../../db/setup.sql"
	executeFile(db, setup)
	_, err := db.Exec("INSERT INTO records (body, created_at, updated_at, patient_record_id, doctor_record_id, specialty) VALUES('Record 3 body text', '2023-01-03 12:00:00', '2023-01-03 12:00:00', 2, 3, 'Odontologia')")

	if err != nil {
		panic(err)
	}
}

func setupTestRecord() (*sql.DB, interfaces.RecordRepository) {
	dbName := "../../db/test.sqlite3"
	os.Remove(dbName)
	db, err := sql.Open("sqlite", dbName)
	if err != nil {
		panic("Failed to open db database connection")
	}
	initialiseDRecord(db)

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := NewRecordRepository(logger, db)
	return db, repo
}

func TestRecordAdd(t *testing.T) {
	// Given
	db, repo := setupTestRecord()
	defer db.Close()

	newRecord := &entities.Record{

		Body:      "Record 3 body text",
		CreatedAt: "2023-01-03 12:00:00",
		UpdatedAt: "2023-01-03 12:00:00",
		PatientId: 2,
		DoctorId:  3,
		Specialty: "Odontologia",
	}

	// WHEN
	_, err := repo.Add(newRecord)

	// THEN
	assert.NoError(t, err)
}

func TestRecordGetAll(t *testing.T) {
	// GIVEN
	db, repo := setupTestRecord()
	defer db.Close()

	expected := []*entities.Record{{
		ID:        1,
		Body:      "Record 3 body text",
		CreatedAt: "2023-01-03T12:00:00Z",
		UpdatedAt: "2023-01-03T12:00:00Z",
		PatientId: 2,
		DoctorId:  3,
		Specialty: "Odontologia",
		Files:     []*entities.File{},
	},
	}

	// WHEN
	result, err := repo.GetAll()

	// THEN
	assert.NoError(t, err)
	assert.Equal(t, expected, result)

}

func TestRecordGetById(t *testing.T) {
	// GIVEN
	db, repo := setupTestRecord()
	defer db.Close()

	expected := &entities.Record{
		ID:        1,
		Body:      "Record 3 body text",
		CreatedAt: "2023-01-03T12:00:00Z",
		UpdatedAt: "2023-01-03T12:00:00Z",
		PatientId: 2,
		DoctorId:  3,
		Specialty: "Odontologia",
	}

	// WHEN
	result, err := repo.GetById(expected.ID)

	// THEN
	assert.NoError(t, err)
	assert.Equal(t, expected, result)

}

func TestUpdateByPatientId(t *testing.T) {
	// GIVEN
	db, repo := setupTestRecord()
	defer db.Close()

	expected := &entities.Record{
		ID:        1,
		Body:      "body updated",
		CreatedAt: "2023-01-03T12:00:00Z",
		UpdatedAt: "2023-01-03T12:00:00Z",
		PatientId: 2,
		DoctorId:  3,
		Specialty: "Odontologia",
	}

	// WHEN
	err := repo.UpdateByPatientID(expected.PatientId, expected.Body)

	// THEN
	assert.NoError(t, err)

}
