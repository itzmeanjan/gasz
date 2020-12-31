package notifier

import (
	"fmt"
	"gas-price-notifier/config"
	"gas-price-notifier/data"
	"gas-price-notifier/pubsub"
	"log"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Start - ...
func Start() {

	redisClient := pubsub.Connect()
	defer redisClient.Close()

	handle := echo.New()

	handle.Use(middleware.Logger())

	v1 := handle.Group("/v1")
	upgrader := websocket.Upgrader{}

	{
		v1.GET("/subscribe", func(c echo.Context) error {

			conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
			if err != nil {
				return err
			}

			// Scheduling graceful connection closing, to be invoked when
			// getting out of this function's scope
			defer conn.Close()

			// All errors encountered in real-time communication
			// handling logic i.e. 👇 are to be stored here, so that they
			// can be returned from this function
			var _err error
			// For each client connected over websocket, this associative
			// array to be maintained, so that we can allow each client
			// subscribe tp different price feeds using single connection
			//
			// They will receive notification as soon as any such criteria gets satisfied
			var subscriptions map[string]*data.PriceSubscription

			subscriptions = make(map[string]*data.PriceSubscription)

			// Handling client request and responding accordingly
			for {

				var payload data.Payload

				// Reading JSON data from client
				if err := conn.ReadJSON(&payload); err != nil {
					_err = err
					break
				}

				// Validating client payload
				if err := payload.Validate(); err != nil {

					if err := conn.WriteJSON(&data.ClientResponse{
						Code:    0,
						Message: "Bad Subscription Request",
					}); err != nil {
						_err = err
						break
					}

					_err = err
					break
				}

				// Kept so that after control gets out of switch case 👇
				// we can check whether we faced any errors with in switch case or not
				//
				// If yes, we need to get out of this execution loop, which will result in automatic
				// closing of underlying network connection
				var facedErrorInSwitchCase bool

				switch payload.Type {

				case "subscription":

					// Client has already subscribed to this event
					_, ok := subscriptions[payload.String()]
					if ok {

						resp := data.ClientResponse{
							Code:    0,
							Message: "Already Subscribed",
						}

						if err := conn.WriteJSON(&resp); err != nil {
							facedErrorInSwitchCase = true
							_err = err
							break
						}

					}

					// Creating subscription entry for this client in associative array
					//
					// To be used in future when `unsubscription` request to be received
					subscriptions[payload.String()] = data.NewPriceSubscription(c.Request().Context(), conn, &payload, redisClient)

				case "unsubscription":

					// Client doesn't have any subscription
					// for this event, so there's no question
					// of unsubscription
					subs, ok := subscriptions[payload.String()]
					if !ok {

						resp := data.ClientResponse{
							Code:    0,
							Message: "Not Subscribed",
						}

						if err := conn.WriteJSON(&resp); err != nil {
							facedErrorInSwitchCase = true
							_err = err
						}

						break
					}

					// Cancelling subscription
					if subs != nil {
						subs.Unsubscribe(c.Request().Context())
					}

					// Removing subscription entry from associative array
					delete(subscriptions, payload.String())

				}

				// If we've faced any errors in switch case 👆
				// we're just breaking out of loop
				if facedErrorInSwitchCase {
					break
				}

			}

			return _err

		})
	}

	if err := handle.Start(fmt.Sprintf(":%s", config.Get("Port"))); err != nil {
		log.Fatalf("[!] Failed to start server : %s\n", err.Error())
	}

}
