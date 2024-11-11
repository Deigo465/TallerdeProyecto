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

type profileRepository struct {
	log      *slog.Logger
	db       *sql.DB
	profiles []*entities.Profile
}

func NewProfileRepository(log *slog.Logger, database *sql.DB) interfaces.ProfileRepository {

	return &profileRepository{log, database, []*entities.Profile{}}
}

func (r *profileRepository) Add(profile *entities.Profile) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        INSERT INTO profiles (
            first_name, 
            mother_last_name, 
            father_last_name, 
            document_number, 
            gender, 
            phone, 
            contact_email, 
            date_of_birth, 
            cmp, 
            specialty, 
            role
        ) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id
    `

	err := r.db.QueryRowContext(ctx, query,
		profile.FirstName,
		profile.MotherLastName,
		profile.FatherLastName,
		profile.DocumentNumber,
		profile.Gender,
		profile.Phone,
		profile.ContactEmail,
		profile.DateOfBirth,
		profile.Cmp,
		profile.Specialty,
		profile.Role,
	).Scan(&profile.ID)

	if err != nil {
		return fmt.Errorf("failed to add profile: %w", err)
	}

	return nil
}
func (r *profileRepository) Update(profile *entities.Profile) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        UPDATE profiles
        SET first_name = ?,
            mother_last_name = ?,
            father_last_name = ?,
            document_number = ?,
            phone = ?,
            contact_email = ?,
			gender=?,
            date_of_birth = ?,
            cmp = ?,
			specialty = ?
        WHERE id = ?
    `
	_, err := r.db.ExecContext(ctx, query,
		profile.FirstName,
		profile.MotherLastName,
		profile.FatherLastName,
		profile.DocumentNumber,
		profile.Phone,
		profile.ContactEmail,
		profile.Gender,
		profile.DateOfBirth,
		profile.Cmp,
		profile.Specialty,
		profile.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update profile: %v", err)
	}
	return nil
}

func (r *profileRepository) GetByDocumentNumber(documentNumber string) (*entities.Profile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT * FROM profiles WHERE document_number = $1"

	row := r.db.QueryRowContext(ctx, query, documentNumber)

	profile := &entities.Profile{}
	err := row.Scan(
		&profile.ID,
		&profile.FirstName,
		&profile.MotherLastName,
		&profile.FatherLastName,
		&profile.DocumentNumber,
		&profile.Gender,
		&profile.Phone,
		&profile.ContactEmail,
		&profile.DateOfBirth,
		&profile.Cmp,
		&profile.Specialty,
		&profile.Role,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("profile with document number %s not found", documentNumber)
		}
		return nil, fmt.Errorf("failed to get profile: %v", err)
	}

	return profile, nil
}

func (r *profileRepository) GetAll() ([]*entities.Profile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT * FROM profiles"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	profiles := []*entities.Profile{}
	for rows.Next() {
		profile := &entities.Profile{}
		if err := rows.Scan(
			&profile.ID,
			&profile.FirstName,
			&profile.MotherLastName,
			&profile.FatherLastName,
			&profile.DocumentNumber,
			&profile.Gender,
			&profile.Phone,
			&profile.ContactEmail,
			&profile.DateOfBirth,
			&profile.Cmp,
			&profile.Specialty,
			&profile.Role,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		profiles = append(profiles, profile)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating through rows: %v", err)
	}

	return profiles, nil
}
func (r *profileRepository) GetById(id int) (*entities.Profile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT * FROM profiles WHERE id = ?"
	row := r.db.QueryRowContext(ctx, query, id)

	profile := &entities.Profile{}
	err := row.Scan(
		&profile.ID,
		&profile.FirstName,
		&profile.MotherLastName,
		&profile.FatherLastName,
		&profile.DocumentNumber,
		&profile.Gender,
		&profile.Phone,
		&profile.ContactEmail,
		&profile.DateOfBirth,
		&profile.Cmp,
		&profile.Specialty,
		&profile.Role,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("profile with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get appointment: %v", err)
	}

	return profile, nil
}
