package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	db "github.com/darkartx/go-project-278/db/generated"
)

type LinkVisitHandler struct {
	queries *db.Queries
}

func NewLinkVisitHandler(queries *db.Queries) *LinkVisitHandler {
	return &LinkVisitHandler{queries: queries}
}

func (h *LinkVisitHandler) Register(rg *gin.RouterGroup) {
	rg.GET("", Range(RangeParam{0, 9}), h.List)
}

func (h *LinkVisitHandler) List(c *gin.Context) {
	param, exists := c.Get("range")
	if !exists {
		param = RangeParam{0, 9}
	}

	rangeParam := param.(RangeParam)

	var visitsCount int64
	var visits []db.Visit
	var err error

	visitsCount, err = h.queries.GetVisitCount(c)
	if err != nil {
		handleDbError(err, c)
		return
	}

	limit := rangeParam.End - rangeParam.Start + 1

	visits, err = h.queries.ListVisits(c, db.ListVisitsParams{
		Limit:  int32(limit),
		Offset: int32(rangeParam.Start),
	})
	if err != nil {
		handleDbError(err, c)
		return
	}

	result := make([]Visit, 0, len(visits))

	for _, item := range visits {
		result = append(
			result,
			Visit{
				Id:        uint64(item.ID),
				LinkId:    uint64(item.LinkID),
				Ip:        item.Ip.String,
				UserAgent: item.UserAgent.String,
				Status:    int(item.Status),
				Referer:   item.Referer.String,
				CreatedAt: item.CreatedAt,
			},
		)
	}

	c.Header("Content-Range", fmt.Sprintf("visits %d-%d/%d", rangeParam.Start, rangeParam.End, visitsCount))
	c.JSON(http.StatusOK, result)
}
