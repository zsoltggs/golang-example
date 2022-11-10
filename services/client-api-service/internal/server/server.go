package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/zsoltggs/golang-example/pkg/pds"
	"net/http"
)

type Server struct {
	pdsService pds.ServiceClient
}

func New(pdsService pds.ServiceClient) *Server {
	return &Server{
		pdsService: pdsService,
	}
}

func (s Server) GetPortByID(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		response.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	vars := mux.Vars(request)
	id, ok := vars["id"]
	if !ok {
		log.Info("missing id parameter")
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx := context.Background()
	port, err := s.getPortByID(ctx, id)
	if err != nil {
		log.WithError(err).Info("unable to get port by id")
		// TODO Handle 404 better
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	portJson, err := json.Marshal(port)
	if err != nil {
		log.WithError(err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	_, err = response.Write(portJson)
	if err != nil {
		log.WithError(err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s Server) getPortByID(ctx context.Context, ID string) (*pds.Port, error) {
	res, err := s.pdsService.GetPortByID(ctx, &pds.GetPortByIDRequest{Id: ID})
	if err != nil {
		return nil, fmt.Errorf("unable to query port by id: %w", err)
	}
	return res.GetPort(), nil
}
