package webserver

import (
	"github.com/gin-gonic/gin"

	"github.com/vodolaz095/dashboard/internal/service"
)

type Transport struct {
	Address           string
	Domain            string
	HeaderForClientIP string
	TrustProxies      []string
	Version           string
	Title             string
	Description       string
	Keywords          []string
	DoIndex           bool
	PathToHeader      string
	PathToFooter      string
	SensorsService    *service.SensorsService
	Debug             bool

	header, footer []byte
	engine         *gin.Engine
}
