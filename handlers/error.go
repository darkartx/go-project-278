package handlers

import "errors"

var (
	ErrorInvalidOriginalUrl = errors.New("invalid original url")
	ErrorInvalidId          = errors.New("invalid id")
)

type Error struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
