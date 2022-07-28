package controllers

import (
	"net/http"

	"github.com/Sogilis/Voogle/src/pkg/clients"
)

type VideoTransformerListHandler struct {
	ServiceDiscovery clients.ServiceDiscovery
}

func (v VideoTransformerListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("ok"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
