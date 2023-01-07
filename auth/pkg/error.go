package pkg

import (
	"encoding/json"

	"reconcip.com.br/microservices/auth/proto"
)

type Error interface {
	error
	StatusCode() int
	AsError() *proto.Error
}

type errorMessage struct {
	Status  int    `json:"status"`
	Title   string `json:"title"`
	Message string `json:"detail"`
}

func NewError(status int, title, message string) Error {
	return errorMessage{status, title, message}
}

func (e errorMessage) StatusCode() int {
	return e.Status
}

func (e errorMessage) Error() string {
	return e.Title
}

func (e errorMessage) AsError() *proto.Error {
	return &proto.Error{
		Status: uint32(e.StatusCode()),
		Title:  e.Title,
		Detail: e.Message,
	}
}

func (e errorMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"status": e.StatusCode(),
		"title":  e.Title,
		"detail": e.Message,
	})
}
