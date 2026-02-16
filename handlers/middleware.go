package handlers

import (
	db "code/db/generated"
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
)

var rangeRegex *regexp.Regexp = regexp.MustCompile(`^\[(\d+),\s*(\d+)\]$`)

type RangeParam struct {
	Start int
	End   int
}

func Range(def RangeParam) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error

		rangeQueryParam := c.Query("range")

		result := def

		if rangeQueryParam != "" {
			result, err = parseRange(rangeQueryParam)
			if err != nil {
				sendError(http.StatusBadRequest, err, c)
				c.Abort()
				return
			}
		}

		c.Set("range", result)
		c.Next()
	}
}

func parseRange(query string) (RangeParam, error) {
	var temp int
	var result RangeParam
	var err error

	matches := rangeRegex.FindStringSubmatch(query)

	if len(matches) != 3 {
		return RangeParam{}, ErrorInvalidRange
	}

	temp, err = strconv.Atoi(matches[1])
	if err != nil {
		return RangeParam{}, ErrorInvalidRange
	}
	result.Start = temp

	temp, err = strconv.Atoi(matches[2])
	if err != nil {
		return RangeParam{}, ErrorInvalidRange
	}
	result.End = temp

	if result.End < result.Start {
		return RangeParam{}, ErrorInvalidRange
	}

	return result, nil
}

func RecordVisit(queries *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		link, exists := c.Get("link")

		if !exists {
			return
		}

		linkId := link.(db.Link).ID
		ip := getClientIpByHeader(c)
		userAgent := c.Request.UserAgent()
		referer := c.Request.Header.Get("Referer")
		status := c.Writer.Status()

		visit, err := queries.CreateVisit(c, db.CreateVisitParams{
			LinkID:    linkId,
			Ip:        sql.NullString{String: ip, Valid: ip != ""},
			UserAgent: sql.NullString{String: userAgent, Valid: userAgent != ""},
			Referer:   sql.NullString{String: referer, Valid: referer != ""},
			Status:    int16(status),
		})

		if err != nil {
			handleDbError(err, c)
		}

		fmt.Println("~~~~~~~~~~~~~~~~~~~~~~")
		fmt.Printf("%v\n", visit)
	}
}

func getClientIpByHeader(c *gin.Context) string {
	headers := [3]string{
		"X-Forwarded-For",
		"x-forwarded-for",
		"X-FORWARDED-FOR",
	}
	var result string

	for _, header := range headers {
		result = c.Request.Header.Get(header)
		if result != "" {
			return result
		}
	}

	return c.ClientIP()
}
