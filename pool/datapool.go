package pool

import (
	"container/list"
	"context"
	"database/sql"
	"log"
	"sync"
	"time"

	// required by postgres connections
	_ "github.com/lib/pq"
)

const sqlStatement = `
	INSERT INTO locations (latitude, longitude, driver_id)
	VALUES ($1, $2, $3)`

// datapool is the default implementation of the Gateway interface
type datapool struct {
	DNS    string
	Sleep  time.Duration
	stack  list.List
	mu     sync.Mutex
	cancel context.CancelFunc
	cond   *sync.Cond
}

func (dp *datapool) newPostgresConn(ctx context.Context) (conn Conn, err error) {
	if conn, err = sql.Open("postgres", dp.DNS); err != nil {
		return
	}

	err = conn.PingContext(ctx)
	return
}

func (dp *datapool) waitForConnectivity(ctx context.Context) (conn Conn, err error) {
	ticker := time.NewTicker(dp.Sleep)
	defer ticker.Stop()

	for conn, err = dp.newPostgresConn(ctx); err != nil; {
		select {
		case <-ticker.C:
			if conn != nil {
				conn.Close()
			}

			log.Printf("Failed to connect to database: %v", err.Error())
			conn, err = dp.newPostgresConn(ctx)

		case <-ctx.Done():
			err = ctx.Err()
			log.Printf("Context has been canceled: %v", err.Error())
			return
		}
	}

	return
}

func (dp *datapool) execQueryForeach(conn Conn) (err error) {
	defer conn.Close()
	ilen := dp.stack.Len()

	// while no error happen, and there still elements to persist
	for err == nil && dp.stack.Len() > 0 {
		elem := dp.stack.Front()
		loc, ok := elem.Value.(Location)
		if ok {
			_, err = conn.Exec(sqlStatement, loc.GetLatitude(), loc.GetLongitude(), loc.GetDriverID())
		}

		// if the elemen has no value of type location, or the location inside has been properly persisted
		// the element must be removed from the pool
		if !ok || err == nil {
			dp.stack.Remove(elem)
		}
	}

	log.Printf("%v of %v locations have been persisted", ilen-dp.stack.Len(), ilen)
	return
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
		dp.execQueryForeach(conn)
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

func (dp *datapool) Insert(loc Location) {
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
