package controllers

import (
	"fmt"
	"net/http"

	"github.com/Sogilis/Voogle/src/pkg/clients"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type WSHandler struct {
	AmqpExchangerStatus clients.IAmqpExchanger
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

	q, err := wsh.AmqpExchangerStatus.QueueDeclare()
	if err != nil {
		log.Error("Could not create queue : ", err)
		return
	}

	msgs, err := wsh.AmqpExchangerStatus.Consume(q)
	if err != nil {
		log.Error("Failed to register a consumer : ", err)
		return
	}

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
		}
	}()

	for {
		// Read message from browser
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		// Print the message to the console
		fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))
	}
}
