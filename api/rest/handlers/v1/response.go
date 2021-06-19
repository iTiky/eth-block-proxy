package v1

import (
	"net/http"

	"github.com/go-chi/render"
)

// ErrResponse keeps HTTP response meta and used to generate error responses.
type ErrResponse struct {
	HTTPStatusCode int `json:"-"`
	//
	StatusText string `json:"status"`
	ErrorText  string `json:"error,omitempty"`
}

// Render implements render.Renderer interface.
func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// NewErrResponseNotFound return the 404 response.
func NewErrResponseNotFound() render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: http.StatusNotFound,
		StatusText:     "Resource not found",
	}
}

// NewErrResponseInvalidRequest return the 400 response.
func NewErrResponseInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     "Invalid request",
		ErrorText:      err.Error(),
	}
}

// NewErrResponseInternal return the 500 response.
func NewErrResponseInternal(err error) render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: http.StatusInternalServerError,
		StatusText:     "Internal error",
		ErrorText:      err.Error(),
	}
}

// NewErrResponseRender return the 422 response.
func NewErrResponseRender(err error) render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: http.StatusUnprocessableEntity,
		StatusText:     "Error rendering response",
		ErrorText:      err.Error(),
	}
}
