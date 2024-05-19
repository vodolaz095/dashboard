package webserver

import (
	"context"
	"net"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/vodolaz095/dashboard/service"
	"github.com/vodolaz095/dashboard/transport/webserver/middlewares"
)

type Transport struct {
	Address     string
	Domain      string
	Version     string
	Title       string
	Description string
	Keywords    []string
	DoIndex     bool

	SensorsService *service.SensorsService
	engine         *gin.Engine
}

func (tr *Transport) Start(ctx context.Context, wg *sync.WaitGroup) (err error) {
	defer wg.Done()
	tr.engine = gin.New()
	middlewares.Secure(tr.engine, tr.Domain)
	middlewares.EmulatePHP(tr.engine)

	err = injectTemplates(tr.engine)
	if err != nil {
		return err
	}

	tr.exposeIndex()
	tr.exposeJSON()
	tr.exposeMetrics()
	tr.exposeEndpoint()
	tr.exposeHealthcheck()

	listener, err := net.Listen("tcp", tr.Address)
	if err != nil {
		return
	}
	go func() {
		<-ctx.Done()
		log.Debug().Msgf("Closing HTTP server on %s...", tr.Address)
		listener.Close()
	}()
	log.Info().Msgf("Starting HTTP server on %s", tr.Address)
	wg.Add(1)
	err = tr.engine.RunListener(listener)
	if strings.Contains(err.Error(), "use of closed network connection") {
		return nil
	}
	return
}
