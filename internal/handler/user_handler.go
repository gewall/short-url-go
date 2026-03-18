package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/gewall/short-url/internal/dto"
	"github.com/gewall/short-url/internal/service"
	"github.com/gewall/short-url/pkg"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user dto.UserReqDTO
	err := render.Bind(r, &user)
	if err != nil {
		render.Render(w, r, pkg.InvalidInput(err))
		return
	}

	u, err := h.service.CreateUser(user)
	if err != nil {
		render.Render(w, r, pkg.InternalServerError(err))
		return
	}

	render.Render(w, r, pkg.CreateResponse(001, "User Created", u))
}

func (h *UserHandler) FindUserByID(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "id")
	log.Println(userId)
	id, err := uuid.Parse(userId)
	if err != nil {
		render.Render(w, r, pkg.InvalidInput(err))
		return
	}

	user, err := h.service.FindUserByID(id)
	switch {
	case errors.Is(err, pkg.ErrRowsEmpty):
		render.Render(w, r, pkg.GetResponse(001, "User not found", nil))
		return
	case err != nil:
		render.Render(w, r, pkg.InternalServerError(err))
		return
	}

	render.Render(w, r, pkg.GetResponse(001, "User Found", user))
}

func (h *UserHandler) FindUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.FindAllUsers()
	switch {
	case err != nil:
		render.Render(w, r, pkg.InternalServerError(err))
		return
	case len(users) == 0:
		render.Render(w, r, pkg.GetResponse(001, "No users found", nil))
		return
	}

	render.Render(w, r, pkg.GetResponse(001, "Users Found", users))
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "id")
	id, err := uuid.Parse(userId)
	if err != nil {
		render.Render(w, r, pkg.InternalServerError(err))
		return
	}

	err = h.service.DeleteUser(id)
	if err != nil {
		render.Render(w, r, pkg.NotFound(err))
		return
	}

	render.Render(w, r, pkg.DeleteResponse(001, "User Deleted", nil))
}
