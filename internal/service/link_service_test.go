package service_test

import (
	"errors"
	"testing"
	"time"

	"github.com/gewall/short-url/internal/domain"
	"github.com/gewall/short-url/internal/dto"
	"github.com/gewall/short-url/internal/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockLinkRepository struct {
	mock.Mock
}

func (m *mockLinkRepository) Create(link domain.Link) (*domain.Link, error) {
	args := m.Called(link)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Link), nil
}

func (m *mockLinkRepository) FindByShortCode(code string) (*domain.Link, error) {
	args := m.Called(code)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Link), nil
}

func (m *mockLinkRepository) FindById(linkId uuid.UUID) (*domain.Link, error) {
	args := m.Called(linkId)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Link), nil
}

func (m *mockLinkRepository) FindAllByUser(userID uuid.UUID) ([]domain.Link, error) {
	args := m.Called(userID)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Link), nil
}

func (m *mockLinkRepository) Update(link domain.Link) (*domain.Link, error) {
	args := m.Called(link)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Link), nil
}

func (m *mockLinkRepository) Delete(linkId uuid.UUID) error {
	args := m.Called(linkId)
	if args.Error(0) != nil {
		return args.Error(0)
	}
	return nil
}

func TestCreateLink_Success(t *testing.T) {
	mockRepo := new(mockLinkRepository)
	linkService := service.NewLinkService(mockRepo)

	link := dto.LinkReqDTO{
		ShortCode:   "abc123",
		OriginalURL: "https://example.com",
		Title:       "ABC",
		ExpiresAt:   time.Now().Add(24 * time.Hour),
	}

	mockRepo.On("Create", mock.MatchedBy(func(l domain.Link) bool {
		return l.ShortCode == link.ShortCode && l.OriginalURL == link.OriginalURL && l.Title == link.Title && l.ExpiresAt.Equal(link.ExpiresAt)
	})).Return(&domain.Link{
		ShortCode:   link.ShortCode,
		OriginalURL: link.OriginalURL,
		Title:       link.Title,
		ExpiresAt:   time.Now().Add(24 * time.Hour),
	}, nil)

	result, err := linkService.CreateLink(link, uuid.NewString())
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, link.ShortCode, result.ShortCode)
	assert.Equal(t, link.OriginalURL, result.OriginalURL)
	assert.Equal(t, link.Title, result.Title)

	mockRepo.AssertExpectations(t)
}

func TestCreateLink_Failure(t *testing.T) {
	mockRepo := new(mockLinkRepository)
	linkService := service.NewLinkService(mockRepo)

	link := dto.LinkReqDTO{
		ShortCode:   "abc123",
		OriginalURL: "https://example.com",
		Title:       "ABC",
		ExpiresAt:   time.Now().Add(24 * time.Hour),
	}

	mockRepo.On("Create", mock.MatchedBy(func(l domain.Link) bool {
		return l.ShortCode == link.ShortCode && l.OriginalURL == link.OriginalURL && l.Title == link.Title && l.ExpiresAt.Equal(link.ExpiresAt)
	})).Return(nil, errors.New("create failed"))

	result, err := linkService.CreateLink(link, uuid.NewString())

	assert.Error(t, err)
	assert.Nil(t, result)

	mockRepo.AssertCalled(t, "Create", mock.MatchedBy(func(l domain.Link) bool {
		return l.ShortCode == link.ShortCode && l.OriginalURL == link.OriginalURL && l.Title == link.Title && l.ExpiresAt.Equal(link.ExpiresAt)
	}))
}

func TestFindShortCode_Success(t *testing.T) {
	mockRepo := new(mockLinkRepository)
	linkService := service.NewLinkService(mockRepo)

	link := dto.LinkReqDTO{
		ShortCode:   "abc123",
		OriginalURL: "https://example.com",
		Title:       "ABC",
		ExpiresAt:   time.Now().Add(24 * time.Hour),
	}

	mockRepo.On("FindByShortCode", link.ShortCode).Return(&domain.Link{
		ShortCode:   link.ShortCode,
		OriginalURL: link.OriginalURL,
		Title:       link.Title,
		ExpiresAt:   link.ExpiresAt,
	}, nil)

	result, err := linkService.FindByShortCode(link.ShortCode)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, link.ShortCode, result.ShortCode)
	assert.Equal(t, link.OriginalURL, result.OriginalURL)
	assert.Equal(t, link.Title, result.Title)
	assert.Equal(t, link.ExpiresAt.Format(time.RFC3339), result.ExpiresAt)

	mockRepo.AssertExpectations(t)
}

func TestFindByShortCode_Failure(t *testing.T) {
	mockRepo := new(mockLinkRepository)
	linkService := service.NewLinkService(mockRepo)

	link := dto.LinkReqDTO{
		ShortCode:   "abc123",
		OriginalURL: "https://example.com",
		Title:       "ABC",
		ExpiresAt:   time.Now().Add(24 * time.Hour),
	}

	mockRepo.On("FindByShortCode", link.ShortCode).Return(nil, errors.New("find failed"))

	result, err := linkService.FindByShortCode(link.ShortCode)

	assert.Error(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestFindById_Success(t *testing.T) {
	mockRepo := new(mockLinkRepository)
	linkService := service.NewLinkService(mockRepo)

	link := dto.LinkReqDTO{
		ShortCode:   "abc123",
		OriginalURL: "https://example.com",
		Title:       "ABC",
	}

	linkId := uuid.New()

	mockRepo.On("FindById", linkId).Return(&domain.Link{
		ShortCode:   link.ShortCode,
		OriginalURL: link.OriginalURL,
		Title:       link.Title,
		ExpiresAt:   link.ExpiresAt,
	}, nil)

	result, err := linkService.FindById(linkId)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, link.ShortCode, result.ShortCode)
	assert.Equal(t, link.OriginalURL, result.OriginalURL)
	assert.Equal(t, link.Title, result.Title)
	assert.Equal(t, link.ExpiresAt.Format(time.RFC3339), result.ExpiresAt)

	mockRepo.AssertExpectations(t)
}

func TestFindById_Failure(t *testing.T) {
	mockRepo := new(mockLinkRepository)
	linkService := service.NewLinkService(mockRepo)

	// link := dto.LinkReqDTO{
	// 	ShortCode:   "abc123",
	// 	OriginalURL: "https://example.com",
	// 	Title:       "ABC",
	// }
	//
	linkId := uuid.New()

	mockRepo.On("FindById", linkId).Return(nil, errors.New("find failed"))

	result, err := linkService.FindById(linkId)

	assert.Error(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestFindAllByUser_Success(t *testing.T) {
	mockRepo := new(mockLinkRepository)
	linkService := service.NewLinkService(mockRepo)

	userID := uuid.New()
	links := []domain.Link{
		{
			ShortCode:   "abc123",
			OriginalURL: "https://example.com",
			Title:       "ABC",
			ExpiresAt:   time.Now().Add(24 * time.Hour),
		},
		{
			ShortCode:   "def456",
			OriginalURL: "https://example.org",
			Title:       "DEF",
			ExpiresAt:   time.Now().Add(24 * time.Hour),
		},
	}

	mockRepo.On("FindAllByUser", userID).Return(links, nil)

	result, err := linkService.FindAllByUser(userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, links[0].ShortCode, result[0].ShortCode)
	assert.Equal(t, links[0].OriginalURL, result[0].OriginalURL)
	assert.Equal(t, links[0].Title, result[0].Title)
	assert.Equal(t, links[0].ExpiresAt.Format(time.RFC3339), result[0].ExpiresAt)
	assert.Equal(t, links[1].ShortCode, result[1].ShortCode)
	assert.Equal(t, links[1].OriginalURL, result[1].OriginalURL)
	assert.Equal(t, links[1].Title, result[1].Title)
	assert.Equal(t, links[1].ExpiresAt.Format(time.RFC3339), result[1].ExpiresAt)

	mockRepo.AssertExpectations(t)
}

func TestFindAllByUser_Failure(t *testing.T) {
	mockRepo := new(mockLinkRepository)
	linkService := service.NewLinkService(mockRepo)

	userID := uuid.New()

	mockRepo.On("FindAllByUser", userID).Return(nil, errors.New("find failed"))

	result, err := linkService.FindAllByUser(userID)

	assert.Error(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestUpdateLink_Success(t *testing.T) {
	mockRepo := new(mockLinkRepository)
	linkService := service.NewLinkService(mockRepo)

	linkID := uuid.New()
	link := &dto.LinkReqUpdateDTO{

		Title: "ABC",
	}

	mockRepo.On("Update", mock.MatchedBy(func(l domain.Link) bool {
		return l.ID == linkID
	})).Return(&domain.Link{
		ShortCode:   "abc123",
		OriginalURL: "https://example.com",
		Title:       "ABC",
		ExpiresAt:   time.Now().Add(24 * time.Hour),
	}, nil)

	result, err := linkService.UpdateLink(linkID, link)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "abc123", result.ShortCode)
	assert.Equal(t, "https://example.com", result.OriginalURL)
	assert.Equal(t, "ABC", result.Title)
	assert.Equal(t, time.Now().Add(24*time.Hour).Format(time.RFC3339), result.ExpiresAt)

	mockRepo.AssertExpectations(t)
}

func TestUpdateLink_Failure(t *testing.T) {
	mockRepo := new(mockLinkRepository)
	linkService := service.NewLinkService(mockRepo)

	linkID := uuid.New()
	link := &dto.LinkReqUpdateDTO{

		Title: "ABC",
	}

	mockRepo.On("Update", mock.MatchedBy(func(l domain.Link) bool {
		return l.ID == linkID
	})).Return(nil, errors.New("update failed"))

	result, err := linkService.UpdateLink(linkID, link)

	assert.Error(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestDeleteLink_Success(t *testing.T) {
	mockRepo := new(mockLinkRepository)
	linkService := service.NewLinkService(mockRepo)

	linkID := uuid.New()

	mockRepo.On("Delete", linkID).Return(nil)

	err := linkService.DeleteLink(linkID)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestDeleteLink_Failure(t *testing.T) {
	mockRepo := new(mockLinkRepository)
	linkService := service.NewLinkService(mockRepo)

	mockRepo.On("Delete", nil).Return(errors.New("delete failed"))

	err := linkService.DeleteLink(uuid.Nil)

	assert.Error(t, err)

	mockRepo.AssertNotCalled(t, "Delete", nil)
}
