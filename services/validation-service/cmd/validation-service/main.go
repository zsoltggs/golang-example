package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/zsoltggs/golang-example/pkg/users"
	"github.com/zsoltggs/golang-example/services/validation-service/internal/database"
	"github.com/zsoltggs/golang-example/services/validation-service/internal/service"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	cli "github.com/jawher/mow.cli"
	log "github.com/sirupsen/logrus"
)

/*
- REST Service
	- Upload JSON Schemas
		- store them at unique URI
	- Validate JSON documents against these URIs
- Clean document before validation
	- Remove keys for which the value is null

POST    /schema/SCHEMAID        - Upload a JSON Schema with unique `SCHEMAID`
GET     /schema/SCHEMAID        - Download a JSON Schema with unique `SCHEMAID`
POST    /validate/SCHEMAID      - Validate a JSON document against the JSON Schema identified by `SCHEMAID`

Use Case:
- Upload schema
	- /schema/SCHEMAID
	- response: {"action": "uploadSchema", "id": "config-schema", "status": "success"} 201
- Validate
	- /validate/SCHEMAID
	- Clean DOC
	- {"action": "validateDocument", "id": "config-schema", "status": "success"} 200
*/

func main() {
	app := cli.App("validation-service", "CRUD api for users")
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
		Value:  "users",
	})

	restPort := app.Int(cli.IntOpt{
		Name:   "rest-port",
		Desc:   "rest api port for health check",
		Value:  8080,
		EnvVar: "REST_PORT",
	})

	app.Action = func() {
		fmt.Println("starting server")

		db, err := database.NewMongo(*mongoConnStr, *mongoDatabase)
		if err != nil {
			log.WithError(err).Panic("unable to connect to mongo")
		}
		usersService := service.New(db)

		grpcServer := grpc.NewServer()
		users.RegisterServiceServer(grpcServer, usersService)
		startGRPCServer(grpcServer, *grpcPort)
		defer grpcServer.GracefulStop()

		ctx := context.Background()
		port := fmt.Sprintf(":%d", *restPort)
		log.Infof("about to start server on port %s", port)
		router := mux.NewRouter()
		// TODO
		router.HandleFunc("/health", nil).
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
