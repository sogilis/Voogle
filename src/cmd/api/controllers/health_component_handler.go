package controllers

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

type HealthComponentHandler struct {
}

// HealthComponentHandler godoc
// @Summary Get component health
// @Description Get component health
// @Tags health
// @Produce plain
// @Success 200 {string} string "OK"
// @Router /health [get]
func (v HealthComponentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug("GET HealthComponentHandler")

	_, err := w.Write([]byte("OK"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
