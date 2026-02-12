package handlers

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
)

var rangeRegex *regexp.Regexp = regexp.MustCompile("^\\[(\\d+),\\s*(\\d+)\\]$")

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
