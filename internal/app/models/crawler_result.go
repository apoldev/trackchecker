package models

import "encoding/json"

type CrawlerResult struct {
	ExecuteTime float64         `json:"execute_time,omitempty"`
	Err         string          `json:"error,omitempty"`
	Result      json.RawMessage `json:"result,omitempty"`
}
