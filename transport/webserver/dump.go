package webserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (tr *Transport) exposeDump() {
	tr.engine.GET("/dump", func(c *gin.Context) {
		tasks := tr.SensorsService.UpdateQueue.Dump()
		c.IndentedJSON(http.StatusOK, tasks)
	})
}
