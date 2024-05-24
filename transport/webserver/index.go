package webserver

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (tr *Transport) exposeIndex() {
	tr.engine.GET("/", func(c *gin.Context) {
		sensors := tr.SensorsService.List()
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":       tr.Title,
			"description": tr.Description,
			"keywords":    strings.Join(tr.Keywords, ", "),
			"doIndex":     tr.DoIndex,
			"sensors":     sensors,
			"version":     tr.Version,
		})
	})
}
