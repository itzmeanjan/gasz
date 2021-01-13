package data

import "sync"

// ActiveConnections - Keeping number of active websocket
// connection
type ActiveConnections struct {
	Count uint64
}

// SafeActiveConnections - Connection Info which can be safely updated
// from multiple worker in concurrent fashion
type SafeActiveConnections struct {
	Connections *ActiveConnections
	Lock        *sync.RWMutex
}

// Increment - Increments count of active connections by some positive integer
func (s *SafeActiveConnections) Increment(by uint64) {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	s.Connections.Count += by
}

// Decrement - Decrements count of active connections by some positive integer
func (s *SafeActiveConnections) Decrement(by uint64) {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	s.Connections.Count -= by
}
