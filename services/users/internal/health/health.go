package health

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type HealthCheckable interface {
	Health() error
}

type HealthService interface {
	Health() error
	HttpHandler(response http.ResponseWriter, request *http.Request)
}

type healthService struct {
	svc HealthCheckable
}

func New(svc HealthCheckable) HealthService {
	return &healthService{
		svc: svc,
	}
}

func (h healthService) Health() error {
	return h.svc.Health()
}

type healthResponse struct {
	Status string `json:"status"`
}

func (s healthService) HttpHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		response.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	err := s.Health()
	if err != nil {
		log.WithError(err).Infof("service is unhealthy")
		response.WriteHeader(http.StatusServiceUnavailable)
	}
	resp := healthResponse{
		Status: "OK",
	}
	respJSON, err := json.Marshal(resp)
	if err != nil {
		log.WithError(err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	_, err = response.Write(respJSON)
	if err != nil {
		log.WithError(err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
}
