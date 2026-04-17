package pkg

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrInternal         = errors.New("internal server error")
	ErrInvalidInput     = errors.New("invalid input")
	ErrRowsEmpty        = errors.New("rows empty")
	ErrInvalidPassOrUsn = errors.New("invalid password or username")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrInvalidToken     = errors.New("invalid token")
	ErrURLNotFound      = errors.New("url not found")
	ErrWorkerQueueFull  = errors.New("worker queue full")
)

type Error struct {
	HttpCode int    `json:"-"`
	Code     int    `json:"code"`
	Message  string `json:"message"`
}

func (err *Error) Render(w http.ResponseWriter, r *http.Request) error {

	render.Status(r, err.HttpCode)

	return nil
}

func InvalidInput(err error) render.Renderer {
	return &Error{
		HttpCode: http.StatusBadRequest,
		Code:     400,
		Message:  err.Error(),
	}
}

func NotFound(err error) render.Renderer {
	return &Error{
		HttpCode: http.StatusNotFound,
		Code:     404,
		Message:  err.Error(),
	}
}

func InternalServerError(err error) render.Renderer {
	return &Error{
		HttpCode: http.StatusInternalServerError,
		Code:     500,
		Message:  err.Error(),
	}
}

func InvalidPasswordOrUsername(err error) render.Renderer {
	return &Error{
		HttpCode: http.StatusUnauthorized,
		Code:     401,
		Message:  err.Error(),
	}
}

func Unauthorized(err error) render.Renderer {
	return &Error{
		HttpCode: http.StatusUnauthorized,
		Code:     401,
		Message:  err.Error(),
	}
}

func URLNotFound(err error) render.Renderer {
	return &Error{
		HttpCode: http.StatusNotFound,
		Code:     404,
		Message:  err.Error(),
	}
}
