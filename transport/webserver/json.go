package webserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (tr *Transport) exposeJSON() {
	tr.engine.GET("/json", func(c *gin.Context) {
		sensors := tr.SensorsService.List()
		c.JSON(http.StatusOK, gin.H{
			"sensors": sensors,
		})
	})
}
