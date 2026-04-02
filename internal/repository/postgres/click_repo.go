package repository

import (
	"context"
	"database/sql"

	"github.com/gewall/short-url/internal/domain"
	"github.com/gewall/short-url/internal/dto"
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

func (r *ClickRepo) Analytics(linkID string) (any, error) {
	query := `SELECT
	COUNT(*) AS total_clicks,
	COUNT(DISTINCT ip_hash) AS unique_clicks,
	COUNT(*) FILTER (WHERE clicked_at >= CURRENT_DATE) AS clicks_today,
	COUNT(*) FILTER (WHERE clicked_at >= NOW() - INTERVAL '7 days') AS clicks_7d,
	COUNT(*) FILTER (WHERE clicked_at >= NOW() - INTERVAL '30 days') AS clicks_30d
	FROM clicks WHERE link_id = $1;`
	row := r.db.QueryRow(query, linkID)
	var result struct {
		TotalClicks  int
		UniqueClicks int
		ClicksToday  int
		Clicks7d     int
		Clicks30d    int
	}
	err := row.Scan(&result.TotalClicks, &result.UniqueClicks, &result.ClicksToday, &result.Clicks7d, &result.Clicks30d)
	if err != nil {
		return nil, err
	}
	return map[string]int{
		"total_clicks":  result.TotalClicks,
		"unique_clicks": result.UniqueClicks,
		"clicks_today":  result.ClicksToday,
		"clicks_7d":     result.Clicks7d,
		"clicks_30d":    result.Clicks30d,
	}, nil
}

func (r *ClickRepo) AnalyticsByDate(linkID string) (any, error) {
	query := `SELECT DATE(hour)   AS date, SUM(click_count) AS clicks FROM click_stats_hourly WHERE link_id = $1 AND hour >= NOW() - INTERVAL '30 days' GROUP BY DATE(hour) ORDER BY date ASC;
`
	rows, err := r.db.Query(query, linkID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []dto.AnalyticsDate
	for rows.Next() {
		var r dto.AnalyticsDate
		if err := rows.Scan(&r.Date, &r.Clicks); err != nil {
			return nil, err
		}
		result = append(result, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *ClickRepo) AnalyticsByHour(linkId string) (any, error) {
	query := `SELECT hour, click_count  AS clicks FROM click_stats_hourly WHERE link_id = $1 AND hour >= DATE_TRUNC('hour', NOW() - INTERVAL '24 hours') ORDER BY hour ASC;`
	rows, err := r.db.Query(query, linkId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []dto.AnalyticsHour
	for rows.Next() {
		var r dto.AnalyticsHour
		if err := rows.Scan(&r.Hour, &r.Clicks); err != nil {
			return nil, err
		}

		result = append(result, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *ClickRepo) AnalyticsByCountry(linkId string) (any, error) {
	query := `SELECT
  country, SUM(click_count)AS clicks, ROUND(SUM(click_count) * 100.0 / SUM(SUM(click_count)) OVER (),2) AS percentage FROM click_stats_country WHERE link_id = $1 AND date >= NOW() - INTERVAL '30 days' GROUP BY country ORDER BY clicks DESC LIMIT 10;`
	rows, err := r.db.Query(query, linkId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []dto.AnalyticsCountry
	for rows.Next() {
		var r dto.AnalyticsCountry
		if err := rows.Scan(&r.Country, &r.Clicks, &r.Percentage); err != nil {
			return nil, err
		}

		result = append(result, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
