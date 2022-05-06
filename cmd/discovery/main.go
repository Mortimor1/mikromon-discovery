package main

import (
	"context"
	"github.com/Mortimor1/mikromon-discovery/internal/config"
	"github.com/Mortimor1/mikromon-discovery/internal/subnet"
	"github.com/Mortimor1/mikromon-discovery/internal/webserver"
	"github.com/Mortimor1/mikromon-discovery/pkg/client/mongodb"
	"github.com/Mortimor1/mikromon-discovery/pkg/logging"
)

func main() {
	logger := logging.GetLogger()
	cfg := config.GetConfig()

	// Init DB
	logger.Info("Connect to database")
	client, err := mongodb.NewClient(context.TODO(), cfg.Db.Mongo.Url, cfg.Db.Mongo.Database)
	if err != nil {
		logger.Fatal(err)
	}

	subnetRepo := subnet.NewSubnetRepository(client.Collection("subnet"))

	server := webserver.NewHttpServer(subnetRepo)

	if err := server.Run(cfg); err != nil {
		logger.Fatalf("error running http server: %s", err.Error())
	}
}
