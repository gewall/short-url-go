package domain

import (
	"time"

	"github.com/google/uuid"
)

type Link struct {
	ID          uuid.UUID `db:"id"`
	UserID      uuid.UUID `db:"user_id"`
	OriginalURL string    `db:"original_url"`
	ShortCode   string    `db:"short_code"`
	Title       string    `db:"title"`
	IsActive    bool      `db:"is_active"`
	ExpiresAt   time.Time `db:"expires_at"`
	CreatedAt   time.Time `db:"created_at"`
}
