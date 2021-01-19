package pool

import (
	"fmt"
	"sync"
)

// NewDatapool returns a brand new datapool as Gateway
func NewDatapool() Pool {
	newfunc := func() interface{} {
		return fmt.Errorf("EOF")
	}

	return &datapool{
		pool: sync.Pool{
			New: newfunc,
		},
	}
}
