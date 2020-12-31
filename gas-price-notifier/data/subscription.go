package data

import (
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

// PriceSubscription - ...
type PriceSubscription struct {
	Client  *websocket.Conn
	Request *Payload
	Redis   *redis.Client
	PubSub  *redis.PubSub
}
