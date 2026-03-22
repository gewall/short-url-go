package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/gewall/short-url/internal/dto"
	"github.com/gewall/short-url/internal/service"
	"github.com/gewall/short-url/pkg"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/mssola/user_agent"
)

type RedirectHandler struct {
	service *service.RedirectService
}

func NewRedirectHandler(service *service.RedirectService) *RedirectHandler {
	return &RedirectHandler{service: service}
}

func (h *RedirectHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var redirect dto.Redirect

	ua := user_agent.New(r.UserAgent())
	code := chi.URLParam(r, "code")
	if code == "" {
		render.Render(w, r, pkg.URLNotFound(pkg.ErrURLNotFound))
		return
	}
	device := "mobile"
	if ua.Mobile() == true {
		device = "mobile"
	} else {
		device = "desktop"
	}
	browser, _ := ua.Browser()
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.Header.Get("X-Real-IP")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}

	redirect = dto.Redirect{
		Code:    chi.URLParam(r, "code"),
		IP:      ip,
		Device:  device,
		Browser: browser,
		OS:      ua.OS(),
		Referer: r.Referer(),
	}

	link, err := h.service.Redirect(ctx, &redirect)
	switch {
	case errors.Is(err, pkg.ErrUserNotFound):
		render.Render(w, r, pkg.URLNotFound(pkg.ErrURLNotFound))
		return
	case err != nil:
		render.Render(w, r, pkg.InternalServerError(err))
		return

	}

	// render.Render(w, r, pkg.GetResponse(4, "redirect", nil))

	http.Redirect(w, r, link.OriginalURL, http.StatusFound)
}
