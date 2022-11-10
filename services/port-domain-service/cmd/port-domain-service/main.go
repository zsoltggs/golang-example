package main

import (
	"fmt"
	"github.com/zsoltggs/golang-example/pkg/pds"
	"github.com/zsoltggs/golang-example/services/port-domain-service/internal/database"
	"github.com/zsoltggs/golang-example/services/port-domain-service/internal/service"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"

	cli "github.com/jawher/mow.cli"
	log "github.com/sirupsen/logrus"
)

func main() {
	app := cli.App("port-domain-service", "this app provides an api for the port domain")
	grpcPort := app.Int(cli.IntOpt{
		Name:   "grpc-port",
		Desc:   "GRPC port",
		Value:  8090,
		EnvVar: "GRPC_PORT",
	})
	mongoConnStr := app.String(cli.StringOpt{
		Name:   "mongo",
		Desc:   "connection string",
		EnvVar: "MONGO",
		Value:  "mongodb://localhost:27017",
	})

	mongoDatabase := app.String(cli.StringOpt{
		Name:   "mongo-database",
		Desc:   "Database name for mongo",
		EnvVar: "MONGO_DB",
		Value:  "port-domain-service",
	})

	app.Action = func() {
		fmt.Println("starting server")

		db, err := database.NewMongo(*mongoConnStr, *mongoDatabase)
		if err != nil {
			log.WithError(err).Panic("unable to connect to mongo")
		}
		pdsService := service.New(db)

		grpcServer := grpc.NewServer()
		pds.RegisterServiceServer(grpcServer, pdsService)
		startGRPCServer(grpcServer, *grpcPort)
		defer grpcServer.GracefulStop()

		waitForShutdown()
	}

	err := app.Run(os.Args)
	if err != nil {
		log.WithError(err).Panic("service stopped")
	}
}

func startGRPCServer(server *grpc.Server, grpcPort int) {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.WithError(err).Panicf("failed to listen on port :%d", grpcPort)
	}

	go func() {
		if err := server.Serve(listen); err != nil {
			log.WithError(err).Panic("failed to serve GRPC connections")
		}
	}()
}

// Graceful shutdown
func waitForShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	log.Warn("shutting down")
}
