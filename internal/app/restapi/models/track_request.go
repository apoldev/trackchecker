// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// TrackRequest track request
//
// swagger:model trackRequest
type TrackRequest struct {

	// tracking numbers
	// Required: true
	TrackingNumbers []string `json:"tracking_numbers"`
}

// Validate validates this track request
func (m *TrackRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateTrackingNumbers(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *TrackRequest) validateTrackingNumbers(formats strfmt.Registry) error {

	if err := validate.Required("tracking_numbers", "body", m.TrackingNumbers); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this track request based on context it is used
func (m *TrackRequest) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *TrackRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *TrackRequest) UnmarshalBinary(b []byte) error {
	var res TrackRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
