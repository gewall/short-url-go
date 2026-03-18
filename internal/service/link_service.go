package service

import (
	"errors"
	"time"

	"github.com/gewall/short-url/internal/domain"
	"github.com/gewall/short-url/internal/dto"
	"github.com/gewall/short-url/pkg"
	"github.com/google/uuid"
)

type LinkRepository interface {
	Create(domain.Link) (*domain.Link, error)
	FindByShortCode(shortCode string) (*domain.Link, error)
	FindById(linkId uuid.UUID) (*domain.Link, error)
	FindAllByUser(userID uuid.UUID) ([]domain.Link, error)
	Update(link domain.Link) (*domain.Link, error)
	Delete(linkId uuid.UUID) error
}

type LinkService struct {
	repo LinkRepository
}

func NewLinkService(repo LinkRepository) *LinkService {
	return &LinkService{repo: repo}
}

func (s *LinkService) CreateLink(link dto.LinkReqDTO, userId string) (*dto.LinkRespDTO, error) {
	var code string
	var _link domain.Link
	_userId, err := uuid.Parse(userId)
	if err != nil {
		return &dto.LinkRespDTO{}, err
	}
	if link.ShortCode == "" {
		code = pkg.GenerateShortCode()
	} else {
		code = link.ShortCode
	}

	_link = domain.Link{
		UserID:      _userId,
		OriginalURL: link.OriginalURL,
		ShortCode:   code,
		Title:       link.Title,
		ExpiresAt:   link.ExpiresAt,
	}

	l, err := s.repo.Create(_link)
	if err != nil {
		return &dto.LinkRespDTO{}, err
	}

	return &dto.LinkRespDTO{
		ID:          l.ID,
		ShortCode:   l.ShortCode,
		OriginalURL: l.OriginalURL,
		Title:       l.Title,
		ExpiresAt:   l.ExpiresAt.Format(time.RFC3339),
	}, nil
}

func (s *LinkService) FindByShortCode(shortCode string) (*dto.LinkRespDTO, error) {
	l, err := s.repo.FindByShortCode(shortCode)
	if err != nil {
		return &dto.LinkRespDTO{}, err
	}

	return &dto.LinkRespDTO{
		ID:          l.ID,
		ShortCode:   l.ShortCode,
		OriginalURL: l.OriginalURL,
		Title:       l.Title,
		ExpiresAt:   l.ExpiresAt.Format(time.RFC3339),
	}, nil
}

func (s *LinkService) FindById(linkId uuid.UUID) (*dto.LinkRespDTO, error) {
	l, err := s.repo.FindById(linkId)

	if errors.Is(err, pkg.ErrRowsEmpty) {
		return &dto.LinkRespDTO{}, pkg.ErrRowsEmpty
	}
	if err != nil {
		return &dto.LinkRespDTO{}, err
	}

	return &dto.LinkRespDTO{
		ID:          l.ID,
		ShortCode:   l.ShortCode,
		OriginalURL: l.OriginalURL,
		Title:       l.Title,
		ExpiresAt:   l.ExpiresAt.Format(time.RFC3339),
	}, nil
}

func (s *LinkService) FindAllByUser(userID uuid.UUID) ([]dto.LinkRespDTO, error) {
	_links, err := s.repo.FindAllByUser(userID)
	switch {
	case errors.Is(err, pkg.ErrRowsEmpty):
		return nil, pkg.ErrRowsEmpty
	case err != nil:
		return nil, err
	}

	var links []dto.LinkRespDTO
	for _, l := range _links {
		links = append(links, dto.LinkRespDTO{
			ID:          l.ID,
			ShortCode:   l.ShortCode,
			OriginalURL: l.OriginalURL,
			Title:       l.Title,
			ExpiresAt:   l.ExpiresAt.Format(time.RFC3339),
		})
	}

	return links, nil
}

func (s *LinkService) UpdateLink(linkID uuid.UUID, link *dto.LinkReqUpdateDTO) (*dto.LinkRespDTO, error) {
	var _link domain.Link
	_link.ID = linkID
	_link.Title = link.Title
	_link.IsActive = link.IsActive == "true"

	l, err := s.repo.Update(_link)
	switch {
	case errors.Is(err, pkg.ErrRowsEmpty):
		return &dto.LinkRespDTO{}, pkg.ErrRowsEmpty
	case err != nil:
		return &dto.LinkRespDTO{}, err
	}

	return &dto.LinkRespDTO{
		ID:          l.ID,
		ShortCode:   l.ShortCode,
		OriginalURL: l.OriginalURL,
		Title:       l.Title,
		ExpiresAt:   l.ExpiresAt.Format(time.RFC3339),
		IsActive:    l.IsActive,
	}, nil
}

func (s *LinkService) DeleteLink(linkID uuid.UUID) error {
	if err := s.repo.Delete(linkID); err != nil {
		return err
	}
	return nil
}
