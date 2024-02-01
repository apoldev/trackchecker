// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/apoldev/trackchecker/internal/app/restapi/models"
)

// PostTrackCreatedCode is the HTTP code returned for type PostTrackCreated
const PostTrackCreatedCode int = 201

/*
PostTrackCreated OK

swagger:response postTrackCreated
*/
type PostTrackCreated struct {

	/*
	  In: Body
	*/
	Payload *models.RequestResult `json:"body,omitempty"`
}

// NewPostTrackCreated creates PostTrackCreated with default headers values
func NewPostTrackCreated() *PostTrackCreated {

	return &PostTrackCreated{}
}

// WithPayload adds the payload to the post track created response
func (o *PostTrackCreated) WithPayload(payload *models.RequestResult) *PostTrackCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post track created response
func (o *PostTrackCreated) SetPayload(payload *models.RequestResult) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostTrackCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*
PostTrackDefault error

swagger:response postTrackDefault
*/
type PostTrackDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostTrackDefault creates PostTrackDefault with default headers values
func NewPostTrackDefault(code int) *PostTrackDefault {
	if code <= 0 {
		code = 500
	}

	return &PostTrackDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the post track default response
func (o *PostTrackDefault) WithStatusCode(code int) *PostTrackDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the post track default response
func (o *PostTrackDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the post track default response
func (o *PostTrackDefault) WithPayload(payload *models.Error) *PostTrackDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post track default response
func (o *PostTrackDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostTrackDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}