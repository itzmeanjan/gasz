package data

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Payload - Payload received from client via websocket connection
type Payload struct {
	Type      string  `json:"type"`
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
		validation.Field(&p.Type, validation.Required, validation.In("subscription", "unsubscription")),
		validation.Field(&p.Field, validation.Required, validation.In("fast", "fastest", "safeLow", "average")),
		validation.Field(&p.Threshold, validation.Required, validation.Min(1.0)),
		validation.Field(&p.Operator, validation.Required, validation.In("<", ">", "<=", ">=", "==")))

}

// String - String representatino of subscription request, to be
// used as unique identifier in HashMap
func (p Payload) String() string {
	return fmt.Sprintf("%s : %s %f", p.Field, p.Operator, p.Threshold)
}

// ClientResponse - Client to be notified with data of this form
type ClientResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// PubSubPayload - Data received from Redis pubsub subscription
// will look like this
type PubSubPayload struct {
	Fast    float64 `json:"fast"`
	Fastest float64 `json:"fastest"`
	SafeLow float64 `json:"safeLow"`
	Average float64 `json:"average"`
}

// GasPriceFeed - Client will receive gas price feed in 👇
// form i.e. when gas price of subscribed category reaches
// threshold specified, it'll sent to client
//
// Data of following form will help browser client in showing notification easily
type GasPriceFeed struct {
	TxType string  `json:"txType"`
	Price  float64 `json:"price"`
	Topic  string  `json:"topic"`
}
