package repository

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/open-wm/blockehr/pkg/domain/entities"
	"github.com/open-wm/blockehr/pkg/interfaces"
	"github.com/stretchr/testify/assert"
)

func initialiseDbUser(db *sql.DB) {
	setup := "../../db/setup.sql"
	executeFile(db, setup)
	_, err := db.Exec("INSERT INTO users(email, password, health_center_id, profile_id) VALUES('doctor@example.com', '123456', 1, 1)")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("INSERT INTO profiles(first_name, mother_last_name, father_last_name, document_number, gender, phone, date_of_birth, cmp, specialty, role) VALUES('doctor', 'Doe', 'Smith', '12345678', 'Male', '123456789', '1970-01-15', 'CMP001', 'Cardiology', 'DOCTOR')")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("INSERT INTO users(email, password, health_center_id, profile_id) VALUES('rod@example.com', '123456', 1, 2)")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("INSERT INTO profiles(first_name, mother_last_name, father_last_name, document_number, gender, phone, date_of_birth, cmp, specialty, role) VALUES('staff', 'Rod', 'Smith', '12345678', 'Male', '123456789', '1970-01-15', 'CMP001', '', 'STAFF')")

	if err != nil {
		panic(err)
	}
}

func setupTestUser() (*sql.DB, interfaces.UserRepository) {
	dbName := "../../db/test.sqlite3"
	os.Remove(dbName)
	db, err := sql.Open("sqlite", dbName)
	if err != nil {
		panic("Failed to open db database connection")
	}
	initialiseDbUser(db)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	repo := NewUserRepository(logger, db)
	return db, repo
}

func TestUserSaveUser(t *testing.T) {
	// Given
	db, repo := setupTestUser()
	defer db.Close()
	newUser := &entities.User{
		Email:          "staff@example.com",
		Password:       "123456",
		HealthCenterId: 1,
		ProfileId:      1,
	}

	// WHEN
	err := repo.SaveUser(newUser)

	// THEN
	assert.NoError(t, err)
}
func TestUserGetAllDoctors(t *testing.T) {
	// Given
	db, repo := setupTestUser()
	defer db.Close()
	expected := []*entities.User{{
		ID:             1,
		Email:          "doctor@example.com",
		Password:       "123456",
		HealthCenterId: 1,
		ProfileId:      1,
	}}

	// WHEN
	result, err := repo.GetAllDoctors()

	// THEN
	assert.NoError(t, err)

	assert.Equal(t, expected, result)

}

func TestUserGetDoctorsForSpeciality(t *testing.T) {
	// Given
	db, repo := setupTestUser()
	defer db.Close()
	expected := []*entities.User{{
		ID:             1,
		Email:          "doctor@example.com",
		Password:       "123456",
		HealthCenterId: 1,
		ProfileId:      1,
	}}

	// WHEN
	result, err := repo.GetDoctorsForSpecialty("Cardiology")

	// THEN
	assert.NoError(t, err)

	assert.Equal(t, expected, result)

}

func TestUserUpdatePassword(t *testing.T) {
	// Given
	db, repo := setupTestUser()
	defer db.Close()
	expected := &entities.User{
		ID:             1,
		Email:          "doctor@example.com",
		Password:       "password",
		HealthCenterId: 1,
		ProfileId:      1,
	}

	// WHEN
	err := repo.UpdatePassword(expected.ID, expected.Password)

	// THEN
	assert.NoError(t, err)
	result, _ := repo.GetAllDoctors()

	for _, doctor := range result {
		fmt.Printf("ID: %d, Email: %s, Password: %s, HealthCenterId: %d, ProfileId: %d\n",
			doctor.ID, doctor.Email, doctor.Password, doctor.HealthCenterId, doctor.ProfileId)
	}

}

func TestGetUserById(t *testing.T) {
	// GIVEN
	db, repo := setupTestUser()
	defer db.Close()

	expected := &entities.User{
		ID:             1,
		Email:          "doctor@example.com",
		Password:       "123456",
		HealthCenterId: 1,
		ProfileId:      1,
	}

	// WHEN
	result, err := repo.GetUserByID(expected.ID)

	// THEN
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}
func TestGetUserByEmailAndPassword(t *testing.T) {
	// GIVEN
	db, repo := setupTestUser()
	defer db.Close()

	expected := &entities.User{
		ID:             1,
		Email:          "doctor@example.com",
		Password:       "123456",
		HealthCenterId: 1,
		ProfileId:      1,
	}

	// WHEN
	result, err := repo.GetUser(expected.Email, expected.Password)

	// THEN
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}
