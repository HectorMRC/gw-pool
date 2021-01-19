package db

// A Conn represents the connection between the backend and any database
type Conn interface {
	Close() error
}
