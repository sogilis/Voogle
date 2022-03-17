package db

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

type database struct {
	Db *sql.DB
}

func OpenConn(user string, userPwd string, addr string, name string) (database, error) {
	//"user:password@tcp(127.0.0.1:3306)/voogle-database"
	dbUrl := user + ":" + userPwd + "@tcp(" + addr + ")/" + name
	log.Info("Open connection to database")
	db, err := sql.Open("mysql", dbUrl)
	if err != nil {
		log.Info("Cannot access database : ", err)
		return database{}, err
	}

	return database{db}, nil
}

func (db database) CloseConn() {
	if db.Db != nil {
		if err := db.Db.Close(); err != nil {
			log.Info("Error while closing database : ", err)
			return
		}
		log.Info("Connection to database closed")
		return
	}
	log.Info("Connection to database not yet open")
}

func (db database) CheckConnection() error {
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
