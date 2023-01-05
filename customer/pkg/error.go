package pkg

import (
	"encoding/json"
	"io"
	"net/http"
)

type Error struct {
	Status  int    `json:"status"`
	Title   string `json:"title"`
	Message string `json:"detail"`
}

func NewError(status int, title, message string) Error {
	return Error{status, title, message}
}

func NewErrorFromResponse(body io.ReadCloser) Error {
	var error Error
	json.NewDecoder(body).Decode(&error)
	return error
}

func (e Error) StatusCode() int {
	return e.Status
}

func (e Error) Headers() http.Header {
	return http.Header{
		"Content-Type": []string{
			"application/problem+json",
		},
	}
}

func (e Error) Error() string {
	return e.Title
}

func (e Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"status": e.StatusCode(),
		"title":  e.Title,
		"detail": e.Message,
	})
}
