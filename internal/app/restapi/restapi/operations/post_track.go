// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// PostTrackHandlerFunc turns a function with the right signature into a post track handler
type PostTrackHandlerFunc func(PostTrackParams) middleware.Responder

// Handle executing the request and returning a response
func (fn PostTrackHandlerFunc) Handle(params PostTrackParams) middleware.Responder {
	return fn(params)
}

// PostTrackHandler interface for that can handle valid post track params
type PostTrackHandler interface {
	Handle(PostTrackParams) middleware.Responder
}

// NewPostTrack creates a new http.Handler for the post track operation
func NewPostTrack(ctx *middleware.Context, handler PostTrackHandler) *PostTrack {
	return &PostTrack{Context: ctx, Handler: handler}
}

/*
	PostTrack swagger:route POST /track postTrack

# Tracking

Create tracking request with some tracking numbers
*/
type PostTrack struct {
	Context *middleware.Context
	Handler PostTrackHandler
}

func (o *PostTrack) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewPostTrackParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
