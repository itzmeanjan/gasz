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
