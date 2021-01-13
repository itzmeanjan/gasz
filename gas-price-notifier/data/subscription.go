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
	Client  *websocket.Conn
	Request *Payload
	Redis   *redis.Client
	PubSub  *redis.PubSub
	Lock    *sync.Mutex
}

// Subscribe - Subscribing to Redis pubsub channel
// so that any time new price feed is posted in channel
// listener will get notified & take proper measurements
// if conditions satisfy
func (ps *PriceSubscription) Subscribe(ctx context.Context) {
	ps.PubSub = ps.Redis.Subscribe(ctx, config.Get("RedisPubSubChannel"))
}

// Listen - Subscribing to Redis pubsub and waiting for message
// to be published, as soon as it's published it's being sent to
// client application, connected via websocket connection
//
//
func (ps *PriceSubscription) Listen(ctx context.Context) {

	// Scheduling unsubscription call here, to be invoked when
	// returning from this function
	defer ps.Unsubscribe(ctx)

	for {

		if ps.Request.Type != "subscription" {
			break
		}

		msg, err := ps.PubSub.ReceiveTimeout(ctx, time.Second)
		if err != nil {
			continue
		}

		var facedErrorInSwitchCase bool

		switch m := msg.(type) {

		case *redis.Subscription:

			resp := ClientResponse{
				Code:    1,
				Message: fmt.Sprintf("Subscribed to `%s`", ps.Request),
			}

			// -- Critical section of code, starts
			ps.Lock.Lock()

			if err := ps.Client.WriteJSON(&resp); err != nil {
				facedErrorInSwitchCase = true
				log.Printf("[!] Failed to communicate with client : %s\n", err.Error())

			}

			// -- Critical section of code, ends
			ps.Lock.Unlock()

		case *redis.Message:

			var pubsubPayload PubSubPayload
			_msg := []byte(m.Payload)

			if err := json.Unmarshal(_msg, &pubsubPayload); err != nil {
				facedErrorInSwitchCase = true
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
			ps.Lock.Lock()

			// Attempting to deliver price feed data, which they've subscribed to
			if err := ps.Client.WriteJSON(ps.GetClientResponse(&pubsubPayload)); err != nil {
				facedErrorInSwitchCase = true
				log.Printf("[!] Failed to communicate with client : %s\n", err.Error())
			}

			// -- Critical section of code, ends
			ps.Lock.Unlock()

		}

		// Checking whether we've encountered any error with in switch case
		//
		// If yes, we can break out of this loop
		if facedErrorInSwitchCase {
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
func (ps *PriceSubscription) isEligibleForDelivery(payload *PubSubPayload) bool {

	// Given obtained gas price of certain category i.e. to which
	// client has subscribed to and criteria specified by them,
	// we'll check whether it's satisfying requirement or not
	//
	// This closure is written so that it becomes easier to use
	// for all price subscription categories i.e. {fast, fastest, safeLow, average}
	checkThreshold := func(price float64) bool {

		status := true

		switch ps.Request.Operator {
		case "<":
			status = price < ps.Request.Threshold
		case ">":
			status = price > ps.Request.Threshold
		case "<=":
			status = price <= ps.Request.Threshold
		case ">=":
			status = price >= ps.Request.Threshold
		case "==":
			status = price == ps.Request.Threshold
		}

		return status

	}

	status := true

	switch ps.Request.Field {
	case "fast":
		status = checkThreshold(payload.Fast)
	case "fastest":
		status = checkThreshold(payload.Fastest)
	case "safeLow":
		status = checkThreshold(payload.SafeLow)
	case "average":
		status = checkThreshold(payload.Average)
	}

	return status

}

// GetClientResponse - Returns gas price response to be sent to client,
// when gas price reaches certain value, which satisfies client set criteria
func (ps *PriceSubscription) GetClientResponse(payload *PubSubPayload) *GasPriceFeed {

	var gasPrice GasPriceFeed

	switch t := ps.Request.Field; t {

	case "fast":
		gasPrice.Price = payload.Fast
		gasPrice.TxType = t
	case "fastest":
		gasPrice.Price = payload.Fastest
		gasPrice.TxType = t
	case "safeLow":
		gasPrice.Price = payload.SafeLow
		gasPrice.TxType = t
	case "average":
		gasPrice.Price = payload.Average
		gasPrice.TxType = t

	}

	return &gasPrice

}

// Unsubscribe - Cancelling price feed subscription for specific user
// and letting client know about it
func (ps *PriceSubscription) Unsubscribe(ctx context.Context) {

	if err := ps.PubSub.Unsubscribe(ctx, config.Get("RedisPubSubChannel")); err != nil {
		log.Printf("[!] Failed to unsubscribe from pubsub topic : %s\n", err.Error())
		return
	}

	resp := ClientResponse{
		Code:    1,
		Message: fmt.Sprintf("Unsubscribed from `%s`", ps.Request),
	}

	ps.Lock.Lock()
	defer ps.Lock.Unlock()

	if err := ps.Client.WriteJSON(&resp); err != nil {
		log.Printf("[!] Failed to communicate with client : %s\n", err.Error())
	}

}

// NewPriceSubscription - Creating new price data subscription for client
// connected over websocket
//
// Whether client will receive notification that depends on whether received price value
// satisfies criteria set by client
func NewPriceSubscription(ctx context.Context, client *websocket.Conn, request *Payload, redisClient *redis.Client, lock *sync.Mutex) *PriceSubscription {

	ps := PriceSubscription{
		Client:  client,
		Request: request,
		Redis:   redisClient,
		Lock:    lock,
	}

	// Subscription object to be stored in 👆 struct
	// after calling this function
	ps.Subscribe(ctx)
	// Running listener i.e. subscriber in different execution thread
	go ps.Listen(ctx)

	return &ps

}
