package pkg

import (
	"net/http"

	"github.com/go-chi/render"
)

type Response struct {
	HttpCode int    `json:"-"`
	Code     int    `json:"code"`
	Message  string `json:"message"`
	Data     any    `json:"data,omitempty"`
}

func (resp *Response) Render(w http.ResponseWriter, r *http.Request) error {

	render.Status(r, resp.HttpCode)

	return nil
}

func CreateResponse(code int, message string, data any) render.Renderer {
	return &Response{
		HttpCode: 201,
		Code:     code,
		Message:  message,
		Data:     data,
	}
}

func GetResponse(code int, message string, data any) render.Renderer {
	return &Response{
		HttpCode: 200,
		Code:     code,
		Message:  message,
		Data:     data,
	}
}

func DeleteResponse(code int, message string, data any) render.Renderer {
	return &Response{
		HttpCode: 204,
		Code:     code,
		Message:  message,
		Data:     data,
	}
}
