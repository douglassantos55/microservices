package pkg

import (
	"encoding/json"
)

type Error struct {
	Status  int    `json:"status"`
	Title   string `json:"title"`
	Message string `json:"detail"`
}

func NewError(status int, title, message string) Error {
	return Error{status, title, message}
}

func (e Error) StatusCode() int {
	return int(e.Status)
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
