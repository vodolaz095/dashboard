package webserver

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/vodolaz095/dashboard/assets"
	"github.com/vodolaz095/dashboard/transport/webserver/middlewares"
)

func (tr *Transport) Start(ctx context.Context, wg *sync.WaitGroup) (err error) {
	defer wg.Done()
	if tr.PathToHeader != "" {
		header, err1 := os.ReadFile(tr.PathToHeader)
		if err1 != nil {
			return fmt.Errorf("error reading header from %s: %w", tr.PathToHeader, err1)
		}
		tr.header = header
	}
	if tr.PathToFooter != "" {
		footer, err2 := os.ReadFile(tr.PathToFooter)
		if err2 != nil {
			return fmt.Errorf("error reading footer from %s: %w", tr.PathToFooter, err2)
		}
		tr.footer = footer
	}
	tr.engine = gin.New()
	tr.engine.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		log.Debug().Msgf("[%s] - \"%s %s %s\" -> code=%d lat=%s size=%d / \"%s\"",
			param.ClientIP,
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.BodySize,
			param.Request.UserAgent(),
		)
		return ""
	}))
	tr.engine.Use(gin.Recovery())
	middlewares.Secure(tr.engine, tr.Domain)
	middlewares.UseCaching(tr.engine)
	middlewares.EmulatePHP(tr.engine)
	tr.engine.TrustedPlatform = gin.PlatformCloudflare
	err = injectTemplates(tr.engine)
	if err != nil {
		return err
	}
	fs := http.FS(assets.Assets)
	tr.engine.StaticFS("/assets", fs)
	tr.engine.GET("/favicon.ico", func(c *gin.Context) {
		c.FileFromFS("favicon.ico", fs)
	})
	tr.engine.GET("/robots.txt", func(c *gin.Context) {
		c.FileFromFS("robots.txt", fs)
	})

	tr.exposeIndex()
	tr.exposeFeed()
	tr.exposeJSON()
	tr.exposeMetrics()
	tr.exposeUpdate()
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
