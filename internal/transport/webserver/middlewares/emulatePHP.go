package middlewares

import (
	"github.com/gin-gonic/gin"
)

// EmulatePHP makes things more funny
func EmulatePHP(router *gin.Engine) {
	router.Use(func(c *gin.Context) {
		c.Header("X-Powered-By", "PHP/5.6.14") //https://schd.io/ENp
	})
}
