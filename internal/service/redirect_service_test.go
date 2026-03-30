package service_test

import (
	"context"
	"errors"
	"net/netip"
	"testing"

	"github.com/gewall/short-url/internal/domain"
	"github.com/gewall/short-url/internal/dto"
	"github.com/gewall/short-url/internal/service"
	"github.com/gewall/short-url/pkg"
	"github.com/google/uuid"
	"github.com/oschwald/geoip2-golang/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockClickRepository struct {
	mock.Mock
}

func (m *mockClickRepository) Create(ctx context.Context, click domain.Clicks) error {
	args := m.Called(ctx, click)
	if args.Error(0) != nil {
		return args.Error(0)
	}
	return nil
}

type mockClickWorker struct {
	mock.Mock
}

func (m *mockClickWorker) Submit(click domain.Clicks) error {
	args := m.Called(click)
	if args.Error(0) != nil {
		return args.Error(0)
	}
	return nil
}

type mockGeoIP struct {
	mock.Mock
}

func (m *mockGeoIP) Country(addr netip.Addr) (*geoip2.Country, error) {
	args := m.Called(addr)

	if args.Get(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*geoip2.Country), nil
}

func TestRedirect_Success(t *testing.T) {
	geodb := new(mockGeoIP)
	linkRepo := new(mockLinkRepository)
	clickRepo := new(mockClickRepository)
	cw := new(mockClickWorker)
	svc := service.NewRedirectService(geodb, linkRepo, clickRepo, cw)

	redirect := &dto.Redirect{
		Code:    "testab",
		IP:      "8.8.8.8",
		Device:  "desktop",
		Browser: "Chrome",
		OS:      "Windows",
		Referer: "https://example.com",
	}
	id := uuid.New()

	linkRepo.On("FindByShortCode", redirect.Code).Return(&domain.Link{ID: id}, nil)

	ip, err := netip.ParseAddr(redirect.IP)
	assert.NoError(t, err)

	geodb.On("Country", ip).
		Return(&geoip2.Country{
			Country: geoip2.CountryRecord{
				ISOCode: "US",
			},
		}, nil)

	ipHash := pkg.GenerateHash(ip.String())
	cw.On("Submit", domain.Clicks{
		LinkID:   id,
		IpHash:   ipHash,
		Device:   redirect.Device,
		City:     "US",
		Country:  "US",
		Browser:  redirect.Browser,
		Os:       redirect.OS,
		Referrer: redirect.Referer,
	}).Return(nil)
	link, err := svc.Redirect(context.Background(), redirect)

	assert.NoError(t, err)
	assert.NotNil(t, link)
	assert.Equal(t, id, link.ID)
	linkRepo.AssertExpectations(t)

}

func TestRedirect_FailureCodeInactive(t *testing.T) {
	geodb := new(mockGeoIP)
	linkRepo := new(mockLinkRepository)
	clickRepo := new(mockClickRepository)
	cw := new(mockClickWorker)
	svc := service.NewRedirectService(geodb, linkRepo, clickRepo, cw)

	redirect := &dto.Redirect{
		Code:    "testab",
		IP:      "8.8.8.8",
		Device:  "desktop",
		Browser: "Chrome",
		OS:      "Windows",
		Referer: "https://example.com",
	}
	id := uuid.New()

	linkRepo.On("FindByShortCode", redirect.Code).Return(nil, pkg.ErrURLNotFound)

	ip, err := netip.ParseAddr(redirect.IP)
	assert.NoError(t, err)

	geodb.On("Country", ip).
		Return(&geoip2.Country{
			Country: geoip2.CountryRecord{
				ISOCode: "US",
			},
		}, nil)

	ipHash := pkg.GenerateHash(ip.String())
	cw.On("Submit", domain.Clicks{
		LinkID:   id,
		IpHash:   ipHash,
		Device:   redirect.Device,
		City:     "US",
		Country:  "US",
		Browser:  redirect.Browser,
		Os:       redirect.OS,
		Referrer: redirect.Referer,
	}).Return(nil)
	_, err = svc.Redirect(context.Background(), redirect)

	assert.Error(t, err)
	assert.Equal(t, pkg.ErrURLNotFound, err)
	linkRepo.AssertExpectations(t)
	geodb.AssertExpectations(t)
	cw.AssertNotCalled(t, "Submit", domain.Clicks{
		LinkID:   id,
		IpHash:   ipHash,
		Device:   redirect.Device,
		City:     "US",
		Country:  "US",
		Browser:  redirect.Browser,
		Os:       redirect.OS,
		Referrer: redirect.Referer,
	})
}

func TestRedirect_FailureGeoIP(t *testing.T) {
	geodb := new(mockGeoIP)
	linkRepo := new(mockLinkRepository)
	clickRepo := new(mockClickRepository)
	cw := new(mockClickWorker)
	svc := service.NewRedirectService(geodb, linkRepo, clickRepo, cw)

	redirect := &dto.Redirect{
		Code:    "testab",
		IP:      "8.8.8.8",
		Device:  "desktop",
		Browser: "Chrome",
		OS:      "Windows",
		Referer: "https://example.com",
	}
	id := uuid.New()

	linkRepo.On("FindByShortCode", redirect.Code).Return(&domain.Link{ID: id}, nil)

	ip, err := netip.ParseAddr(redirect.IP)
	assert.NoError(t, err)

	geodb.On("Country", ip).
		Return(nil,
			errors.New("Country not found"))

	ipHash := pkg.GenerateHash(ip.String())
	cw.On("Submit", domain.Clicks{
		LinkID:   id,
		IpHash:   ipHash,
		Device:   redirect.Device,
		City:     "US",
		Country:  "US",
		Browser:  redirect.Browser,
		Os:       redirect.OS,
		Referrer: redirect.Referer,
	}).Return(nil)
	_, err = svc.Redirect(context.Background(), redirect)

	assert.Error(t, err)

	geodb.AssertExpectations(t)
	linkRepo.AssertNotCalled(t, "FindByShortCode", redirect.Code)
	cw.AssertNotCalled(t, "Submit", domain.Clicks{
		LinkID:   id,
		IpHash:   ipHash,
		Device:   redirect.Device,
		City:     "US",
		Country:  "US",
		Browser:  redirect.Browser,
		Os:       redirect.OS,
		Referrer: redirect.Referer,
	})
}

func TestRedirect_FailureWorker(t *testing.T) {
	geodb := new(mockGeoIP)
	linkRepo := new(mockLinkRepository)
	clickRepo := new(mockClickRepository)
	cw := new(mockClickWorker)
	svc := service.NewRedirectService(geodb, linkRepo, clickRepo, cw)

	redirect := &dto.Redirect{
		Code:    "testab",
		IP:      "8.8.8.8",
		Device:  "desktop",
		Browser: "Chrome",
		OS:      "Windows",
		Referer: "https://example.com",
	}
	id := uuid.New()

	linkRepo.On("FindByShortCode", redirect.Code).Return(&domain.Link{ID: id}, nil)

	ip, err := netip.ParseAddr(redirect.IP)
	assert.NoError(t, err)

	geodb.On("Country", ip).
		Return(&geoip2.Country{
			Country: geoip2.CountryRecord{
				ISOCode: "US",
			},
		}, nil)

	ipHash := pkg.GenerateHash(ip.String())
	cw.On("Submit", domain.Clicks{
		LinkID:   id,
		IpHash:   ipHash,
		Device:   redirect.Device,
		City:     "US",
		Country:  "US",
		Browser:  redirect.Browser,
		Os:       redirect.OS,
		Referrer: redirect.Referer,
	}).Return(errors.New("Worker Fail"))
	link, err := svc.Redirect(context.Background(), redirect)

	assert.Error(t, err)
	assert.Nil(t, link)

	linkRepo.AssertExpectations(t)

}
