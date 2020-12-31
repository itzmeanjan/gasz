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

			var _err error

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

					if err := conn.WriteJSON(&data.ErrorResponse{
						Message: "Bad Subscription Request",
					}); err != nil {
						_err = err
						break
					}

					_err = err
					break
				}

				// Writing subscription confirmation message
				if err := conn.WriteJSON(&data.ErrorResponse{
					Message: "Subscribed",
				}); err != nil {
					_err = err
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
