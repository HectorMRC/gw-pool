package pool

import (
	"time"
)

// NewDatapool returns a brand new datapool as Gateway
func NewDatapool(sleep time.Duration, open ConnFunc) Pool {
	dp := &datapool{
		Open:  open,
		Sleep: sleep,
	}

	dp.cond.L = dp.mu.RLocker()
	return dp
}

// NewLocation builds a brand new Location
func NewLocation(latitude, longitude, driverID int) Location {
	return &Coordinates{
		Latitude:  latitude,
		Longitude: longitude,
		DriverID:  driverID,
	}
}
