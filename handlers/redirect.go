package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	db "code/db/generated"
)

type RedirectHandler struct {
	queries *db.Queries
}

func NewRedirectHandler(queries *db.Queries) *RedirectHandler {
	return &RedirectHandler{queries: queries}
}

func (h *RedirectHandler) Register(r *gin.Engine) {
	r.GET("/r/:code", RecordVisit(h.queries), h.Get)
}

func (h *RedirectHandler) Get(c *gin.Context) {
	shortName := c.Param("code")

	link, err := h.queries.GetLinkByShortName(c, shortName)

	if err != nil {
		handleDbError(err, c)
		return
	}

	c.Set("link", link)
	c.Redirect(http.StatusFound, link.OriginalUrl)
}
