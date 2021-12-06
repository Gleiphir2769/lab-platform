package ping

import (
	"errors"
	"lab-platform/core/loadconfig"
	"lab-platform/lib/auth"
	"lab-platform/lib/logger"
	"lab-platform/util"
	"sync"
	"time"


	"go.uber.org/zap"
)

var Ping *ping
var once sync.Once

type ping struct{}

func NewPing() *ping {
	once.Do(func() {
		Ping = &ping{}
	})
	return Ping
}

// Ping the server to make sure the router is working.
func (p *ping) Start() {
	logger.GetLog().Debug("pkg.ping.go:Start()")
	go start()
}

func start() {
	if err := pingServer(); err != nil {
		logger.GetLog().Fatal("The router has no response, or it might took too long to start up.", zap.Error(err))
	}
	logger.GetLog().Info("The router has been deployed successfully.")
}

// pingServer pings the hulk server to make sure the router is working.
func pingServer() error {
	for i := 0; i < loadconfig.Config.MaxPingCount; i++ {
		// Ping the server by sending a GET request to `/health`.
		resp, err := util.HTTPGetBasicAuth(
			loadconfig.Config.URL+"/check/health",
			nil,
			nil,
			util.NewAuth(auth.Config.AdminUsername, auth.Config.AdminPassword))
		if err == nil && resp.StatusCode == 200 {
			return nil
		}

		// Sleep for a second to continue the next ping.
		logger.GetLog().Info("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}
	return errors.New("can not connect to the router")
}
