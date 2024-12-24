package store

import (
	"context"
	"database/sql"
	"errors"
)

type Role struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Level       int    `json:"level"`
	Description string `json:"description"`
}

type RoleStore struct {
	db *sql.DB
}

func (s *RoleStore) GetByName(ctx context.Context, slug string) (*Role, error) {
	query := `select id, name, level, description from roles where name = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var role Role
	err := s.db.QueryRowContext(ctx, query, slug).Scan(
		&role.ID,
		&role.Name,
		&role.Level,
		&role.Description,
	)

	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &role, nil
}
