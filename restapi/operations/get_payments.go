// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetPaymentsHandlerFunc turns a function with the right signature into a get payments handler
type GetPaymentsHandlerFunc func(GetPaymentsParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetPaymentsHandlerFunc) Handle(params GetPaymentsParams) middleware.Responder {
	return fn(params)
}

// GetPaymentsHandler interface for that can handle valid get payments params
type GetPaymentsHandler interface {
	Handle(GetPaymentsParams) middleware.Responder
}

// NewGetPayments creates a new http.Handler for the get payments operation
func NewGetPayments(ctx *middleware.Context, handler GetPaymentsHandler) *GetPayments {
	return &GetPayments{Context: ctx, Handler: handler}
}

/*GetPayments swagger:route GET /payments getPayments

List payments

*/
type GetPayments struct {
	Context *middleware.Context
	Handler GetPaymentsHandler
}

func (o *GetPayments) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetPaymentsParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
