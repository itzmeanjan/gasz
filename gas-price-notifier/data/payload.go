package data

// Payload - Payload received from client via websocket connection
type Payload struct {
	Field     string  `json:"field"`
	Threshold float64 `json:"threshold"`
	Operator  string  `json:"operator"`
}

// ErrorResponse - When some error is encountered while processing
// client request, client to be notified with data of this form
type ErrorResponse struct {
	Message string `json:"message"`
}
