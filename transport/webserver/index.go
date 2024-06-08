package webserver

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vodolaz095/dashboard/model"
)

func (tr *Transport) exposeIndex() {
	tr.engine.GET("/", func(c *gin.Context) {
		var sensors []model.Sensor
		tags := c.Request.URL.Query()
		if len(tags) == 0 {
			sensors = tr.SensorsService.List()
		} else {
			needle := make(map[string]string, 0)
			for k := range tags {
				needle[k] = strings.Join(tags[k], " ")
			}
			sensors = tr.SensorsService.ListByTags(needle)
		}
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":       tr.Title,
			"description": tr.Description,
			"keywords":    strings.Join(tr.Keywords, ", "),
			"doIndex":     tr.DoIndex,
			"sensors":     sensors,
			"version":     tr.Version,
			"header":      string(tr.header),
			"footer":      string(tr.footer),
		})
	})
}
