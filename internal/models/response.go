package models

type Response struct {
	Result  interface{} `json:"result,omitempty"`
	Deleted interface{} `json:"deleted,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}
