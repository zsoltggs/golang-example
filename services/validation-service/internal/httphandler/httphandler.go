package httphandler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/zsoltggs/golang-example/services/validation-service/internal/service"
	"github.com/zsoltggs/golang-example/services/validation-service/pkg/validationmodels"
	"io"
	"net/http"
)

type Handler struct {
	svc        service.Service
	httpServer *http.Server
	port       int
}

func NewHandler(svc service.Service, port int) *Handler {
	h := Handler{
		svc:  svc,
		port: port,
	}
	portStr := fmt.Sprintf(":%d", port)
	router := mux.NewRouter()
	router.HandleFunc(validationmodels.GetSchemaPath, h.getSchema).
		Methods("GET")
	router.HandleFunc(validationmodels.PostSchemaPath, h.postSchema).
		Methods("POST")
	router.HandleFunc(validationmodels.PostValidatePath, h.validateDoc).
		Methods("POST")
	httpServer := http.Server{
		Addr:    portStr,
		Handler: router,
	}
	h.httpServer = &httpServer

	return &h
}

/*
POST    /schema/SCHEMAID        - Upload a JSON Schema with unique `SCHEMAID`
GET     /schema/SCHEMAID        - Download a JSON Schema with unique `SCHEMAID`
POST    /validate/SCHEMAID      - Validate a JSON document against the JSON Schema identified by `SCHEMAID`
*/

func (h Handler) ListenAndServe() {
	if err := h.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}

func (h Handler) Shutdown(ctx context.Context) {
	_ = h.httpServer.Shutdown(ctx)
}

func (h Handler) postSchema(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		response.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	vars := mux.Vars(request)
	schemaID, ok := vars["schemaId"]
	if !ok {
		log.Info("missing schemaID parameter")
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(request.Body)
	if err != nil {
		log.WithError(err).Info("error unmarshaling request body")
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx := context.Background()
	resp, err := h.svc.UpsertSchema(ctx, &validationmodels.UpsertSchemaRequest{
		SchemaID: schemaID,
		Schema:   string(body),
	})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidJSON):
			responseJson, err := json.Marshal(resp.HttpResponse)
			if err != nil {
				log.WithError(err)
				response.WriteHeader(http.StatusInternalServerError)
				return
			}

			response.Header().Set("Content-Type", "application/json")
			response.WriteHeader(http.StatusBadRequest)
			_, err = response.Write(responseJson)
			if err != nil {
				log.WithError(err)
				response.WriteHeader(http.StatusInternalServerError)
				return
			}
		default:
			log.WithError(err).Info("unable to validate document")
			response.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	responseJson, err := json.Marshal(resp.HttpResponse)
	if err != nil {
		log.WithError(err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusCreated)
	_, err = response.Write(responseJson)
	if err != nil {
		log.WithError(err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h Handler) getSchema(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		response.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	vars := mux.Vars(request)
	schemaID, ok := vars["schemaId"]
	if !ok {
		log.Info("missing schemaID parameter")
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx := context.Background()
	resp, err := h.svc.GetSchemaByID(ctx, &validationmodels.GetSchemaRequest{
		ID: schemaID,
	})
	if err != nil {
		log.WithError(err).Info("unable to get schema document")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	responseJson, err := json.Marshal(resp)
	if err != nil {
		log.WithError(err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	_, err = response.Write(responseJson)
	if err != nil {
		log.WithError(err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h Handler) validateDoc(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		response.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	vars := mux.Vars(request)
	schemaID, ok := vars["schemaId"]
	if !ok {
		log.Info("missing schemaID parameter")
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(request.Body)
	if err != nil {
		log.WithError(err).Info("error unmarshaling request body")
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx := context.Background()
	resp, err := h.svc.ValidateDocument(ctx, &validationmodels.ValidateRequest{
		SchemaID: schemaID,
		Document: string(body),
	})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidJSON),
			errors.Is(err, service.ErrValidationError):
			responseJson, err := json.Marshal(resp.HttpResponse)
			if err != nil {
				log.WithError(err)
				response.WriteHeader(http.StatusInternalServerError)
				return
			}

			response.Header().Set("Content-Type", "application/json")
			response.WriteHeader(http.StatusBadRequest)
			_, err = response.Write(responseJson)
			if err != nil {
				log.WithError(err)
				response.WriteHeader(http.StatusInternalServerError)
				return
			}
		default:
			log.WithError(err).Info("unable to validate document")
			response.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	responseJson, err := json.Marshal(resp.HttpResponse)
	if err != nil {
		log.WithError(err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	_, err = response.Write(responseJson)
	if err != nil {
		log.WithError(err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
}
