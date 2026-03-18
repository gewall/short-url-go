package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gewall/short-url/internal/dto"
	"github.com/gewall/short-url/internal/service"
	"github.com/gewall/short-url/pkg"
	"github.com/go-chi/render"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var user dto.UserReqDTO
	if err := render.Bind(r, &user); err != nil {
		render.Render(w, r, pkg.InvalidInput(err))
		return
	}

	if err := h.service.SignUp(&user); err != nil {
		render.Render(w, r, pkg.InvalidInput(err))
		return
	}

	render.Render(w, r, pkg.CreateResponse(002, "Sign up successful", nil))
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var user dto.UserReqDTO
	if err := render.Bind(r, &user); err != nil {
		render.Render(w, r, pkg.InvalidInput(err))
		return
	}

	token, err := h.service.SignIn(&user)
	if errors.Is(err, pkg.ErrInvalidPassOrUsn) {
		render.Render(w, r, pkg.InvalidPasswordOrUsername(err))
		return
	}
	if err != nil {
		render.Render(w, r, pkg.InvalidInput(err))
		return
	}

	refCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    token.RefreshToken,
		Expires:  time.Now().Add(time.Hour * 24 * 3),
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &refCookie)
	render.Render(w, r, pkg.GetResponse(002, "Sign in successful", token))
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {

	refToken, err := r.Cookie("refresh_token")
	if err != nil {
		render.Render(w, r, pkg.InvalidInput(err))
		return
	}

	token, err := h.service.RefreshToken(refToken.Value)
	if err != nil {
		render.Render(w, r, pkg.Unauthorized(err))
		return
	}

	refCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    token.RefreshToken,
		Expires:  time.Now().Add(time.Hour * 24 * 3),
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &refCookie)

	render.Render(w, r, pkg.CreateResponse(002, "Refresh token successful", token))
}
