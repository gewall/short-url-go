package dto

import "net/http"

type AuthRespDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"-"`
}

func (a *AuthRespDTO) Bind(r *http.Request) error {
	return nil
}
