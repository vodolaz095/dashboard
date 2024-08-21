package webserver

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/vodolaz095/dashboard/model"
)

func (tr *Transport) listFilteredSensors(c *gin.Context) (sensors []model.Sensor, filtered bool) {
	tags := c.Request.URL.Query()
	if len(tags) == 0 {
		sensors = tr.SensorsService.List()
	} else {
		needle := make(map[string]string, 0)
		for k := range tags {
			needle[k] = strings.Join(tags[k], " ")
		}
		sensors = tr.SensorsService.ListByTags(needle)
		filtered = true
	}
	return
}
