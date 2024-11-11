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

type fileRepository struct {
	log   *slog.Logger
	db    *sql.DB
	files []*entities.File
}

func NewFileRepository(log *slog.Logger, database *sql.DB) interfaces.FileRepository {
	return &fileRepository{log, database, []*entities.File{}}
}

func (r *fileRepository) Add(file *entities.File) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "INSERT INTO files (url, name, filesize, mimetype,record_id) VALUES ($1, $2, $3, $4,$5)"
	_, err := r.db.ExecContext(ctx, query,
		file.Url,
		file.Name,
		file.FileSize,
		file.MimeType,
		file.RecordId,
	)

	if err != nil {
		return fmt.Errorf("failed to add files: %v", err)
	}

	return nil

}
func (r *fileRepository) GetByRecordId(id int) ([]*entities.File, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT * FROM files WHERE record_id = $1"
	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	var files []*entities.File
	for rows.Next() {
		file := &entities.File{}
		err := rows.Scan(
			&file.ID,
			&file.Url,
			&file.Name,
			&file.FileSize,
			&file.MimeType,
			&file.RecordId,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		files = append(files, file)
	}

	return files, nil

}

func (r *fileRepository) GetById(id int) (*entities.File, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT id, url, name, filesize, mimetype, record_id FROM files WHERE id = $1"
	row := r.db.QueryRowContext(ctx, query, id)

	file := &entities.File{}
	err := row.Scan(
		&file.ID,
		&file.Url,
		&file.Name,
		&file.FileSize,
		&file.MimeType,
		&file.RecordId,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("file with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get file: %v", err)
	}

	return file, nil
}
