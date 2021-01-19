package pool

import "github.com/HectorMRC/gw-pool/location"

// A Pool represents a controller between some data to persistence and the database
type Pool interface {
	Insert(loc location.Location)
	Init() error
}
