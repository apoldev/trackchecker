package models

type TrackingNumber struct {
	ID   int64  `json:"id,omitempty"`
	UUID string `json:"uuid,omitempty"`
	Code string `json:"code"`
}
