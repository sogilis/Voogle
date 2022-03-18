package clients

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

type IMariadbClient interface {
	CloseConn() error
	CheckConnection() error
	GetDb() *sql.DB
}

var _ IMariadbClient = &mariadbClient{}

type mariadbClient struct {
	Db    *sql.DB
	dbUrl string
}

func NewMariadbClient(user, userPwd, addr, name string) (IMariadbClient, error) {
	db := &mariadbClient{
		Db:    nil,
		dbUrl: user + ":" + userPwd + "@tcp(" + addr + ")/" + name,
	}

	database, err := sql.Open("mysql", db.dbUrl)
	if err != nil {
		return nil, err
	}

	db.Db = database
	return db, nil
}

func (db *mariadbClient) CloseConn() error {
	if db.Db != nil {
		if err := db.Db.Close(); err != nil {
			log.Error("Error while closing database : ", err)
			return err
		}
		log.Info("Connection to database closed")
		return nil
	}
	log.Info("Connection to database not yet open")
	return nil
}

func (db *mariadbClient) CheckConnection() error {
	// Check the server version
	if db.Db == nil {
		return errors.New("connexion to database not yet open")
	}

	var version string
	if err := db.Db.QueryRow("SELECT VERSION();").Scan(&version); err != nil {
		return err
	}

	log.Info("Connected to: ", version)
	return nil
}

func (db *mariadbClient) GetDb() *sql.DB {
	return db.Db
}
