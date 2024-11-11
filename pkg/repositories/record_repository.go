package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/open-wm/blockehr/pkg/domain/entities"
	"github.com/open-wm/blockehr/pkg/interfaces"
)

type recordRepository struct {
	log     *slog.Logger
	db      *sql.DB
	records []*entities.Record
}

func NewRecordRepository(log *slog.Logger, database *sql.DB) interfaces.RecordRepository {
	return &recordRepository{log, database, []*entities.Record{}}
}

func (r *recordRepository) Add(record *entities.Record) (*entities.Record, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        INSERT INTO records (body, created_at, updated_at, patient_record_id, doctor_record_id, specialty) 
        VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
    `

	err := r.db.QueryRowContext(ctx, query,
		record.Body,
		record.CreatedAt,
		record.UpdatedAt,
		record.PatientId,
		record.DoctorId,
		record.Specialty,
	).Scan(&record.ID)

	if err != nil {
		return nil, fmt.Errorf("failed to add record: %v", err)
	}
	return record, nil
}

func (r *recordRepository) GetAll() ([]*entities.Record, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT id, body, created_at, updated_at, patient_record_id, doctor_record_id, specialty FROM records"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf("Error closing rows: %v", closeErr)
		}
	}()

	records := []*entities.Record{}
	for rows.Next() {
		record := &entities.Record{}
		if scanErr := rows.Scan(
			&record.ID,
			&record.Body,
			&record.CreatedAt,
			&record.UpdatedAt,
			&record.PatientId,
			&record.DoctorId,
			&record.Specialty,
		); scanErr != nil {
			return nil, fmt.Errorf("failed to scan row: %v", scanErr)
		}

		// get all files as well
		query := "SELECT id, url, filesize, name, mimetype, record_id FROM files WHERE record_id = $1"
		rows, err := r.db.QueryContext(ctx, query, record.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to execute query: %v", err)
		}
		defer func() {
			if closeErr := rows.Close(); closeErr != nil {
				log.Printf("Error closing rows: %v", closeErr)
			}
		}()

		files := []*entities.File{}
		for rows.Next() {
			file := &entities.File{}
			if scanErr := rows.Scan(
				&file.ID,
				&file.Url,
				&file.FileSize,
				&file.Name,
				&file.MimeType,
				&file.RecordId,
			); scanErr != nil {
				return nil, fmt.Errorf("failed to scan row: %v", scanErr)
			}
			files = append(files, file)
		}
		record.Files = files
		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating through rows: %v", err)
	}

	return records, nil
}
func (r *recordRepository) GetById(id int) (*entities.Record, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT id, body, created_at, updated_at, patient_record_id, doctor_record_id, specialty FROM records WHERE id = $1"
	row := r.db.QueryRowContext(ctx, query, id)

	record := &entities.Record{}
	err := row.Scan(
		&record.ID,
		&record.Body,
		&record.CreatedAt,
		&record.UpdatedAt,
		&record.PatientId,
		&record.DoctorId,
		&record.Specialty,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("record with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get record: %v", err)
	}

	return record, nil
}
func (r *recordRepository) UpdateByPatientID(patientID int, newBody string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "UPDATE records SET body = $1 WHERE patient_record_id = $2"
	_, err := r.db.ExecContext(ctx, query, newBody, patientID)
	if err != nil {
		return fmt.Errorf("failed to execute update query: %v", err)
	}

	return nil
}
