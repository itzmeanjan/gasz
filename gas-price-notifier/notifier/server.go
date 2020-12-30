package notifier

import (
	"fmt"
	"gas-price-notifier/config"
	"gas-price-notifier/data"
	"gas-price-notifier/pubsub"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Start - ...
func Start() {

	redisClient := pubsub.Connect()
	defer redisClient.Close()

	handle := echo.New()

	handle.Use(middleware.Logger())
	handle.Validator = &data.CustomValidator{Validator: validator.New()}

	v1 := handle.Group("/v1")

	{
		v1.GET("/subscribe", func(c echo.Context) error {

			var payload data.Payload

			if err := c.Bind(&payload); err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}

			if err := c.Validate(&payload); err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}

			return c.JSON(http.StatusOK, payload)

		})
	}

	if err := handle.Start(fmt.Sprintf(":%s", config.Get("Port"))); err != nil {
		log.Fatalf("[!] Failed to start server : %s\n", err.Error())
	}

}
