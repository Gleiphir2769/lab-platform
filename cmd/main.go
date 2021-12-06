package main

import (
	"lab-platform/core"
	"lab-platform/lib/config"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config.Parse()
	if config.PrintVersion() {
		return
	}

	c, err := core.Init()
	if err != nil {
		panic(err)
	}
	c.Start()
	c.Run()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals
	c.Shutdown()
}
