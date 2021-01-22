package pool

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"
)

var timeout = time.Second

type connMutex struct {
	open bool
	locs []Location
}

func (conn *connMutex) Open() {
	conn.open = true
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

func (conn *connMutex) Exec(query string, vars ...interface{}) (_ sql.Result, err error) {
	if !conn.open {
		err = fmt.Errorf("Connection refused")
		return
	}

	if len(vars) != 3 {
		err = fmt.Errorf("Not enought values")
		return
	}

	loc := NewLocation(vars[0].(int), vars[1].(int), vars[2].(int))
	conn.locs = append(conn.locs, loc)
	return nil, nil
}

func (conn *connMutex) asConn() (Conn, error) {
	return conn, nil
}

func newConnMutex(open bool) *connMutex {
	return &connMutex{
		open: open,
	}
}

func newSubject(sleep time.Duration, open ConnFunc) *datapool {
	dp := &datapool{
		Open:  open,
		Sleep: sleep,
	}

	dp.cond.L = dp.mu.RLocker()
	return dp
}

func TestReset(t *testing.T) {
	conn := newConnMutex(false)
	subject := newSubject(timeout, conn.asConn)

	// setting up a dummy context
	ctx, cancel := context.WithCancel(context.TODO())
	subject.cancel = cancel

	loc := NewLocation(0, 0, 0)
	subject.stack.Init().PushBack(loc)

	subject.Reset()
	defer subject.Stop()

	if subject.cancel == nil {
		t.Errorf("Items must be set")
	}

	// reset must call the current CancelFunc, so the new context must be canceled
	if ctx.Err() == nil {
		t.Errorf("Reset function must kill the current context, if any")
	}

	// reset must clear the current stack
	if subject.stack.Len() != 0 {
		t.Errorf("Stop function must clear the stack")
	}
}

func TestStop(t *testing.T) {
	conn := newConnMutex(false)
	subject := newSubject(timeout, conn.asConn)

	subject.Reset()
	defer subject.Stop()

	// replacing the current CancelFunc with the one of a new context
	ctx, cancel := context.WithCancel(context.TODO())
	defer subject.cancel()
	subject.cancel = cancel

	subject.Stop()
	// stop must call the current CancelFunc, so the new context must be canceled
	if ctx.Err() == nil {
		t.Errorf("Stop function must kill the current context")
	}
}

func TestInsert(t *testing.T) {
	conn := newConnMutex(false)
	subject := newSubject(timeout, conn.asConn)

	want := 5
	// Inserting location into the list
	for i := 0; i < want; i++ {
		loc := NewLocation(123, 123, i)
		subject.Insert(loc)
	}

	if got := subject.stack.Len(); got != want {
		t.Errorf("Got %v items into the stack, want %v", got, want)
	}

	want = 0
	// Checking the insertion order
	for it := subject.stack.Front(); it != nil; it = it.Next() {
		if got := it.Value.(Location).GetDriverID(); got != want {
			t.Errorf("Got driver_id %v, want %v", got, want)
		}

		want++
	}
}

func TestInsert_conn_restored(t *testing.T) {
	conn := newConnMutex(false) // no connection
	subject := newSubject(timeout, conn.asConn)

	subject.Reset()
	defer subject.Stop()

	want := 5
	// Inserting location into the list
	for i := 0; i < want; i++ {
		loc := NewLocation(123, 123, i+1)
		subject.Insert(loc)
	}

	conn.Open()
	// some time to make sure the subject got enough time to store all locations
	time.Sleep(100 * time.Microsecond)

	// once connection is restored, all locations must into the stack must be
	// persisted into the database and cleared from the stack
	if got := len(conn.locs); got != want {
		t.Errorf("Got %v items into the database, want %v", got, want)
	}

	want = 0
	if got := subject.stack.Len(); got != want {
		t.Errorf("Got %v items into the stack, want %v", got, want)
	}
}

func TestInsert_conn_open(t *testing.T) {
	conn := newConnMutex(true)
	subject := newSubject(timeout, conn.asConn)

	subject.Reset()
	defer subject.Stop()

	want := 5
	// Inserting location into the list
	for i := 0; i < want; i++ {
		loc := NewLocation(123, 123, i+1)
		subject.Insert(loc)
	}

	// some time to make sure the subject got enough time to store all locations
	time.Sleep(100 * time.Microsecond)

	// all locations must be persisted into the database
	if got := len(conn.locs); got != want {
		t.Errorf("Got %v items into the database, want %v", got, want)
	}

	want = 0
	if got := subject.stack.Len(); got != want {
		t.Errorf("Got %v items into the stack, want %v", got, want)
	}

}
