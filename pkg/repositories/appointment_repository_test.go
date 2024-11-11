package repository

import (
	"database/sql"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/open-wm/blockehr/pkg/domain/entities"
	"github.com/open-wm/blockehr/pkg/interfaces"
	"github.com/stretchr/testify/assert"
	_ "modernc.org/sqlite"
)

func executeFile(db *sql.DB, fileName string) {
	b, err := os.ReadFile(fileName)
	if err != nil {
		panic("Failed to read file: " + fileName)
	}
	_, err = db.Exec(string(b))
	if err != nil {
		panic("Failed to execute file:" + fileName)
	}

}
func initialiseDB(db *sql.DB) {
	// read seed.sql file and execute it
	setup := "../../db/setup.sql"
	executeFile(db, setup)

	_, err := db.Exec("INSERT INTO appointments (specialty, status, date_appointment, doctor_id, patient_id, description) VALUES ('Cardiology', 0, '2021-06-01 10:00:00', 1, 1, 'Regular checkup')")
	if err != nil {
		panic(err)
	}
}

func setupTest() (*sql.DB, interfaces.AppointmentRepository) {
	dbName := "../../db/test.sqlite3"
	os.Remove(dbName)
	db, err := sql.Open("sqlite", dbName)
	if err != nil {
		panic("Failed to open db database connection")
	}
	initialiseDB(db)

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := NewAppointmentRepository(logger, db)
	return db, repo
}

func TestGetAllAppointments(t *testing.T) {
	// GIVEN
	db, repo := setupTest()
	defer db.Close()

	parseString := "2021-06-01 10:00:00"
	tomorrow, _ := time.Parse("2006-01-02 15:04:05", parseString)
	expected := []*entities.Appointment{
		{ID: 1, Specialty: "Cardiology", Status: entities.PENDING, StartsAt: tomorrow, DoctorId: 1, PatientId: 1, Description: "Regular checkup"},
	}

	// WHEN
	result, err := repo.GetAll()

	// THEN
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestGetAppointmentById(t *testing.T) {
	// GIVEN
	db, repo := setupTest()
	defer db.Close()

	parseString := "2021-06-01 10:00:00"
	tomorrow, _ := time.Parse("2006-01-02 15:04:05", parseString)
	expected := &entities.Appointment{ID: 1, Specialty: "Cardiology", Status: entities.PENDING, StartsAt: tomorrow, DoctorId: 1, PatientId: 1, Description: "Regular checkup"}

	// WHEN
	result, err := repo.GetById(expected.ID)

	// THEN
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestAddAppointment(t *testing.T) {
	// Given
	db, repo := setupTest()
	defer db.Close()

	tomorrow := time.Now().AddDate(0, 0, 1)
	newAppointment := &entities.Appointment{ID: 1, Specialty: "Cardiology", Status: entities.PENDING, StartsAt: tomorrow, DoctorId: 1, PatientId: 1, Description: "Regular checkup"}

	// WHEN
	err := repo.Add(newAppointment)

	// THEN
	assert.NoError(t, err)
}

func TestUpdateAppointment(t *testing.T) {
	// GIven
	db, repo := setupTest()
	defer db.Close()

	appointment := &entities.Appointment{
		ID:     1,
		Status: entities.PAID,
	}

	// WHEN
	err := repo.Update(appointment)

	// THEN
	assert.NoError(t, err)

	// check that it was updated
	updated, _ := repo.GetById(appointment.ID)
	assert.Equal(t, entities.PAID, updated.Status)
}
