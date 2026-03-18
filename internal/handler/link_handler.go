package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gewall/short-url/internal/dto"
	"github.com/gewall/short-url/internal/service"
	"github.com/gewall/short-url/pkg"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type LinkHandler struct {
	service *service.LinkService
}

func NewLinkHandler(service *service.LinkService) *LinkHandler {
	return &LinkHandler{service: service}
}

func (h *LinkHandler) CreateLink(w http.ResponseWriter, r *http.Request) {
	var link dto.LinkReqDTO
	userId := r.Context().Value("userId").(string)
	if userId == "" {
		render.Render(w, r, pkg.Unauthorized(pkg.ErrUnauthorized))
		return
	}
	if err := render.Bind(r, &link); err != nil {
		render.Render(w, r, pkg.InvalidInput(err))
		return
	}
	l, err := h.service.CreateLink(link, userId)
	if err != nil {
		render.Render(w, r, pkg.InternalServerError(err))
		return
	}

	render.Render(w, r, pkg.CreateResponse(3, "Link created successfully", l))
}

func (h *LinkHandler) FindById(w http.ResponseWriter, r *http.Request) {
	linkIdRaw := chi.URLParam(r, "id")
	linkId, err := uuid.Parse(linkIdRaw)
	if linkIdRaw == "" || err != nil {
		render.Render(w, r, pkg.InvalidInput(fmt.Errorf("link id is required")))
		return
	}
	l, err := h.service.FindById(linkId)
	switch {
	case errors.Is(err, pkg.ErrRowsEmpty):
		render.Render(w, r, pkg.GetResponse(3, "No link found", nil))
		return
	case err != nil:
		render.Render(w, r, pkg.InternalServerError(err))
		return
	}
	render.Render(w, r, pkg.GetResponse(3, "Link found successfully", l))
}

func (h *LinkHandler) FindAllByUserId(w http.ResponseWriter, r *http.Request) {
	userIdRaw := r.Context().Value("userId").(string)
	userId, err := uuid.Parse(userIdRaw)
	if userIdRaw == "" || err != nil {
		render.Render(w, r, pkg.Unauthorized(pkg.ErrUnauthorized))
		return
	}
	l, err := h.service.FindAllByUser(userId)
	switch {
	case errors.Is(err, pkg.ErrRowsEmpty):
		render.Render(w, r, pkg.GetResponse(3, "No links found", nil))
		return
	case err != nil:
		render.Render(w, r, pkg.InternalServerError(err))
		return
	}
	render.Render(w, r, pkg.GetResponse(3, "Links found successfully", l))
}

func (h *LinkHandler) UpdateLink(w http.ResponseWriter, r *http.Request) {
	var link dto.LinkReqUpdateDTO
	linkIdRaw := chi.URLParam(r, "id")
	linkId, err := uuid.Parse(linkIdRaw)
	if linkIdRaw == "" || err != nil {
		render.Render(w, r, pkg.InvalidInput(fmt.Errorf("link id is required")))
		return
	}
	if err := render.Bind(r, &link); err != nil {
		render.Render(w, r, pkg.InvalidInput(err))
		return
	}
	l, err := h.service.UpdateLink(linkId, &link)
	switch {
	case errors.Is(err, pkg.ErrRowsEmpty):
		render.Render(w, r, pkg.GetResponse(3, "Link not found", nil))
		return
	case err != nil:
		render.Render(w, r, pkg.InternalServerError(err))
		return
	}
	render.Render(w, r, pkg.GetResponse(3, "Link updated successfully", l))
}

func (h *LinkHandler) DeleteLink(w http.ResponseWriter, r *http.Request) {
	linkIdRaw := chi.URLParam(r, "id")
	linkId, err := uuid.Parse(linkIdRaw)
	if linkIdRaw == "" || err != nil {
		render.Render(w, r, pkg.InvalidInput(fmt.Errorf("link id is required")))
		return
	}
	if err := h.service.DeleteLink(linkId); err != nil {
		render.Render(w, r, pkg.InternalServerError(err))
		return
	}
	render.Render(w, r, pkg.DeleteResponse(3, "Link deleted successfully", nil))
}
