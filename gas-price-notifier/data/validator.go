package data

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// CustomValidator - Custom payload validator
type CustomValidator struct {
	Validator *validator.Validate
}

// Validate - Custom payload validator's validate method to be
// called when new request is received
func (cv *CustomValidator) Validate(i interface{}) error {
	return echo.NewHTTPError(http.StatusBadRequest, cv.Validator.Struct(i).Error())
}
