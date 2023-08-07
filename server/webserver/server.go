package webserver

import (
	"github.com/charlesbases/library/sonyflake"
)

// Request .
type Request struct {
}

// Response .
type Response struct {
	ID   sonyflake.ID `json:"id"`
	Code Code         `json:"code,omitempty"`
	Data interface{}  `json:"data,omitempty"`
}
