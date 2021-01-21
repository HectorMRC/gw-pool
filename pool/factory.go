package pool

import (
	"time"
)

// NewDatapool returns a brand new datapool as Gateway
func NewDatapool(DNS string, sleep time.Duration) Pool {
	return &datapool{
		DNS:   DNS,
		Sleep: sleep,
	}
}
