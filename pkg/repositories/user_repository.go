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

type usersRepository struct {
	log   *slog.Logger
	db    *sql.DB
	users []*entities.User
}

func NewUserRepository(log *slog.Logger, database *sql.DB) interfaces.UserRepository {
	return &usersRepository{log, database, []*entities.User{}}
}

func (r *usersRepository) SaveUser(user *entities.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := "INSERT INTO users (email, password, health_center_id, profile_id) VALUES ($1, $2, $3, $4)"
	_, err := r.db.ExecContext(ctx, query,
		user.Email,
		user.Password,
		user.HealthCenterId,
		user.ProfileId,
	)

	if err != nil {
		return fmt.Errorf("failed to add user: %v", err)
	}
	return nil
}
func (r *usersRepository) UpdatePassword(id int, newPassword string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "UPDATE users SET password = $1 WHERE profile_id = $2"
	_, err := r.db.ExecContext(ctx, query, newPassword, id)
	if err != nil {
		return fmt.Errorf("failed to update password: %v", err)
	}
	return nil
}

func (r *usersRepository) GetUser(email string, password string) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT id, email, password, health_center_id, profile_id FROM users WHERE email = $1 AND password = $2"
	row := r.db.QueryRowContext(ctx, query, email, password)

	user := &entities.User{}
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.HealthCenterId,
		&user.ProfileId,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return user, nil
}

func (r *usersRepository) GetUserByID(ID int) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT id, email, password, health_center_id, profile_id FROM users WHERE id = $1"
	row := r.db.QueryRowContext(ctx, query, ID)

	user := &entities.User{}
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.HealthCenterId,
		&user.ProfileId,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with ID %d not found", ID)
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return user, nil
}

func (r *usersRepository) GetUserByProfileID(ID int) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT id, email, password, health_center_id, profile_id FROM users WHERE profile_id = $1"
	row := r.db.QueryRowContext(ctx, query, ID)

	user := &entities.User{}
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.HealthCenterId,
		&user.ProfileId,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with ID %d not found", ID)
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return user, nil
}

func (r *usersRepository) GetAllDoctors() ([]*entities.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        SELECT u.id, u.email, u.password, u.health_center_id, u.profile_id
        FROM users u
        JOIN profiles p ON u.profile_id = p.id
        WHERE p.role = 'DOCTOR'
    `
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	doctors := []*entities.User{}
	for rows.Next() {
		doctor := &entities.User{}
		if err := rows.Scan(
			&doctor.ID,
			&doctor.Email,
			&doctor.Password,
			&doctor.HealthCenterId,
			&doctor.ProfileId,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		doctors = append(doctors, doctor)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating through rows: %v", err)
	}

	return doctors, nil
}

func (r *usersRepository) GetDoctorsForSpecialty(specialty string) ([]*entities.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        SELECT u.id, u.email, u.password, u.health_center_id, u.profile_id
        FROM users u
        JOIN profiles p ON u.profile_id = p.id
        WHERE p.specialty = $1
    `
	rows, err := r.db.QueryContext(ctx, query, specialty)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	doctors := []*entities.User{}
	for rows.Next() {
		doctor := &entities.User{}
		if err := rows.Scan(
			&doctor.ID,
			&doctor.Email,
			&doctor.Password,
			&doctor.HealthCenterId,
			&doctor.ProfileId,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		doctors = append(doctors, doctor)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating through rows: %v", err)
	}

	return doctors, nil
}
