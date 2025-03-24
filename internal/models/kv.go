package models

type KeyValue struct {
	Key   string                 `json:"key"`
	Value map[string]interface{} `json:"value"`
}
