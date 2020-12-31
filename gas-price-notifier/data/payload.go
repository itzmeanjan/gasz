package data

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Payload - Payload received from client via websocket connection
type Payload struct {
	Field     string  `json:"field"`
	Threshold float64 `json:"threshold"`
	Operator  string  `json:"operator"`
}

// Validate - Validates payload received from client
// and in case of success returns `nil`
//
// Error handling required in caller side
func (p *Payload) Validate() error {

	return validation.ValidateStruct(p,
		validation.Field(&p.Field, validation.Required, validation.In("fast", "fastest", "safeLow", "average")),
		validation.Field(&p.Threshold, validation.Required, validation.Min(1.0)),
		validation.Field(&p.Operator, validation.Required, validation.In("<", ">", "<=", ">=", "==")))

}

// ErrorResponse - When some error is encountered while processing
// client request, client to be notified with data of this form
type ErrorResponse struct {
	Message string `json:"message"`
}
