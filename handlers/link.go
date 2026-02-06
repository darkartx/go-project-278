package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
	input, err := h.parseAndValidateParams(c)

	if err != nil {
		handleError(http.StatusBadRequest, err, c)
		return
	}

	baseUrl := getBaseUrl(c)

	originalUrl := input.OriginalUrl
	shortName := "test"
	shortUrl := baseUrl + "/r/" + shortName

	c.JSON(http.StatusCreated, Link{1, originalUrl, shortName, shortUrl})
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
		input, err = h.parseAndValidateParams(c)
	}

	if err != nil {
		handleError(http.StatusBadRequest, err, c)
		return
	}

	baseUrl := getBaseUrl(c)

	originalUrl := input.OriginalUrl
	shortName := "test"
	shortUrl := baseUrl + "/r/" + shortName

	c.JSON(http.StatusCreated, Link{id, originalUrl, shortName, shortUrl})
}

func (h *LinkHandler) Delete(c *gin.Context) {
	_, err := parseId(c)

	if err != nil {
		handleError(http.StatusBadRequest, err, c)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *LinkHandler) parseAndValidateParams(c *gin.Context) (LinkParams, error) {
	var params LinkParams

	if err := c.ShouldBindJSON(&params); err != nil {
		return LinkParams{}, err
	}

	return params, nil
}

func parseId(c *gin.Context) (int64, error) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
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
	scheme := "http://"

	if c.Request.TLS != nil {
		scheme = "https://"
	}

	return scheme + c.Request.Host
}
