package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"google.golang.org/protobuf/proto"

	"github.com/Sogilis/Voogle/src/pkg/clients"
	contracts "github.com/Sogilis/Voogle/src/pkg/contracts/v1"

	"github.com/Sogilis/Voogle/src/cmd/api/config"
	jsonDTO "github.com/Sogilis/Voogle/src/cmd/api/dto/json"
	"github.com/Sogilis/Voogle/src/cmd/api/dto/protobuf"
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

	err = conn.WriteMessage(websocket.TextMessage, []byte("Connexion is a success."))
	if err != nil {
		log.Error("Cannot send message : ", err)
		return
	}

	msgs, q, err := GetConsumer(&wsh)
	if err != nil {
		log.Error("Could not create Consumer : ", err)
	}

	HandleMessage(&wsh, *q, msgs, conn)
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

var GetConsumer = func(wsh *WSHandler) (<-chan amqp.Delivery, *amqp.Queue, error) {

	q, err := wsh.AmqpExchangerStatus.QueueDeclare()
	if err != nil {
		log.Error("Could not create queue : ", err)
		return nil, nil, err
	}

	msgs, err := wsh.AmqpExchangerStatus.Consume(q)
	if err != nil {
		log.Error("Failed to register a consumer : ", err)
		return nil, nil, err
	}
	return msgs, &q, nil
}

var HandleMessage = func(wsh *WSHandler, q amqp.Queue, msgs <-chan amqp.Delivery, conn *websocket.Conn) {
	// Read message from client
	go func() {
		for {
			// Read message from browser
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Error("Could not read message : ", err)
				return
			}
			err = wsh.AmqpExchangerStatus.QueueBind(q, string(msg))
			if err != nil {
				log.Error("Could not bind queue : ", err)
			}
		}
	}()

	// Transfer message from queue to client
	for {
		for d := range msgs {
			videoProto := &contracts.Video{}
			if err := proto.Unmarshal([]byte(d.Body), videoProto); err != nil {
				log.Error("Fail to unmarshal video event : ", err)
				continue
			}
			video := protobuf.VideoProtobufToVideo(videoProto)
			video.Title = d.RoutingKey
			msg, err := json.Marshal(jsonDTO.VideoToStatusJson(video))
			if err != nil {
				log.Error("Failed to marshall response to front :", err)
			}
			err = conn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Error("Cannot send message : ", err)
				return
			}
		}
	}
}
