package webserver

import (
	"github.com/gin-gonic/gin"
	"github.com/vodolaz095/dashboard/service"
)

type Transport struct {
	Address      string
	Domain       string
	Version      string
	Title        string
	Description  string
	Keywords     []string
	DoIndex      bool
	PathToHeader string
	PathToFooter string

	header, footer []byte
	SensorsService *service.SensorsService
	engine         *gin.Engine
}
