// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// TrackingNumber tracking number
//
// swagger:model trackingNumber
type TrackingNumber struct {

	// code
	Code string `json:"code,omitempty"`

	// uuid
	UUID string `json:"uuid,omitempty"`
}

// Validate validates this tracking number
func (m *TrackingNumber) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this tracking number based on context it is used
func (m *TrackingNumber) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *TrackingNumber) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *TrackingNumber) UnmarshalBinary(b []byte) error {
	var res TrackingNumber
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
