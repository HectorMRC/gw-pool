package pool

import (
	"time"
)

// NewDatapool returns a brand new datapool as Gateway
func NewDatapool(open ConnFunc, sleep time.Duration) Pool {
	return &datapool{
		Open:  open,
		Sleep: sleep,
	}
}
