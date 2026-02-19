package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrorInvalidId            = errors.New("invalid id")
	ErrorShortNameAlreadyUsed = errors.New("short name already in use")
	ErrorInvalidRange         = errors.New("invalid range param")
	ErrorInvalidRequest       = errors.New("invalid request")
)

type ErrorFieldErrors struct {
	Errors map[string]error
}

func (e ErrorFieldErrors) Error() string {
	return "invalid request"
}

func (e *ErrorFieldErrors) Add(field string, err error) {
	e.Errors[field] = err
}

func NewErrorFieldErrors() ErrorFieldErrors {
	errors := make(map[string]error)
	return ErrorFieldErrors{errors}
}

func sendError(code int, err error, c *gin.Context) {
	var fieldErrors ErrorFieldErrors
	var result Error

	if errors.As(err, &fieldErrors) {
		result.Errors = make(map[string]string)
		for field, fe := range fieldErrors.Errors {
			result.Errors[field] = fe.Error()
		}
	} else {
		result.Error = err.Error()
	}

	c.JSON(code, result)
}

func handleDbError(err error, c *gin.Context) {
	if err == nil {
		return
	}

	if errors.Is(err, sql.ErrNoRows) {
		sendNotFound(c)
		return
	}

	sendServerError(c)
}

func sendNotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, Error{
		Error: "Not found",
	})
}

func sendServerError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, Error{
		Error: "Something went wrong",
	})
}
