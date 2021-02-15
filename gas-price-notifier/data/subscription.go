package data

import (
	"context"
	"encoding/json"
	"fmt"
	"gas-price-notifier/config"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

// PriceSubscription - Manages whole life cycle of price feed subscription
// for each client
//
// Functions defined on this struct, are supposed to be invoked for subscribing to and unsubscribing from
// Redis pubsub topic, where price feed data is being published
type PriceSubscription struct {
	Client    *websocket.Conn
	Redis     *redis.Client
	PubSub    *redis.PubSub
	Topics    map[string]*Payload
	TopicLock *sync.RWMutex
	ConnLock  *sync.Mutex
}

// Subscribe - Subscribing to Redis pubsub channel
// so that any time new price feed is posted in channel
// listener will get notified & take proper measurements
// if conditions satisfy
func (ps *PriceSubscription) Subscribe(req *Payload) {

	// -- Safely reading from associative array
	// shared among multiple go routines
	ps.TopicLock.RLock()

	_, ok := ps.Topics[req.String()]

	ps.TopicLock.RUnlock()
	// -- Safely read whether client is alreadY subscribed
	// to topic or not

	// Client is already subscribed to topic
	if ok {

		resp := ClientResponse{
			Code:    0,
			Message: "Already Subscribed",
		}

		// -- Critical section code, starts
		ps.ConnLock.Lock()

		if err := ps.Client.WriteJSON(&resp); err != nil {
			log.Printf("[!] Failed to communicate with client : %s\n", err.Error())
		}

		ps.ConnLock.Unlock()
		// -- Critical section code, ends

		return

	}

	// -- Attempting to safely write to shared
	// associative array
	ps.TopicLock.Lock()

	ps.Topics[req.String()] = req

	ps.TopicLock.Unlock()
	// -- Done with safe writing to shared memory space

	// Attempting to let client know
	// subscription has been confirmed, for requested
	// topic
	resp := ClientResponse{
		Code:    1,
		Message: fmt.Sprintf("Subscribed to `%s`", req.String()),
	}

	// -- Critical section code, starts
	ps.ConnLock.Lock()

	if err := ps.Client.WriteJSON(&resp); err != nil {
		log.Printf("[!] Failed to communicate with client : %s\n", err.Error())
	}

	ps.ConnLock.Unlock()
	// -- Critical section code, ends

}

// Unsubscribe - Cancelling price feed subscription for specific user
// and letting client know about it
func (ps *PriceSubscription) Unsubscribe(req *Payload) {

	ps.TopicLock.RLock()

	_, ok := ps.Topics[req.String()]

	ps.TopicLock.RUnlock()

	if !ok {

		resp := ClientResponse{
			Code:    0,
			Message: "Not Subscribed",
		}

		// -- Critical section code, starts
		ps.ConnLock.Lock()

		if err := ps.Client.WriteJSON(&resp); err != nil {
			log.Printf("[!] Failed to communicate with client : %s\n", err.Error())
		}

		ps.ConnLock.Unlock()
		// -- Critical section code, ends

		return

	}

	// -- Safely removing topic entry from map, i.e. client
	// not subscribed to this topic anymore
	ps.TopicLock.Lock()

	delete(ps.Topics, req.String())

	ps.TopicLock.Unlock()
	// -- Done with removing entry from shared map

	// Attempting to let client know
	// unsubscription has been confirmed, for requested
	// topic
	resp := ClientResponse{
		Code:    1,
		Message: fmt.Sprintf("Unsubscribed from `%s`", req.String()),
	}

	// -- Critical section code, starts
	ps.ConnLock.Lock()

	if err := ps.Client.WriteJSON(&resp); err != nil {
		log.Printf("[!] Failed to communicate with client : %s\n", err.Error())
	}

	ps.ConnLock.Unlock()
	// -- Critical section code, ends

}

// Listen - Subscribing to Redis pubsub and waiting for message
// to be published, as soon as it's published it's being sent to
// client application, connected via websocket connection
//
// Of course criteria evaluation is performed before sending
// notification to client
//
// Listening criteria to be specified by client application
// during subscription phase
func (ps *PriceSubscription) Listen(ctx context.Context) {

	for {

		msg, err := ps.PubSub.ReceiveTimeout(ctx, time.Second)
		if err != nil {
			continue
		}

		var stopListening bool

		switch m := msg.(type) {

		case *redis.Subscription:

			if m.Kind == "subscribe" {

				log.Printf(fmt.Sprintf("[*] Subscribed to %s\n", config.Get("RedisPubSubChannel")))
				break

			}

			if m.Kind == "unsubscribe" {

				stopListening = true
				log.Printf(fmt.Sprintf("[*] Unsubscribed from %s\n", config.Get("RedisPubSubChannel")))

				break

			}

		case *redis.Message:

			var pubsubPayload PubSubPayload
			_msg := []byte(m.Payload)

			if err := json.Unmarshal(_msg, &pubsubPayload); err != nil {

				stopListening = true
				log.Printf("[!] Failed to decode received data from pubsub channel : %s\n", err.Error())

				break
			}

			// If not satisfying criteria, then we're not attempting to deliver
			//
			// Otherwise, delivery attempt to be made
			if !ps.isEligibleForDelivery(&pubsubPayload) {
				break
			}

			// -- Critical section of code, starts
			ps.ConnLock.Lock()

			// Attempting to deliver price feed data, which they've subscribed to
			if err := ps.Client.WriteJSON(ps.GetClientResponse(&pubsubPayload)); err != nil {

				stopListening = true
				log.Printf("[!] Failed to communicate with client : %s\n", err.Error())

			}

			// -- Critical section of code, ends
			ps.ConnLock.Unlock()

		}

		// If something went wrong in execution flow with in `switch-case`
		// block, we're going to get out of listener loop
		if stopListening {
			break
		}

	}

}

// Checking whether price data feed received is eligible for delivery, by comparing
// with evaluation criteria provided by user, when subscribing to price feed
//
// We'll simply check whether price value of certain category i.e. {fast, fastest, safeLow, average},
// to which client has subscribed to, is {<, >, <=, >=, ==} to gas price they have provided us with
//
// If yes, we're going to deliver this piece of data to client
func (ps *PriceSubscription) isEligibleForDelivery(payload *PubSubPayload) (bool, *Payload) {

	// -- Closure starting here
	//
	// Given obtained gas price of certain category i.e. to which
	// client has subscribed to and criteria specified by them,
	// we'll check whether it's satisfying requirement or not
	//
	// This closure is written so that it becomes easier to use
	// for all price subscription categories i.e. {fast, fastest, safeLow, average}
	checkThreshold := func(price float64, operator string, threshold float64) bool {

		var status bool

		switch operator {

		case "<":
			status = price < threshold
		case ">":
			status = price > threshold
		case "<=":
			status = price <= threshold
		case ">=":
			status = price >= threshold
		case "==":
			status = price == threshold

		}

		return status

	}
	// -- Closure ending here

	// To be returned back to caller
	//
	// If `status` is true, `request` will also be
	// non-nil
	var status bool
	var request *Payload

	for _, v := range ps.Topics {

		switch v.Field {

		case "fast":
			status = checkThreshold(payload.Fast, v.Operator, v.Threshold)
		case "fastest":
			status = checkThreshold(payload.Fastest, v.Operator, v.Threshold)
		case "safeLow":
			status = checkThreshold(payload.SafeLow, v.Operator, v.Threshold)
		case "average":
			status = checkThreshold(payload.Average, v.Operator, v.Threshold)

		}

		if status {

			request = v
			break

		}

	}

	return status, request

}

// GetClientResponse - Returns gas price response to be sent to client,
// when gas price reaches certain value, which satisfies client set criteria
func (ps *PriceSubscription) GetClientResponse(payload *PubSubPayload) *GasPriceFeed {

	var gasPrice GasPriceFeed

	switch ps.Request.Field {

	case "fast":
		gasPrice.Price = payload.Fast
	case "fastest":
		gasPrice.Price = payload.Fastest
	case "safeLow":
		gasPrice.Price = payload.SafeLow
	case "average":
		gasPrice.Price = payload.Average

	}

	gasPrice.TxType = ps.Request.Field
	gasPrice.Topic = ps.Request.String()

	return &gasPrice

}

// NewPriceSubscription - Creating new price data subscription for client
// connected over websocket
//
// Whether client will receive notification that depends on whether received price value
// satisfies criteria set by client
func NewPriceSubscription(ctx context.Context, client *websocket.Conn, request *Payload, redisClient *redis.Client, topicLock *sync.RWMutex, connLock *sync.Mutex) *PriceSubscription {

	ps := PriceSubscription{
		Client:    client,
		Redis:     redisClient,
		PubSub:    redisClient.Subscribe(ctx, config.Get("RedisPubSubChannel")),
		Topics:    make(map[string]*Payload),
		TopicLock: topicLock,
		ConnLock:  connLock,
	}

	// Running listener i.e. subscriber in different execution thread
	go ps.Listen(ctx)

	return &ps

}
