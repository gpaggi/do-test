package models

import "encoding/json"

// Response data structure
type Response struct {
	HttpCode int
	Output   json.RawMessage
}
