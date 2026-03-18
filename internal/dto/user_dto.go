package dto

import (
	"net/http"
	"time"

	"github.com/gewall/short-url/pkg"
	"github.com/google/uuid"
)

type UserReqDTO struct {
	Username string `json:"username" validate:"required,min=3,max=255"`
	Password string `json:"password" validate:"required,min=6,max=255"`
}

type UserRespDTO struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *UserReqDTO) Render(r *http.Request) error {
	return nil
}

func (u *UserReqDTO) Bind(r *http.Request) error {
	if err := pkg.ValidateStruct(u); err != nil {
		return err
	}
	return nil
}

func (u *UserRespDTO) Render(r *http.Request) error {
	return nil
}
