package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrorInvalidOriginalUrl   = errors.New("invalid original url")
	ErrorInvalidId            = errors.New("invalid id")
	ErrorInvalidShortName     = errors.New("invalid short name")
	ErrorShortNameAlreadyUsed = errors.New("short name is already used")
	ErrorInvalidRange         = errors.New("invalid range param")
)

func sendError(code int, err error, c *gin.Context) {
	message := ""

	if err != nil {
		message = err.Error()
	}

	c.JSON(code, Error{
		Error:   http.StatusText(code),
		Message: message,
	})
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
		Error:   http.StatusText(http.StatusNotFound),
		Message: "Not found",
	})
}

func sendServerError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, Error{
		Error:   http.StatusText(http.StatusInternalServerError),
		Message: "Something went wrong",
	})
}
