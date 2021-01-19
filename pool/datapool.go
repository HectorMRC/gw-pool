package pool

import (
	"fmt"
	"sync"

	"github.com/HectorMRC/gw-pool/location"
)

// Errors provided by the datapool
const (
	ErrAlreadyInit = "The pool has already been initialized"
)

// datapool is the default implementation of the Gateway interface
type datapool struct {
	pool sync.Pool
	mu   sync.Mutex
	run  bool
}

func (dp *datapool) Init() error {
	dp.mu.Lock()
	defer dp.mu.Unlock()

	if dp.run {
		return fmt.Errorf(ErrAlreadyInit)
	}

	return nil
}

func (dp *datapool) Insert(loc location.Location) {
	if loc != nil {
		dp.pool.Put(loc)
	}
}
