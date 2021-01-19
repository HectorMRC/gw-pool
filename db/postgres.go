package db

import (
	"database/sql"
	"fmt"
	"os"

	// required by postgres connections
	_ "github.com/lib/pq"
)

const (
	errPostgresDNS = "Postgres DNS must be set."
	errConnFailed  = "Failed to open a DB connection: %s"

	envPostgresKey = "DATABASE_DNS"
)

// NewPostgresConn provides a new connection to a postgres database
func NewPostgresConn() (conn Conn, err error) {
	postgresDNS, exists := os.LookupEnv(envPostgresKey)
	if !exists {
		err = fmt.Errorf(errPostgresDNS)
		return
	}

	return sql.Open("postgres", postgresDNS)
}
