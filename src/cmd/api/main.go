package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Sogilis/Voogle/src/pkg/clients"
	"github.com/Sogilis/Voogle/src/pkg/events"

	"github.com/Sogilis/Voogle/src/cmd/api/config"
	"github.com/Sogilis/Voogle/src/cmd/api/db/dao"
	"github.com/Sogilis/Voogle/src/cmd/api/eventhandler"
	"github.com/Sogilis/Voogle/src/cmd/api/router"
)

const GORILLA_MUX_SHUTDOWN_TIMEOUT time.Duration = time.Second * 2
const GOROUTINE_FLUSH_TIMEOUT time.Duration = time.Millisecond * 100

func main() {
	log.Info("Starting Voogle API")

	// Retrieve environment variables
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("Failed to parse environment variables : ", err)
	}
	if cfg.DevMode {
		log.SetReportCaller(true)
		log.SetLevel(log.DebugLevel)
	}

	// Create routers
	routerClients, routerDAOs := createRouters(cfg)
	defer routerDAOs.Db.Close()
	defer routerDAOs.VideosDAO.Close()
	defer routerDAOs.UploadsDAO.Close()

	// Start service discovery
	go func() {
		serviceInfos := clients.ServiceInfos{
			Name:    "api",
			Address: cfg.LocalAddr,
			Port:    int(cfg.Port),
			Tags:    nil,
		}
		if err := routerClients.ServiceDiscovery.StartServiceDiscovery(serviceInfos); err != nil {
			log.Fatal("Discovery Service crash : ", err)
		}
	}()

	// Start API server
	log.Info("Starting server on port : ", cfg.Port)
	srv := &http.Server{
		Handler: router.NewRouter(cfg, routerClients, routerDAOs),
		Addr:    fmt.Sprintf("0.0.0.0:%v", cfg.Port),
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server crashed with error : ", err)
		}
	}()

	// Start encoder event listener
	go eventhandler.ConsumeEvents(cfg, routerClients.AmqpExchangerStatus, &routerDAOs.VideosDAO)

	// Wait for SIGINT.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig

	routerClients.ServiceDiscovery.Stop()

	// Graceful shutdown for api server
	ctxServer, cancelServer := context.WithTimeout(context.Background(), GORILLA_MUX_SHUTDOWN_TIMEOUT)
	defer cancelServer()

	if err = srv.Shutdown(ctxServer); err != nil {
		// Error from closing listeners, or context timeout:
		log.Info("HTTP server Shutdown: ", err)
	}

	log.Infof("Receive signal %v. Shutting down properly", sig)
	time.Sleep(GOROUTINE_FLUSH_TIMEOUT)
}

func createRouters(cfg config.Config) (*router.Clients, *router.DAOs) {
	s3Client, err := clients.NewS3Client(cfg.S3Host, cfg.S3Region, cfg.S3Bucket, cfg.S3AuthKey, cfg.S3AuthPwd)
	if err != nil {
		log.Fatal("Failed to create S3 client: ", err)
	}

	// amqpClient for new uploaded video (api->encoder)
	amqpURL := "amqp://" + cfg.RabbitmqUser + ":" + cfg.RabbitmqPwd + "@" + cfg.RabbitmqAddr + "/"
	amqpClientVideoUpload, err := clients.NewAmqpClient(amqpURL)
	if err != nil {
		log.Fatal("Failed to create RabbitMQ client: ", err)
	}

	amqpExchangerStatus, err := clients.NewAmqpExchanger(cfg.RabbitmqAddr, cfg.RabbitmqUser, cfg.RabbitmqPwd, events.VideoUpdated)
	if err != nil {
		log.Fatal("Failed to create RabbitMQ client: ", err)
	}

	// Use "?parseTime=true" to match golang time.Time with Mariadb DATETIME types
	db, err := sql.Open("mysql", cfg.MariadbUser+":"+cfg.MariadbUserPwd+"@tcp("+cfg.MariadbHost+":"+cfg.MariadbPort+")/"+cfg.MariadbName+"?parseTime=true")
	if err != nil {
		log.Fatal("Failed to open connection with database: ", err)
	}

	videosDAO, err := dao.CreateVideosDAO(context.Background(), db)
	if err != nil {
		log.Fatal("Failed to create videos DAO : ", err)
	}

	uploadsDAO, err := dao.CreateUploadsDAO(context.Background(), db)
	if err != nil {
		log.Fatal("Failed to create uploads DAO : ", err)
	}

	discoveryClient, err := clients.NewServiceDiscovery(cfg.ConsulHost)
	if err != nil {
		log.Fatal("Cannot create consul client : ", err)
	}

	routerClients := &router.Clients{
		S3Client:            s3Client,
		AmqpClient:          amqpClientVideoUpload,
		AmqpExchangerStatus: amqpExchangerStatus,
		ServiceDiscovery:    discoveryClient,
		UUIDGen:             clients.NewUuidGenerator(),
	}

	routerDAOs := &router.DAOs{
		Db:         db,
		VideosDAO:  *videosDAO,
		UploadsDAO: *uploadsDAO,
	}

	return routerClients, routerDAOs
}
