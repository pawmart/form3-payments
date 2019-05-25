// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/pawmart/form3-payments/models"
)

// GetPaymentsIDOKCode is the HTTP code returned for type GetPaymentsIDOK
const GetPaymentsIDOKCode int = 200

/*GetPaymentsIDOK Payment details

swagger:response getPaymentsIdOK
*/
type GetPaymentsIDOK struct {

	/*
	  In: Body
	*/
	Payload *models.PaymentDetailsResponse `json:"body,omitempty"`
}

// NewGetPaymentsIDOK creates GetPaymentsIDOK with default headers values
func NewGetPaymentsIDOK() *GetPaymentsIDOK {

	return &GetPaymentsIDOK{}
}

// WithPayload adds the payload to the get payments Id o k response
func (o *GetPaymentsIDOK) WithPayload(payload *models.PaymentDetailsResponse) *GetPaymentsIDOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get payments Id o k response
func (o *GetPaymentsIDOK) SetPayload(payload *models.PaymentDetailsResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetPaymentsIDOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetPaymentsIDNotFoundCode is the HTTP code returned for type GetPaymentsIDNotFound
const GetPaymentsIDNotFoundCode int = 404

/*GetPaymentsIDNotFound Resource not found

swagger:response getPaymentsIdNotFound
*/
type GetPaymentsIDNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.APIError `json:"body,omitempty"`
}

// NewGetPaymentsIDNotFound creates GetPaymentsIDNotFound with default headers values
func NewGetPaymentsIDNotFound() *GetPaymentsIDNotFound {

	return &GetPaymentsIDNotFound{}
}

// WithPayload adds the payload to the get payments Id not found response
func (o *GetPaymentsIDNotFound) WithPayload(payload *models.APIError) *GetPaymentsIDNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get payments Id not found response
func (o *GetPaymentsIDNotFound) SetPayload(payload *models.APIError) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetPaymentsIDNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}