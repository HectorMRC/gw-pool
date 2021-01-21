package pool

import (
	"time"
)

// NewDatapool returns a brand new datapool as Gateway
func NewDatapool(sleep time.Duration) Pool {
	return &datapool{
		sleep: sleep,
	}
}
