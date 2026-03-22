package domain

import (
	"time"

	"github.com/google/uuid"
)

type Clicks struct {
	ID        int       `db:"id"`
	LinkID    uuid.UUID `db:"link_id"`
	ClickedAt time.Time `db:"clicked_at"`
	IpHash    string    `db:"ip_hash"`
	Country   string    `db:"country"`
	City      string    `db:"city"`
	Device    string    `db:"device"`
	Browser   string    `db:"browser"`
	Os        string    `db:"os"`
	Referrer  string    `db:"referrer"`
}

type ClickStatsHourly struct {
	ID         int       `db:"id"`
	LinkID     uuid.UUID `db:"link_id"`
	Hour       time.Time `db:"hour"`
	ClickCount int       `db:"click_count"`
}

type ClickStatsCountry struct {
	ID         int       `db:"id"`
	LinkID     uuid.UUID `db:"link_id"`
	Country    string    `db:"country"`
	ClickCount int       `db:"click_count"`
	Date       time.Time `db:"date"`
}

type ClickStatsDevice struct {
	ID         int       `db:"id"`
	LinkID     uuid.UUID `db:"link_id"`
	Device     string    `db:"device"`
	Browser    string    `db:"browser"`
	Os         string    `db:"os"`
	ClickCount int       `db:"click_count"`
	Date       time.Time `db:"date"`
}
