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

func OpenConn(user, userPwd, addr, name string) (database, error) {
	dbUrl := user + ":" + userPwd + "@tcp(" + addr + ")/" + name
	log.Info("Open connection to database")
	db, err := sql.Open("mysql", dbUrl)
	if err != nil {
		log.Info("Cannot access database : ", err)
		return database{}, err
	}

	return database{db}, nil
}

func (db database) CloseConn() error {
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
