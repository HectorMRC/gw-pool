package pool

import (
	"context"
	"database/sql"

	"github.com/HectorMRC/gw-pool/location"
)

// A Conn represents the connection to any database
type Conn interface {
	Close() error
	PingContext(context.Context) error
	Exec(string, ...interface{}) (sql.Result, error)
}

// A Pool represents a controller between some data to persistence and the database
type Pool interface {
	Insert(loc location.Location)
	Reset()
	Stop()
}
