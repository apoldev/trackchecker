package models

import (
	"encoding/json"
)

type Crawler struct {
	TrackingNumber
	Status  string          `json:"status"`
	Results []CrawlerResult `json:"results"`
}

type CrawlerResult struct {
	Spider         string          `json:"spider"`
	TrackingNumber string          `json:"tracking_number"`
	ExecuteTime    float64         `json:"execute_time"`
	Err            string          `json:"error,omitempty"`
	Result         json.RawMessage `json:"result,omitempty"`
}
