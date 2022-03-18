package clients

import (
	"database/sql"
)

var _ IMariadbClient = mariadbClientDummy{}

type mariadbClientDummy struct {
	closeConn       func() error
	checkConnection func() error
	getDb           func() *sql.DB
}

func NewMariadbClientDummy(closeConn func() error, checkConnection func() error, getDb func() *sql.DB) IMariadbClient {
	return mariadbClientDummy{closeConn, checkConnection, getDb}
}

func (m mariadbClientDummy) CloseConn() error {
	if m.closeConn != nil {
		return m.closeConn()
	}
	return nil
}

func (m mariadbClientDummy) CheckConnection() error {
	if m.checkConnection != nil {
		return m.checkConnection()
	}
	return nil
}

func (m mariadbClientDummy) GetDb() *sql.DB {
	if m.getDb != nil {
		return m.getDb()
	}
	return nil
}
