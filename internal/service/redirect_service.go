package service

import (
	"context"
	"errors"
	"log"
	"net/netip"

	"github.com/gewall/short-url/internal/domain"
	"github.com/gewall/short-url/internal/dto"
	"github.com/gewall/short-url/internal/worker"

	"github.com/gewall/short-url/pkg"
	"github.com/oschwald/geoip2-golang/v2"
)

type ClickRepository interface {
	Create(context.Context, domain.Clicks) error
}

type RedirectService struct {
	geoip    *geoip2.Reader
	LinkRepo LinkRepository
	repo     ClickRepository
	cw       *worker.ClickWorker
}

func NewRedirectService(geoip *geoip2.Reader, link LinkRepository, repo ClickRepository, cw *worker.ClickWorker) *RedirectService {
	return &RedirectService{geoip: geoip, LinkRepo: link, repo: repo, cw: cw}
}

func (s *RedirectService) Redirect(ctx context.Context, redirect *dto.Redirect) (*domain.Link, error) {
	ip, err := netip.ParseAddr(redirect.IP)
	if err != nil {
		return nil, err
	}
	rec, err := s.geoip.Country(ip)
	if err != nil {
		return nil, err
	}
	if !rec.HasData() {
		return nil, pkg.ErrRowsEmpty
	}
	ipHash := pkg.GenerateHash(ip.String())
	redirect.Country = rec.Country.ISOCode

	link, err := s.LinkRepo.FindByShortCode(redirect.Code)
	switch {
	case errors.Is(err, pkg.ErrURLNotFound):
		return nil, pkg.ErrURLNotFound
	case err != nil:
		return nil, err
	}
	click := domain.Clicks{
		LinkID:   link.ID,
		IpHash:   ipHash,
		Country:  redirect.Country,
		City:     redirect.Country,
		Device:   redirect.Device,
		Browser:  redirect.Browser,
		Os:       redirect.OS,
		Referrer: redirect.Referer,
	}

	if err := s.cw.Submit(click); err != nil {
		return nil, err
	}

	// if err := s.repo.Create(ctx, click); err != nil {
	// 	return err
	// }

	return link, nil
}

func RedirectProcessJob(ctx context.Context, click domain.Clicks, repo any) error {
	log.Printf("click: %+v", click)
	if err := repo.(ClickRepository).Create(ctx, click); err != nil {
		return err
	}
	return nil
}
