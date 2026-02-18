package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/darkartx/go-project-278/internal"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	db "github.com/darkartx/go-project-278/db/generated"
)

const (
	shortNameMin = 6
	shortNameMax = 10
)

type LinkHandler struct {
	queries *db.Queries
}

func NewLinkHandler(queries *db.Queries) *LinkHandler {
	return &LinkHandler{queries: queries}
}

func (h *LinkHandler) Register(rg *gin.RouterGroup) {
	rg.POST("", h.Create)
	rg.GET("/:id", h.Get)
	rg.GET("", Range(RangeParam{0, 9}), h.List)
	rg.PUT("/:id", h.Update)
	rg.DELETE("/:id", h.Delete)
}

func (h *LinkHandler) List(c *gin.Context) {
	param, exists := c.Get("range")
	if !exists {
		param = RangeParam{0, 9}
	}

	rangeParam := param.(RangeParam)

	var linksCount int64
	var links []db.Link
	var err error

	linksCount, err = h.queries.GetLinkCount(c)
	if err != nil {
		handleDbError(err, c)
		return
	}

	limit := rangeParam.End - rangeParam.Start + 1

	links, err = h.queries.ListLinks(c, db.ListLinksParams{
		Limit:  int32(limit),
		Offset: int32(rangeParam.Start),
	})
	if err != nil {
		handleDbError(err, c)
		return
	}

	result := make([]Link, 0, len(links))

	for _, item := range links {
		result = append(
			result,
			Link{
				Id:          uint64(item.ID),
				OriginalUrl: item.OriginalUrl,
				ShortName:   item.ShortName,
				ShortUrl:    makeShortUrl(item.ShortName, c),
			},
		)
	}

	c.Header("Content-Range", fmt.Sprintf("links %d-%d/%d", rangeParam.Start, rangeParam.End, linksCount))
	c.JSON(http.StatusOK, result)
}

func (h *LinkHandler) Create(c *gin.Context) {
	input, err := parseAndValidateParams(c)

	if err != nil {
		sendError(http.StatusBadRequest, err, c)
		return
	}

	var shortName string
	if len(input.ShortName) > 0 {
		shortName = input.ShortName
	} else {
		shortName = internal.GenerateShortName(shortNameMin, shortNameMax)
	}

	link, err := h.queries.CreateLink(c, db.CreateLinkParams{
		OriginalUrl: input.OriginalUrl,
		ShortName:   shortName,
	})

	if err != nil {
		handleLinkCreateError(err, c)
		return
	}

	c.JSON(http.StatusCreated, Link{
		Id:          uint64(link.ID),
		OriginalUrl: link.OriginalUrl,
		ShortName:   link.ShortName,
		ShortUrl:    makeShortUrl(link.ShortName, c),
	})
}

func (h *LinkHandler) Get(c *gin.Context) {
	id, err := parseId(c)

	if err != nil {
		sendError(http.StatusBadRequest, err, c)
		return
	}

	var link db.Link
	if link, err = h.queries.GetLink(c, int64(id)); err != nil {
		handleDbError(err, c)
		return
	}

	c.JSON(http.StatusOK, Link{uint64(link.ID), link.OriginalUrl, link.ShortName, makeShortUrl(link.ShortName, c)})
}

func (h *LinkHandler) Update(c *gin.Context) {
	id, err := parseId(c)

	if err != nil {
		sendError(http.StatusBadRequest, err, c)
		return
	}

	var input LinkParams
	input, err = parseAndValidateParams(c)

	if err != nil {
		sendError(http.StatusUnprocessableEntity, err, c)
		return
	}

	var shortName string
	if len(input.ShortName) > 0 {
		shortName = input.ShortName
	} else {
		shortName = internal.GenerateShortName(shortNameMin, shortNameMax)
	}

	var link db.Link
	link, err = h.queries.UpdateLink(c, db.UpdateLinkParams{ID: int64(id), OriginalUrl: input.OriginalUrl, ShortName: shortName})
	if err != nil {
		handleDbError(err, c)
		return
	}

	c.JSON(http.StatusOK, Link{id, link.OriginalUrl, link.ShortName, makeShortUrl(link.ShortName, c)})
}

func (h *LinkHandler) Delete(c *gin.Context) {
	id, err := parseId(c)

	if err != nil {
		sendError(http.StatusBadRequest, err, c)
		return
	}

	if _, err = h.queries.GetLink(c, int64(id)); err != nil {
		handleDbError(err, c)
		return
	}

	if err = h.queries.DeleteLink(c, int64(id)); err != nil {
		handleDbError(err, c)
		return
	}

	c.Status(http.StatusNoContent)
}

func parseAndValidateParams(c *gin.Context) (LinkParams, error) {
	var params LinkParams

	if err := c.ShouldBindJSON(&params); err != nil {
		return LinkParams{}, err
	}

	return params, nil
}

func handleLinkCreateError(err error, c *gin.Context) {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		// Unique constraint
		if pgErr.Code == pgerrcode.UniqueViolation {
			sendError(http.StatusUnprocessableEntity, ErrorShortNameAlreadyUsed, c)
			return
		}
	}

	handleDbError(err, c)
}
