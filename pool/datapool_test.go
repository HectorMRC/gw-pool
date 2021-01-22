package pool_test

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/HectorMRC/gw-pool/pool"
)

var timeout = 1 * time.Second

type connMutex struct {
	open bool
}

func (conn *connMutex) Close() error {
	conn.open = false
	return nil
}

func (conn *connMutex) PingContext(context.Context) (err error) {
	if !conn.open {
		err = fmt.Errorf("Connection is not enabled")
	}

	return
}

func (conn *connMutex) Exec(string, ...interface{}) (sql.Result, error) {
	return nil, nil
}

func (conn *connMutex) asConn() (pool.Conn, error) {
	return conn, nil
}

func newConnMutex() *connMutex {
	return &connMutex{true}
}

func ExamplePool() {
	conn := newConnMutex()
	subject := pool.NewDatapool(timeout, conn.asConn)
	subject.Reset()

	for i := 0; i < 5; i++ {
		loc := pool.NewLocation(123, 123, i+1)
		subject.Insert(loc)
	}

	// output:
}
