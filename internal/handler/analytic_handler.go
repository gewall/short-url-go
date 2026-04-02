package handler

import (
	"fmt"
	"net/http"

	"github.com/gewall/short-url/internal/service"
	"github.com/gewall/short-url/pkg"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type AnalyticHandler struct {
	service *service.AnalyticService
}

func NewAnalyticHandler(service *service.AnalyticService) *AnalyticHandler {
	return &AnalyticHandler{service: service}
}

func (h *AnalyticHandler) Analytics(w http.ResponseWriter, r *http.Request) {
	linkIdRaw := chi.URLParam(r, "id")
	linkId, err := uuid.Parse(linkIdRaw)
	if linkIdRaw == "" || err != nil {
		render.Render(w, r, pkg.InvalidInput(fmt.Errorf("link id is required")))
		return
	}

	link, err := h.service.Analytics(linkId)
	if err != nil {
		render.Render(w, r, pkg.InternalServerError(err))
		return
	}

	render.Render(w, r, pkg.GetResponse(5, "Analytics Found", link))
}

func (h *AnalyticHandler) TimeSeries(w http.ResponseWriter, r *http.Request) {
	linkIdRaw := chi.URLParam(r, "id")
	linkId, err := uuid.Parse(linkIdRaw)
	if linkIdRaw == "" || err != nil {
		render.Render(w, r, pkg.InvalidInput(fmt.Errorf("link id is required")))
		return
	}

	series, err := h.service.TimeSeries(linkId)
	if err != nil {
		render.Render(w, r, pkg.InternalServerError(err))
		return
	}

	render.Render(w, r, pkg.GetResponse(5, "Time Series Found", series))
}

func (h *AnalyticHandler) Country(w http.ResponseWriter, r *http.Request) {
	linkIdRaw := chi.URLParam(r, "id")
	linkId, err := uuid.Parse(linkIdRaw)
	if linkIdRaw == "" || err != nil {
		render.Render(w, r, pkg.InvalidInput(fmt.Errorf("link id is required")))
		return
	}

	series, err := h.service.Country(linkId)
	if err != nil {
		render.Render(w, r, pkg.InternalServerError(err))
		return
	}

	render.Render(w, r, pkg.GetResponse(5, "Country Analytics Found", series))
}
