package main

import (
	"github.com/Mortimor1/mikromon-discovery/internal/config"
	"github.com/Mortimor1/mikromon-discovery/internal/discovery"
	"github.com/Mortimor1/mikromon-discovery/internal/webserver"
	"github.com/Mortimor1/mikromon-discovery/pkg/logging"
	"github.com/carlescere/scheduler"
)

func main() {
	logger := logging.GetLogger()
	cfg := config.GetConfig()

	server := new(webserver.Server)

	d := new(discovery.Discovery)

	job := func() {
		logger.Println("Start discovery...")
		go d.Run(cfg)
	}

	_, err := scheduler.Every(int(cfg.Discovery.Interval.Seconds())).Seconds().Run(job)

	if err != nil {
		logger.Fatalf("error running job discovery: %s", err.Error())
	}

	if err := server.Run(cfg); err != nil {
		logger.Fatalf("error running http server: %s", err.Error())
	}
}
