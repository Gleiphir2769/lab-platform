package core

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"lab-platform/core/loadconfig"
	"lab-platform/lib/logger"
	"lab-platform/ping"
	"net/http"
	"runtime"
	"sync"
	"time"
)

var once sync.Once
var Core *core

type core struct {
	server *http.Server
}

func NewCore(g *gin.Engine) *core {
	once.Do(func() {
		Core = &core{
			server: &http.Server{
				Addr:    loadconfig.Config.Addr,
				Handler: g,
			},
		}
	})
	return Core
}

func (c *core) Start() {
	logger.Log.Start()

	// Ping the server to make sure the router is working.
	ping.NewPing().Start()
}

func (c *core) Run() {
	logger.GetLog().Debug("handler.internal.core.go:Run()")

	time.Sleep(300 * time.Millisecond)
	// ConfigRuntime sets the number of operating system threads.
	logger.GetLog().Info("real GOMAXPROCS", zap.String("ConfigRuntime", fmt.Sprintf("Running with %d CPUs", runtime.GOMAXPROCS(-1))))
	logger.GetLog().Info("Start to listening the incoming requests on http address.", zap.String("port", loadconfig.Config.Addr))
	// server start
	go func() {
		if err := c.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.GetLog().Fatal(fmt.Sprintf("listen: %s\n", err))
		}
	}()
}

func (c *core) Shutdown() {
	logger.GetLog().Debug("handler.internal.core.go:Shutdown()")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := c.server.Shutdown(ctx); err != nil {
		logger.GetLog().Fatal(fmt.Sprintf("Server forced to shutdown:%s\n", err))
	}

	logger.GetLog().Info("Server exiting")
}

