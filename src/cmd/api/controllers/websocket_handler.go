package controllers

import (
	"net/http"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type WSHandler struct {
}

func (wsh WSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	upgrader := websocket.Upgrader{}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("Cannot upgrade : ", err)
		return
	}

	err = conn.WriteMessage(websocket.TextMessage, []byte("Connection is a success."))
	if err != nil {
		log.Error("Cannot send message : ", err)
		return
	}
}
