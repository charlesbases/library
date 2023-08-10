package webserver

// Request .
type Request struct {
}

// Response .
type Response struct {
	ID   string      `json:"id"`
	Code Code        `json:"code,omitempty"`
	Data interface{} `json:"data,omitempty"`
}
