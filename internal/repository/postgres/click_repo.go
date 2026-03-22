package repository

import (
	"context"
	"database/sql"

	"github.com/gewall/short-url/internal/domain"
)

type ClickRepo struct {
	db *sql.DB
}

func NewClickRepo(db *sql.DB) *ClickRepo {
	return &ClickRepo{db: db}
}

func (r *ClickRepo) Create(ctx context.Context, click domain.Clicks) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. raw log
	_, err = tx.ExecContext(ctx, `
        INSERT INTO clicks (link_id, ip_hash, country, city, device, browser, os, referrer)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		click.LinkID, click.IpHash, click.Country, click.City,
		click.Device, click.Browser, click.Os, click.Referrer,
	)
	if err != nil {
		return err
	}

	// 2. hourly
	_, err = tx.ExecContext(ctx, `
        INSERT INTO click_stats_hourly (link_id, hour, click_count)
        VALUES ($1, DATE_TRUNC('hour', NOW()), 1)
        ON CONFLICT (link_id, hour)
        DO UPDATE SET click_count = click_stats_hourly.click_count + 1`,
		click.LinkID,
	)
	if err != nil {
		return err
	}

	// 3. country
	_, err = tx.ExecContext(ctx, `
        INSERT INTO click_stats_country (link_id, country, date, click_count)
        VALUES ($1, $2, CURRENT_DATE, 1)
        ON CONFLICT (link_id, country, date)
        DO UPDATE SET click_count = click_stats_country.click_count + 1`,
		click.LinkID, click.Country,
	)
	if err != nil {
		return err
	}

	// 4. device
	_, err = tx.ExecContext(ctx, `
        INSERT INTO click_stats_device (link_id, device, browser, os, date, click_count)
        VALUES ($1, $2, $3, $4, CURRENT_DATE, 1)
        ON CONFLICT (link_id, device, browser, os, date)
        DO UPDATE SET click_count = click_stats_device.click_count + 1`,
		click.LinkID, click.Device, click.Browser, click.Os,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}
