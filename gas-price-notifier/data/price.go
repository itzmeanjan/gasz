package data

import "sync"

// GasPrice - Safely read/ update latest gas price
// to be used when answering HTTP GET queries
type GasPrice struct {
	Latest *PubSubPayload
	Lock   *sync.RWMutex
}

// Get - Safely read latest gas price
func (g *GasPrice) Get() PubSubPayload {

	g.Lock.RLock()
	defer g.Lock.RUnlock()

	return *g.Latest

}

// Put - Safely updated gas price to latest value received
func (g *GasPrice) Put(latest *PubSubPayload) {

	g.Lock.Lock()
	defer g.Lock.Unlock()

	g.Latest = latest

}
