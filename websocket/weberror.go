package websocket

import "encoding/json"

// WebError .
type WebError struct {
	ErrCode int    `json:"err_code,omitempty"`
	ErrMsg  string `json:"err_msg,omitempty"`
}

// Error .
func (err *WebError) Error() string {
	data, _ := json.Marshal(err)
	return string(data)
}
