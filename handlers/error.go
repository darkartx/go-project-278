package handlers

import "errors"

var (
	ErrorInvalidOriginalUrl = errors.New("invalid original url")
	ErrorInvalidId          = errors.New("invalid id")
	ErrorInvalidShortName   = errors.New("invalid short name")
)

type Error struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
