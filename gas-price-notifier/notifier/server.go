package notifier

import (
	"fmt"
	"gas-price-notifier/config"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Start - ...
func Start() {

	handle := echo.New()

	handle.Use(middleware.Logger())

	v1 := handle.Group("/v1")

	{
		v1.GET("/", func(c echo.Context) error {
			return c.JSON(http.StatusOK, struct {
				Name string `json:"name"`
			}{Name: "Anjan"})
		})
	}

	if err := handle.Start(fmt.Sprintf(":%s", config.Get("Port"))); err != nil {
		log.Fatalf("[!] Failed to start server : %s\n", err.Error())
	}

}
