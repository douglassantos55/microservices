package pkg

import (
	"encoding/json"

	"reconcip.com.br/microservices/inventory/proto"
)

type Error struct {
	Status int    `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

func NewError(status int, title, detail string) error {
	return Error{status, title, detail}
}

func (e Error) StatusCode() int {
	return int(e.Status)
}

func (e Error) Error() string {
	return e.Title
}

func NewErrorFromReply(err *proto.Error) Error {
	return Error{
		Status: int(err.GetStatus()),
		Title:  err.GetTitle(),
		Detail: err.GetDetail(),
	}
}

func (e Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"status": e.StatusCode(),
		"title":  e.Title,
		"detail": e.Detail,
	})
}
