package data

import "sync/atomic"

// ActiveConnections - Keeping number of active websocket
// connection
type ActiveConnections struct {
	Count uint64
}

// Increment - Atomically increments count of active connections by some positive integer
func (a *ActiveConnections) Increment(by uint64) {
	atomic.AddUint64(&a.Count, by)
}

// Decrement - Atomically decrements count of active connections by some positive integer
func (a *ActiveConnections) Decrement(by uint64) {
	atomic.AddUint64(&a.Count, ^uint64(by-1))
}
