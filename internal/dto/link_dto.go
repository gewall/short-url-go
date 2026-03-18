package dto

import (
	"net/http"
	"time"

	"github.com/gewall/short-url/pkg"
	"github.com/google/uuid"
)

type LinkReqDTO struct {
	OriginalURL string    `json:"original_url" validate:"required,url"`
	ShortCode   string    `json:"short_code" validate:"max=6"`
	Title       string    `json:"title" validate:"required"`
	ExpiresAt   time.Time `json:"expires_at" validate:"required"`
}

type LinkReqUpdateDTO struct {
	Title    string `json:"title" validate:"required"`
	IsActive string `json:"is_active" validate:"boolean"`
}

type LinkRespDTO struct {
	ID          uuid.UUID `json:"id"`
	ShortCode   string    `json:"short_code"`
	OriginalURL string    `json:"original_url"`
	Title       string    `json:"title"`
	ExpiresAt   string    `json:"expires_at"`
	IsActive    bool      `json:"is_active"`
}

func (lreq *LinkReqDTO) Bind(r *http.Request) error {
	if err := pkg.ValidateStruct(lreq); err != nil {
		return err
	}
	return nil
}

func (lreqUp *LinkReqUpdateDTO) Bind(r *http.Request) error {
	if err := pkg.ValidateStruct(lreqUp); err != nil {
		return err
	}
	return nil
}

func (lres *LinkRespDTO) Render(r *http.Request) error {
	return nil
}
