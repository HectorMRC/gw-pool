package pool

import (
	"time"
)

// NewDatapool returns a brand new datapool as Gateway
func NewDatapool(sleep time.Duration, open ConnFunc) Pool {
	return &datapool{
		Open:  open,
		Sleep: sleep,
	}
}

// NewLocation builds a brand new Location
func NewLocation(latitude, longitude, driverID int) Location {
	return &Coordinates{
		Latitude:  latitude,
		Longitude: longitude,
		DriverID:  driverID,
	}
}
