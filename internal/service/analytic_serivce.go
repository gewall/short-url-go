package service

import (
	"github.com/gewall/short-url/internal/dto"
	"github.com/google/uuid"
)

type AnalyticRepository interface {
	Analytics(string) (any, error)
	AnalyticsByHour(string) (any, error)
	AnalyticsByDate(string) (any, error)
	AnalyticsByCountry(string) (any, error)
}

type AnalyticService struct {
	repo AnalyticRepository
}

func NewAnalyticService(repo AnalyticRepository) *AnalyticService {
	return &AnalyticService{repo: repo}
}

func (s *AnalyticService) Analytics(linkID uuid.UUID) (*dto.AnalyticsResp, error) {
	res, err := s.repo.Analytics(linkID.String())
	if err != nil {
		return nil, err
	}

	analytics := dto.AnalyticsResp{
		TotalClicks:  res.(map[string]int)["total_clicks"],
		UniqueClicks: res.(map[string]int)["unique_clicks"],
		ClicksToday:  res.(map[string]int)["clicks_today"],
		Clicks7d:     res.(map[string]int)["clicks_7d"],
		Clicks30d:    res.(map[string]int)["clicks_30d"],
	}

	return &analytics, nil
}

func (s *AnalyticService) TimeSeries(linkID uuid.UUID) (*dto.AnalyticsTimeSeries, error) {
	hour, err := s.repo.AnalyticsByHour(linkID.String())
	if err != nil {
		return nil, err
	}

	date, err := s.repo.AnalyticsByDate(linkID.String())
	if err != nil {
		return nil, err
	}

	series := dto.AnalyticsTimeSeries{
		Hours: hour.([]dto.AnalyticsHour),
		Dates: date.([]dto.AnalyticsDate),
	}

	return &series, nil
}

func (s *AnalyticService) Country(linkId uuid.UUID) ([]dto.AnalyticsCountry, error) {
	res, err := s.repo.AnalyticsByCountry(linkId.String())
	if err != nil {
		return nil, err
	}

	country := res.([]dto.AnalyticsCountry)

	return country, nil
}
