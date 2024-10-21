package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vindosVP/snauth/internal/models"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Storage {
	return &Storage{db: db}
}

func (s *Storage) CreateUser(ctx context.Context, email string, hPassword []byte) (int64, error) {
	var id int64
	query := `INSERT INTO users (email, hashed_password, created_at, banned, deleted) 
				VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := s.db.QueryRow(ctx, query, email, hPassword, time.Now(), false, false).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Storage) UserByEmail(ctx context.Context, email string) (*models.User, error) {
	u := &models.User{}
	query := `SELECT id, email, hashed_password, created_at, banned, deleted FROM users WHERE email = $1`
	row := s.db.QueryRow(ctx, query, email)
	err := row.Scan(&u.Id, &u.Email, &u.HPassword, &u.CreatedAt, &u.Banned, &u.Deleted)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Storage) UserByID(ctx context.Context, id int64) (*models.User, error) {
	u := &models.User{}
	query := `SELECT id, email, hashed_password, created_at, banned, deleted FROM users WHERE id = $1`
	row := s.db.QueryRow(ctx, query, id)
	err := row.Scan(&u.Id, &u.Email, &u.HPassword, &u.CreatedAt, &u.Banned, &u.Deleted)
	if err != nil {
		return nil, err
	}
	return u, nil
}
