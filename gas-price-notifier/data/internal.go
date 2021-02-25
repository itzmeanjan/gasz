package data

import (
	"context"
	"encoding/json"
	"fmt"
	"gas-price-notifier/config"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

// SubscribeToPriceFeed - This function can be invoked as a seperate go routine, which will subscribe
// to latest gas price event & updated state safely, which can be read from API layer & responded back
// to client(s) in concurrent safe fashion
func SubscribeToPriceFeed(ctx context.Context, redisClient *redis.Client, latestGasPrice *GasPrice) {

	pubsub := redisClient.Subscribe(ctx, config.Get("RedisPubSubChannel"))

OUTER:
	for {

	INNER:
		select {

		// If caller context is cancelled, we can unsubscribe
		case <-ctx.Done():

			if err := pubsub.Unsubscribe(ctx, config.Get("RedisPubSubChannel")); err != nil {

				log.Printf("[❌] Failed to unsubscribe : %s\n", err.Error())

			}
			break OUTER

		default:

			msg, err := pubsub.ReceiveTimeout(ctx, time.Second)
			if err != nil {
				continue
			}

			switch m := msg.(type) {

			case *redis.Subscription:

				if m.Kind == "subscribe" {

					log.Printf(fmt.Sprintf("[✅] Subscribed to %s\n", config.Get("RedisPubSubChannel")))
					break INNER

				}

				if m.Kind == "unsubscribe" {

					log.Printf(fmt.Sprintf("[❌] Unsubscribed from %s\n", config.Get("RedisPubSubChannel")))
					break OUTER

				}

			case *redis.Message:

				var pubsubPayload PubSubPayload
				_msg := []byte(m.Payload)

				if err := json.Unmarshal(_msg, &pubsubPayload); err != nil {

					log.Printf("[❌] Failed to decode received data from pubsub channel : %s\n", err.Error())
					break OUTER
				}

				// Update gas price to latest
				latestGasPrice.Put(&pubsubPayload)

			}

		}

	}

}
