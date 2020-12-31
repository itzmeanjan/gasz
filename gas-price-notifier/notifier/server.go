package notifier

import (
	"fmt"
	"gas-price-notifier/config"
	"gas-price-notifier/data"
	"gas-price-notifier/pubsub"
	"log"
	"net/http"

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

	{
		v1.POST("/subscribe", func(c echo.Context) error {

			var payload data.Payload

			if err := c.Bind(&payload); err != nil {

				return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
					Message: "Bad Payload",
				})

			}

			if err := payload.Validate(); err != nil {

				return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
					Message: "Bad Payload",
				})

			}

			return c.JSON(http.StatusOK, payload)

		})
	}

	if err := handle.Start(fmt.Sprintf(":%s", config.Get("Port"))); err != nil {
		log.Fatalf("[!] Failed to start server : %s\n", err.Error())
	}

}
