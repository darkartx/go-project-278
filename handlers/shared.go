package handlers

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func parseId(c *gin.Context) (uint64, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		return 0, ErrorInvalidId
	}

	return id, nil
}

func getBaseUrl(c *gin.Context) string {
	result := c.Request.Header.Get("Referer")

	if result != "" {
		return result
	}

	scheme := "http"

	if c.Request.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	} else if c.Request.TLS != nil {
		scheme = "https"
	}

	return fmt.Sprintf("%s://%s/", scheme, c.Request.Host)
}

func makeShortUrl(shortName string, c *gin.Context) string {
	baseUrl := getBaseUrl(c)
	return fmt.Sprint(baseUrl, "r/", shortName)
}
