package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/zsoltggs/golang-example/client-api-service/internal/parser"
	"github.com/zsoltggs/golang-example/client-api-service/internal/server"
	"github.com/zsoltggs/golang-example/port-domain-service/pkg/generated/github.com/zsoltggs/golang-example/pkg/pds"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	cli "github.com/jawher/mow.cli"
	log "github.com/sirupsen/logrus"
)

func main() {
	app := cli.App("client-api-service", "this app will read the ports.json file")
	inputFileName := app.String(cli.StringOpt{
		Name:   "input-file",
		Desc:   "Input ports file",
		EnvVar: "INPUT_PORTS_FILE",
		Value:  "resources/ports.json",
	})
	portDomainServiceApi := app.String(cli.StringOpt{
		Name:   "port-domain-service-api",
		Desc:   "Port domain service api address",
		EnvVar: "PORT_DOMAIN_SERVICE_API",
		Value:  "localhost:9000",
	})
	restPort := app.Int(cli.IntOpt{
		Name:   "rest-port",
		Desc:   "rest api port",
		Value:  8080,
		EnvVar: "REST_PORT",
	})

	app.Action = func() {
		ctx := context.Background()
		inputFile, err := os.Open(*inputFileName)
		if err != nil {
			log.WithError(err).Panicf("unable to open file with name %q", *inputFileName)
		}

		ctxTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		portDomainConnection := initialiseGRPCConnection(ctxTimeout, portDomainServiceApi)
		defer portDomainConnection.Close()
		pdsService := pds.NewServiceClient(portDomainConnection)

		parserSvc := parser.New(pdsService)
		err = parserSvc.Parse(ctx, inputFile)
		if err != nil {
			log.WithError(err).Panicf("unable to parse file")
		}
		log.Info("ports import finished")

		serviceServer := server.New(pdsService)
		port := fmt.Sprintf(":%d", *restPort)
		log.Infof("about to start server on port %s", port)
		router := mux.NewRouter()
		router.HandleFunc("/port/{id}", serviceServer.GetPortByID).
			Methods("GET")
		httpServer := http.Server{
			Addr:    port,
			Handler: router,
		}
		go func() {
			if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("listen: %s\n", err)
			}
		}()
		defer func() { _ = httpServer.Shutdown(ctx) }()

		log.Info("ready")
		waitForShutdown()

	}

	err := app.Run(os.Args)
	if err != nil {
		log.WithError(err).Panic("service stopped")
	}
}

func initialiseGRPCConnection(ctx context.Context, addr *string) *grpc.ClientConn {
	connection, err := grpc.DialContext(ctx, *addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock())
	if err != nil {
		log.WithError(err).Panicf("failed to dial to %s", *addr)
	}

	return connection
}

// Graceful shutdown
func waitForShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	log.Warn("shutting down")
}
