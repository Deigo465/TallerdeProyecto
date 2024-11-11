package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/open-wm/blockehr/pkg/domain/entities"
	"github.com/open-wm/blockehr/pkg/interfaces"
)

type appointmentRepository struct {
	log          *slog.Logger
	db           *sql.DB
	appointments []*entities.Appointment
}

func NewAppointmentRepository(log *slog.Logger, database *sql.DB) interfaces.AppointmentRepository {
	return &appointmentRepository{log, database, []*entities.Appointment{}}
}

func parseDate(dateStr string) (time.Time, error) {
	formats := []string{
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05 -07:00",
		"2006-01-02T15:04:05.000000-07:00",
		"2006-01-02T15:04:05.00000-07:00",
		"2006-01-02T15:04:05.0000-07:00",
		"2006-01-02T15:04:05.000-07:00",
		"2006-01-02T15:04:05-07:00",
		"2006-01-02 15:04:05 -0700 -0700",
	}

	var parseErr error
	for _, format := range formats {
		date, err := time.Parse(format, dateStr)
		if err == nil {
			return date, nil
		}
		parseErr = err
	}
	return time.Time{}, parseErr
}

func (r *appointmentRepository) Add(appointment *entities.Appointment) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := "INSERT INTO appointments (specialty, status, date_appointment, doctor_id, patient_id, description) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := r.db.ExecContext(ctx, query,
		appointment.Specialty,
		appointment.Status,
		appointment.StartsAt,
		appointment.DoctorId,
		appointment.PatientId,
		appointment.Description,
	)

	if err != nil {
		return fmt.Errorf("failed to add appointment: %v", err)
	}
	return nil

}
func (r *appointmentRepository) Update(appointment *entities.Appointment) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
			UPDATE appointments
			SET date_appointment = $1,
			status = $2
			WHERE id = $3
		`
	_, err := r.db.ExecContext(ctx, query, appointment.StartsAt, appointment.Status, appointment.ID)
	if err != nil {
		return fmt.Errorf("failed to update appointment status: %v", err)
	}
	return nil

}
func (r *appointmentRepository) GetAll() ([]*entities.Appointment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := "SELECT id, specialty, status, date_appointment, doctor_id, patient_id, description FROM appointments"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	appointments := []*entities.Appointment{}
	for rows.Next() {
		var appointment entities.Appointment
		var dateStr string
		err := rows.Scan(
			&appointment.ID,
			&appointment.Specialty,
			&appointment.Status,
			&dateStr,
			&appointment.DoctorId,
			&appointment.PatientId,
			&appointment.Description,
		)
		if err != nil {
			r.log.Error("Error scanning appointment:", "error", err)
			return nil, err
		}

		date, err := parseDate(dateStr)
		if err != nil {
			r.log.Error("Error", "error", err.Error(), "date", dateStr)
			return nil, err
		}
		appointment.StartsAt = date
		appointments = append(appointments, &appointment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error when iterating over query results %v", err)
	}

	return appointments, nil
}
func (r *appointmentRepository) GetById(id int) (*entities.Appointment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT id, specialty, status, date_appointment, doctor_id, patient_id, description FROM appointments WHERE id = $1"
	row := r.db.QueryRowContext(ctx, query, id)

	appointment := &entities.Appointment{}
	var dateStr string

	err := row.Scan(
		&appointment.ID,
		&appointment.Specialty,
		&appointment.Status,
		&dateStr,
		&appointment.DoctorId,
		&appointment.PatientId,
		&appointment.Description,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("appointment with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get appointment: %v", err)
	}
	date, err := parseDate(dateStr)
	if err != nil {
		r.log.Error("Error", "error", err.Error(), "date", dateStr)
		return nil, err
	}
	appointment.StartsAt = date
	return appointment, nil
}

func (r *appointmentRepository) GetByDoctorID(doctorID int) ([]*entities.Appointment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT * FROM appointments WHERE doctor_id = $1"
	rows, err := r.db.QueryContext(ctx, query, doctorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get appointments for doctor with ID %d: %v", doctorID, err)
	}
	defer rows.Close()

	var appointments []*entities.Appointment
	for rows.Next() {
		appointment := &entities.Appointment{}
		var dateStr string
		err := rows.Scan(
			&appointment.ID,
			&appointment.Specialty,
			&appointment.Status,
			&dateStr,
			&appointment.DoctorId,
			&appointment.PatientId,
			&appointment.Description)
		if err != nil {
			return nil, fmt.Errorf("failed to scan appointment: %v", err)
		}

		date, err := parseDate(dateStr)
		if err != nil {
			r.log.Error("Error", "error", err.Error(), "date", dateStr)
			return nil, err
		}
		appointment.StartsAt = date
		appointments = append(appointments, appointment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	return appointments, nil
}

func (r *appointmentRepository) GetByPatientID(patientId int) ([]*entities.Appointment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT * FROM appointments WHERE patient_id = $1"
	rows, err := r.db.QueryContext(ctx, query, patientId)
	if err != nil {
		return nil, fmt.Errorf("failed to get appointments for patient with ID %d: %v", patientId, err)
	}
	defer rows.Close()

	var appointments []*entities.Appointment
	for rows.Next() {
		appointment := &entities.Appointment{}

		var dateStr string
		err := rows.Scan(
			&appointment.ID,
			&appointment.Specialty,
			&appointment.Status,
			&dateStr,
			&appointment.DoctorId,
			&appointment.PatientId,
			&appointment.Description)
		if err != nil {
			return nil, fmt.Errorf("failed to scan appointment: %v", err)
		}

		date, err := parseDate(dateStr)
		if err != nil {
			r.log.Error("Error", "error", err.Error(), "date", dateStr)
			return nil, err
		}
		appointment.StartsAt = date
		appointments = append(appointments, appointment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	return appointments, nil
}
