package dto

import "time"

type AnalyticsResp struct {
	TotalClicks  int `json:"total_clicks"`
	UniqueClicks int `json:"unique_clicks"`
	ClicksToday  int `json:"clicks_today"`
	Clicks7d     int `json:"clicks_7d"`
	Clicks30d    int `json:"clicks_30d"`
}

type AnalyticsTimeSeries struct {
	Hours []AnalyticsHour `json:"hours"`
	Dates []AnalyticsDate `json:"dates"`
}

type AnalyticsHour struct {
	Hour   time.Time `json:"hour"`
	Clicks int       `json:"clicks"`
}

type AnalyticsDate struct {
	Date   time.Time `json:"date"`
	Clicks int       `json:"clicks"`
}

type AnalyticsCountry struct {
	Country    string  `json:"country"`
	Clicks     int     `json:"clicks"`
	Percentage float64 `json:"percentage"`
}
