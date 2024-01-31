package models

type TrackingNumber struct {
	RequestID string `json:"id,omitempty"`
	UUID      string `json:"uuid,omitempty"`
	Code      string `json:"code"`
}
