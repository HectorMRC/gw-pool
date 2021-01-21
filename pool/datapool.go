package pool

import (
	"container/list"
	"context"
	"log"
	"sync"
	"time"

	"github.com/HectorMRC/gw-pool/db"
	"github.com/HectorMRC/gw-pool/location"
)

// datapool is the default implementation of the Gateway interface
type datapool struct {
	stack  list.List
	mu     sync.Mutex
	cancel context.CancelFunc
	sleep  time.Duration
	cond   *sync.Cond
}

func (dp *datapool) newConnection(ctx context.Context) (conn db.Conn, err error) {
	if conn, err = db.NewPostgresConn(); err != nil {
		return
	}

	err = conn.PingContext(ctx)
	return
}

func (dp *datapool) waitForConnectivity(ctx context.Context) (conn db.Conn, err error) {
	ticker := time.NewTicker(dp.sleep)
	defer ticker.Stop()

	for conn, err = dp.newConnection(ctx); err != nil; {
		select {
		case <-ticker.C:
			conn.Close()
			log.Printf("Failed to connect to database: %v", err.Error())
			conn, err = dp.newConnection(ctx)
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	return
}

func (dp *datapool) executeQuery(conn db.Conn) error {
	defer conn.Close()
	return nil
}

func (dp *datapool) scheduler(ctx context.Context) {
	for ctx.Err() == nil {
		// each iteration waits for the stack to be non empty
		dp.cond.L.Lock()
		if dp.stack.Len() == 0 {
			dp.cond.Wait()
		}

		// waiting for database connectivity
		conn, err := dp.waitForConnectivity(ctx)
		if err != nil {
			log.Printf("%v", err.Error())
			break
		}

		// persisting all current data
		dp.executeQuery(conn)
		dp.cond.L.Unlock()
	}
}

func (dp *datapool) kill() {
	if dp.cancel != nil {
		dp.cancel()
		dp.cancel = nil
	}
}

func (dp *datapool) Reset() {
	dp.mu.Lock()
	defer dp.mu.Unlock()

	// if the pool was already initialized, it requires to kill the current goroutines
	// before resetting it
	dp.kill()

	ctx := context.Background()
	ctx, dp.cancel = context.WithCancel(ctx)
	dp.cond = sync.NewCond(&sync.Mutex{})
	dp.stack.Init()

	go dp.scheduler(ctx)
}

func (dp *datapool) Stop() {
	dp.mu.Lock()
	defer dp.mu.Unlock()

	dp.kill()
}

func (dp *datapool) Insert(loc location.Location) {
	dp.mu.Lock()
	defer dp.mu.Unlock()

	dp.cond.L.Lock()
	if loc != nil {
		dp.stack.PushBack(loc)
	}

	// wakes up the scheduler goroutine
	dp.cond.Broadcast()
	dp.cond.L.Unlock()
}
