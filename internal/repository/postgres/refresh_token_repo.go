package repository

import (
	"database/sql"
	"time"

	"github.com/gewall/short-url/internal/domain"
)

type refreshTokenRepo struct {
	db *sql.DB
}

func NewRefreshTokenRepo(db *sql.DB) *refreshTokenRepo {
	return &refreshTokenRepo{db: db}
}

func (r *refreshTokenRepo) Create(token domain.RefreshToken) error {
	query := `INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(query, token.UserId, token.TokenHash, token.ExpiresAt)
	if err != nil {

		return err
	}
	return nil
}

func (r *refreshTokenRepo) Find(token string) (*domain.RefreshToken, error) {
	query := `SELECT * FROM refresh_tokens WHERE token_hash = $1 AND revoked_at IS NULL AND expires_at > $2`
	var rt domain.RefreshToken
	err := r.db.QueryRow(query, token, time.Now()).Scan(&rt.ID, &rt.UserId, &rt.TokenHash, &rt.ExpiresAt, &rt.CreatedAt, &rt.RevokedAt)
	if err != nil {
		return nil, err
	}
	return &rt, nil
}

func (r *refreshTokenRepo) UpdateRevoke(token string) error {
	query := `UPDATE refresh_tokens SET revoked_at = $1 WHERE token_hash = $2`
	_, err := r.db.Exec(query, time.Now(), token)
	if err != nil {
		return err
	}
	return nil
}

func (r *refreshTokenRepo) Delete(token string) error {
	query := `DELETE FROM refresh_tokens WHERE token_hash = $1`
	_, err := r.db.Exec(query, token)
	if err != nil {
		return err
	}
	return nil
}
