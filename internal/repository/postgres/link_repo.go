package repository

import (
	"database/sql"
	"errors"

	"github.com/gewall/short-url/internal/domain"
	"github.com/gewall/short-url/pkg"
	"github.com/google/uuid"
)

type linkRepo struct {
	db *sql.DB
}

func NewLinkRepo(db *sql.DB) *linkRepo {
	return &linkRepo{db: db}
}

func (r *linkRepo) Create(link domain.Link) (*domain.Link, error) {
	var _link domain.Link
	query := `INSERT INTO links (user_id, original_url, short_code, title, expires_at) VALUES ($1, $2, $3, $4, $5) RETURNING *`

	err := r.db.QueryRow(query, link.UserID, link.OriginalURL, link.ShortCode, link.Title, link.ExpiresAt).Scan(&_link.ID, &_link.UserID, &_link.OriginalURL, &_link.ShortCode, &_link.Title, &_link.IsActive, &_link.ExpiresAt, &_link.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &_link, nil
}

func (r *linkRepo) FindByShortCode(shortCode string) (*domain.Link, error) {
	var link domain.Link
	query := `SELECT * FROM links WHERE short_code = $1`

	err := r.db.QueryRow(query, shortCode).Scan(&link.ID, &link.UserID, &link.OriginalURL, &link.ShortCode, &link.Title, &link.IsActive, &link.ExpiresAt, &link.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &link, nil
}

func (r *linkRepo) FindById(linkId uuid.UUID) (*domain.Link, error) {
	var link domain.Link
	query := `SELECT * FROM links WHERE id = $1`

	err := r.db.QueryRow(query, linkId).Scan(&link.ID, &link.UserID, &link.OriginalURL, &link.ShortCode, &link.Title, &link.IsActive, &link.ExpiresAt, &link.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, pkg.ErrRowsEmpty
	}
	if err != nil {
		return nil, err
	}

	return &link, nil
}

func (r *linkRepo) FindAllByUser(userID uuid.UUID) ([]domain.Link, error) {
	var links []domain.Link
	query := `SELECT * FROM links WHERE user_id = $1`

	rows, err := r.db.Query(query, userID)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, pkg.ErrRowsEmpty
	case err != nil:
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var link domain.Link
		err := rows.Scan(&link.ID, &link.UserID, &link.OriginalURL, &link.ShortCode, &link.Title, &link.IsActive, &link.ExpiresAt, &link.CreatedAt)
		if err != nil {
			return nil, err
		}
		links = append(links, link)
	}

	return links, nil
}

func (r *linkRepo) Update(link domain.Link) (*domain.Link, error) {
	var _link domain.Link
	query := `UPDATE links SET title = $1, is_active = $2 WHERE id = $3 RETURNING *`

	err := r.db.QueryRow(query, link.Title, link.IsActive, link.ID).Scan(&_link.ID, &_link.UserID, &_link.OriginalURL, &_link.ShortCode, &_link.Title, &_link.IsActive, &_link.ExpiresAt, &_link.CreatedAt)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, pkg.ErrRowsEmpty
	case err != nil:
		return nil, err
	}

	return &_link, nil
}

func (r *linkRepo) Delete(linkId uuid.UUID) error {
	query := `DELETE FROM links WHERE id = $1`

	_, err := r.db.Exec(query, linkId)
	return err
}
