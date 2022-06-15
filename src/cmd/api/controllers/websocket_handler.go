package controllers

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/streadway/amqp"

	"github.com/Sogilis/Voogle/src/pkg/clients"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"

	"github.com/Sogilis/Voogle/src/cmd/api/config"
)

type WSHandler struct {
	Config              config.Config
	AmqpExchangerStatus clients.IAmqpExchanger
}

func (wsh WSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	upgrader := websocket.Upgrader{}

	upgrader.CheckOrigin = func(r *http.Request) bool {

		decodedData, err := decodeAuthorization(r)
		if err != nil {
			log.Error("Could not decode data", err)
			return false
		}

		givenUser, givenPass := extractCredentials(decodedData)

		// We check the data match our expectation.
		if givenUser == wsh.Config.UserAuth &&
			givenPass == wsh.Config.PwdAuth {
			return true
		}
		return false
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("Cannot upgrade : ", err)
		return
	}
	defer conn.Close()

	err = conn.WriteMessage(websocket.TextMessage, []byte("Connection is a success."))
	if err != nil {
		log.Error("Cannot send message : ", err)
		return
	}

	msgs, err := GetConsumer(&wsh)
	if err != nil {
		log.Error("Could not create Consumer : ", err)
	}

	ReadMessage(msgs, conn)
}

func decodeAuthorization(r *http.Request) (decodedData []byte, err error) {
	authCookie, err := r.Cookie("Authorization")
	if err != nil {
		log.Error("No such cookie", err)
	}
	strCookie := authCookie.String()
	auth := strCookie[len("Authorization="):]
	log.Debug(auth)
	decodedData, err = base64.StdEncoding.DecodeString(auth[len("Basic%20"):])
	if err != nil {
		return nil, err
	}

	return decodedData, nil
}

func extractCredentials(data []byte) (username string, password string) {
	creds := bytes.SplitN(data, []byte(":"), 2)
	givenUser := string(creds[0])
	givenPass := string(creds[1])
	return givenUser, givenPass
}

var GetConsumer = func(wsh *WSHandler) (<-chan amqp.Delivery, error) {

	q, err := wsh.AmqpExchangerStatus.QueueDeclare()
	if err != nil {
		log.Error("Could not create queue : ", err)
		return nil, err
	}

	msgs, err := wsh.AmqpExchangerStatus.Consume(q)
	if err != nil {
		log.Error("Failed to register a consumer : ", err)
		return nil, err
	}
	return msgs, nil
}

var ReadMessage = func(msgs <-chan amqp.Delivery, conn *websocket.Conn) {
	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
		}
	}()

	for {
		// Read message from browser
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Error("Could not read message : ", err)
			return
		}

		// Print the message to the console
		fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))
	}
}
