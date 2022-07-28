package controllers

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	jsonDTO "github.com/Sogilis/Voogle/src/cmd/api/dto/json"
	"github.com/Sogilis/Voogle/src/pkg/clients"
)

type VideoTransformerListHandler struct {
	ServiceDiscovery clients.ServiceDiscovery
}

type TransformerServiceListResponse struct {
	Services []jsonDTO.TransformerServiceJson `json:"services"`
}

// VideoTransformerListHandler godoc
// @Summary Get list of existing services
// @Description Get list of existing services
// @Tags services
// @Produce json
// @Success 200 {object} TransformerServiceListResponse "Service list"
// @Router /api/v1/videos/transformer [get]
func (v VideoTransformerListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug("GET VideoTransformerListHandler")
	existingServices := v.ServiceDiscovery.GetExistingServices()
	services := []jsonDTO.TransformerServiceJson{}
	for _, service := range existingServices {
		services = append(services, jsonDTO.TransformerServiceToTransformerServiceJson(service))
	}
	payload, err := json.Marshal(TransformerServiceListResponse{Services: services})
	if err != nil {
		log.Error("Unable to parse data struct in json ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(payload)
}
