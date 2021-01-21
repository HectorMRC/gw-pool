package pool

import (
	"context"
	"database/sql"
)

// Location is the current position for any driver
type Location interface {
	GetLatitude() int
	GetLongitude() int
	GetDriverID() int
}

// A Conn represents the connection to any database
type Conn interface {
	Close() error
	PingContext(context.Context) error
	Exec(string, ...interface{}) (sql.Result, error)
}

// A Pool represents a controller between some data to persistence and the database
type Pool interface {
	Insert(Location)
	Reset()
	Stop()
}
