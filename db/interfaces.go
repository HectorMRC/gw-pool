package db

import "context"

// A Conn represents the connection between the backend and any database
type Conn interface {
	Close() error
	PingContext(context.Context) error
}
