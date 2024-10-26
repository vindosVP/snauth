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

func (s *Storage) SetDeletedToUser(ctx context.Context, userId int64, isDeleted bool) (bool, error) {
	var deleted bool
	query := `UPDATE users SET is_deleted = $1 WHERE id = $2 RETURNING is_deleted`
	err := s.db.QueryRow(ctx, query, isDeleted, userId).Scan(&deleted)
	if err != nil {
		return false, err
	}
	return deleted, nil
}

func (s *Storage) SetBannedToUser(ctx context.Context, userId int64, isBanned bool) (bool, error) {
	var banned bool
	query := `UPDATE users SET is_banned = $1 WHERE id = $2 RETURNING is_banned`
	err := s.db.QueryRow(ctx, query, isBanned, userId).Scan(&banned)
	if err != nil {
		return false, err
	}
	return banned, nil
}

func (s *Storage) SetAdminToUser(ctx context.Context, userId int64, isAdmin bool) (bool, error) {
	var admin bool
	query := `UPDATE users SET is_admin = $1 WHERE id = $2 RETURNING is_admin`
	err := s.db.QueryRow(ctx, query, isAdmin, userId).Scan(&admin)
	if err != nil {
		return false, err
	}
	return admin, nil
}

func (s *Storage) CreateUser(ctx context.Context, email string, hPassword []byte) (int64, error) {
	var id int64
	var usersCount int64
	isAdmin := false

	query := "SELECT COUNT(id) FROM users"
	err := s.db.QueryRow(ctx, query).Scan(&usersCount)
	if err != nil {
		return 0, err
	}
	if usersCount == 0 {
		isAdmin = true
	}
	query = `INSERT INTO users (email, hashed_password, created_at, is_banned, is_deleted, is_admin) 
				VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err = s.db.QueryRow(ctx, query, email, hPassword, time.Now(), false, false, isAdmin).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Storage) UserByEmail(ctx context.Context, email string) (*models.User, error) {
	u := &models.User{}
	query := `SELECT id, email, hashed_password, created_at, is_banned, is_deleted, is_admin 
				FROM users WHERE email = $1`
	row := s.db.QueryRow(ctx, query, email)
	err := row.Scan(&u.Id, &u.Email, &u.HPassword, &u.CreatedAt, &u.IsBanned, &u.IsDeleted, &u.IsAdmin)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Storage) UserByID(ctx context.Context, id int64) (*models.User, error) {
	u := &models.User{}
	query := `SELECT id, email, hashed_password, created_at, is_banned, is_deleted, is_admin 
				FROM users WHERE id = $1`
	row := s.db.QueryRow(ctx, query, id)
	err := row.Scan(&u.Id, &u.Email, &u.HPassword, &u.CreatedAt, &u.IsBanned, &u.IsDeleted, &u.IsAdmin)
	if err != nil {
		return nil, err
	}
	return u, nil
}
