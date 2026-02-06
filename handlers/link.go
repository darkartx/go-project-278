package handlers

import (
	"code/internal"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const (
	shortNameMin = 6
	shortNameMax = 10
)

type LinkHandler struct {
}

func NewLinkHandler() *LinkHandler {
	return &LinkHandler{}
}

func (h *LinkHandler) Register(rg *gin.RouterGroup) {
	rg.POST("", h.Create)
	rg.GET("/:id", h.Get)
	rg.GET("", h.List)
	rg.PUT("/:id", h.Update)
	rg.DELETE("/:id", h.Delete)
}

func (h *LinkHandler) List(c *gin.Context) {
	baseUrl := getBaseUrl(c)

	links := []Link{
		{1, "http://google.com", "test", baseUrl + "/r/test"},
		{2, "http://google.com", "test2", baseUrl + "/r/test2"},
	}

	c.JSON(http.StatusOK, links)
}

func (h *LinkHandler) Create(c *gin.Context) {
	input, err := parseAndValidateParams(c)

	if err != nil {
		handleError(http.StatusBadRequest, err, c)
		return
	}

	baseUrl := getBaseUrl(c)

	var shortName string
	if len(input.ShortName) > 0 {
		shortName = input.ShortName
	} else {
		shortName = internal.GenerateShortName(shortNameMin, shortNameMax)
	}

	shortUrl := fmt.Sprint(baseUrl, "/r/", shortName)

	c.JSON(http.StatusCreated, Link{1, input.OriginalUrl, shortName, shortUrl})
}

func (h *LinkHandler) Get(c *gin.Context) {
	id, err := parseId(c)

	if err != nil {
		handleError(http.StatusBadRequest, err, c)
		return
	}

	baseUrl := getBaseUrl(c)

	c.JSON(http.StatusOK, Link{id, "https://google.com", "test", baseUrl + "/r/test"})
}

func (h *LinkHandler) Update(c *gin.Context) {
	var input LinkParams

	id, err := parseId(c)

	if err == nil {
		input, err = parseAndValidateParams(c)
	}

	if err != nil {
		handleError(http.StatusBadRequest, err, c)
		return
	}

	baseUrl := getBaseUrl(c)

	var shortName string
	if len(input.ShortName) > 0 {
		shortName = input.ShortName
	} else {
		shortName = internal.GenerateShortName(shortNameMin, shortNameMax)
	}

	shortUrl := fmt.Sprint(baseUrl, "/r/", shortName)

	c.JSON(http.StatusOK, Link{id, input.OriginalUrl, shortName, shortUrl})
}

func (h *LinkHandler) Delete(c *gin.Context) {
	_, err := parseId(c)

	if err != nil {
		handleError(http.StatusBadRequest, err, c)
		return
	}

	c.Status(http.StatusNoContent)
}

func parseAndValidateParams(c *gin.Context) (LinkParams, error) {
	var params LinkParams

	if err := c.ShouldBindJSON(&params); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)

		if ok && len(validationErrors) > 0 {
			firstError := validationErrors[0]
			switch firstError.StructField() {
			case "OriginalUrl":
				err = ErrorInvalidOriginalUrl
			case "ShortName":
				err = ErrorInvalidShortName
			}
		}

		return LinkParams{}, err
	}

	return params, nil
}

func parseId(c *gin.Context) (uint64, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		return 0, ErrorInvalidId
	}

	return id, nil
}

func handleError(code int, err error, c *gin.Context) {
	c.JSON(code, Error{
		Error:   http.StatusText(code),
		Message: err.Error(),
	})
}

func getBaseUrl(c *gin.Context) string {
	scheme := "http"

	if c.Request.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	} else if c.Request.TLS != nil {
		scheme = "https"
	}

	return fmt.Sprintf("%s://%s", scheme, c.Request.Host)
}
