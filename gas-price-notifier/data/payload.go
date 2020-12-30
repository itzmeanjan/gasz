package data

// Payload - Payload received from client via websocket connection
type Payload struct {
	Field     string  `json:"field" validate:"required"`
	Threshold float64 `json:"threshold" validate:"required"`
	Operator  string  `json:"operator" validate:"required"`
}
