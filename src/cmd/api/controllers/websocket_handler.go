package controllers

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"

	"github.com/Sogilis/Voogle/src/pkg/clients"
	contracts "github.com/Sogilis/Voogle/src/pkg/contracts/v1"

	"github.com/Sogilis/Voogle/src/cmd/api/config"
	jsonDTO "github.com/Sogilis/Voogle/src/cmd/api/dto/json"
	"github.com/Sogilis/Voogle/src/cmd/api/dto/protobuf"
)

type WSHandler struct {
	Config                config.Config
	AmqpVideoStatusUpdate clients.AmqpClient
}

// wshandler godoc
// @Summary Send Update to Front
// @Description Send Update to Front
// @Tags websocket
// @Accept plain
// @Produce plain
// @Param Cookie header string true "Authentication cookie"
// @Success 101 {string} string
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 500 {string} string
// @Router /ws [get]
func (wsh WSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	log.Debug("WS WSHandler new connection", r.Host)

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
	randomQueueName := wsh.AmqpVideoStatusUpdate.GetRandomQueueName()
	HandleMessage(context.Background(), &wsh, randomQueueName, conn)
}

var HandleMessage = func(ctx context.Context, wsh *WSHandler, randomQueueName string, conn *websocket.Conn) {

	ctx, clear := context.WithCancel(ctx)

	conn.SetCloseHandler(func(code int, text string) error {
		log.Debugf("Connection closed with code %v : %v", code, text)
		clear()
		return nil
	})

	// Read message from client
	go wsh.handleClientMessage(ctx, clear, randomQueueName, conn)

	// Transfer message from queue to client
	go wsh.handleUpdateMessage(ctx, randomQueueName, conn)

	// Ping Client to ensure connection is still needed
	wsh.pingClient(ctx, clear, conn, time.Duration(5)*time.Second)

	conn.Close()
}

func decodeAuthorization(r *http.Request) (decodedData []byte, err error) {
	authCookie, err := r.Cookie("Authorization")
	if err != nil {
		log.Error("No such cookie", err)
	}
	strCookie := authCookie.String()
	auth := strCookie[len("Authorization="):]
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

func (wsh *WSHandler) handleClientMessage(ctx context.Context, clear context.CancelFunc, randomQueueName string, conn *websocket.Conn) {
	for {
		select {
		case <-ctx.Done():
			return
		default: // Read message from browser
			_, msg, err := conn.ReadMessage()
			if err != nil {
				if _, ok := err.(*websocket.CloseError); ok {
					log.Debug("Close message received.")
					clear()
				} else {
					log.Error("Could not read message : ", err)
				}
			}
			err = wsh.AmqpVideoStatusUpdate.QueueBind(randomQueueName, string(msg))
			if err != nil {
				log.Error("Could not bind queue : ", err)
			}
		}
	}
}

func (wsh *WSHandler) handleUpdateMessage(ctx context.Context, randomQueueName string, conn *websocket.Conn) {
	session := wsh.AmqpVideoStatusUpdate.WithRedial()

	for {
		var client clients.AmqpClient
		select {
		case <-ctx.Done():
			return
		case client = <-session:
			msgs, err := client.Consume(randomQueueName)
			if err != nil {
				log.Error("Failed to consume RabbitMQ client: ", err)
				continue
			}
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
				if err := d.Acknowledger.Ack(d.DeliveryTag, false); err != nil {
					log.Error("Failed to Ack message ", video.ID, " - ", err)
					continue
				}
			}
		}
		client.Close()
	}
}

func (wsh *WSHandler) pingClient(ctx context.Context, clear context.CancelFunc, conn *websocket.Conn, timeout time.Duration) {
	lastCheck := time.Now()
	conn.SetPongHandler(func(appData string) error {
		lastCheck = time.Now()
		return nil
	})
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if time.Now().After(lastCheck.Add(timeout)) {
				clear()
				return
			} else {
				err := conn.WriteMessage(websocket.PingMessage, []byte("pingClient"))
				if err != nil {
					log.Error("Could not ping the client : ", err)
				}
			}
		}
		time.Sleep(timeout * 9 / 10)
	}
}
