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

type healthCenterRepository struct {
	log           *slog.Logger
	db            *sql.DB
	healthCenters []*entities.HealthCenter
}

func NewHealthCenterRepository(log *slog.Logger, database *sql.DB) interfaces.HealthCenterRepository {

	return &healthCenterRepository{log, database, []*entities.HealthCenter{}}
}

func (r *healthCenterRepository) Add(healthCenter *entities.HealthCenter) error {
	/*
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		query := "INSERT INTO health_centers (name, district, address) VALUES ($1, $2, $3)"
		_, err := r.db.ExecContext(ctx, query,
			healthCenter.Name,
			healthCenter.District,
			healthCenter.Address,
		)

		if err != nil {
			return fmt.Errorf("failed to add HealthCenter: %v", err)
		}
	*/
	return nil
}

func (r *healthCenterRepository) Update(healthCenter *entities.HealthCenter) error {

	return nil
}

func (r *healthCenterRepository) GetAll() ([]*entities.HealthCenter, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT id, name, district, address FROM health_centers"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	healthCenters := []*entities.HealthCenter{}
	for rows.Next() {
		healthCenter := &entities.HealthCenter{}
		if err := rows.Scan(
			&healthCenter.ID,
			&healthCenter.Name,
			&healthCenter.District,
			&healthCenter.Address,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		healthCenters = append(healthCenters, healthCenter)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating through rows: %v", err)
	}

	return healthCenters, nil
}
func (r *healthCenterRepository) GetByID(id int) (*entities.HealthCenter, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT id, name, district, address FROM health_centers WHERE id = $1"
	row := r.db.QueryRowContext(ctx, query, id)

	healthCenter := &entities.HealthCenter{}
	err := row.Scan(
		&healthCenter.ID,
		&healthCenter.Name,
		&healthCenter.District,
		&healthCenter.Address,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("health center with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get health center: %v", err)
	}

	return healthCenter, nil
}
