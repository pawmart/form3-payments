package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/imdario/mergo"
	"github.com/pawmart/form3-payments/models"
	o "github.com/pawmart/form3-payments/restapi/operations"
	"github.com/pawmart/form3-payments/storage"
)

// GetHealth handling.
func GetHealth(params o.GetHealthParams) middleware.Responder {

	err := storage.New().Ping()

	result := new(models.Health)
	result.Status = "up"

	if err != nil {
		result.Status = "down"
	}

	return o.NewGetHealthOK().WithPayload(result)
}

// GetPayments handling.
func GetPayments(params o.GetPaymentsParams) middleware.Responder {

	var err error

	s := storage.New()
	result := s.FindPayments(params)
	if err != nil {
		log.Print("could not list resources, db query problems", params)
	}

	resp := new(models.PaymentDetailsListResponse)
	resp.Data = result
	resp.Links = &models.Links{Self: new(o.GetPaymentsURL).String()}

	return o.NewGetPaymentsOK().WithPayload(resp)
}

// FetchPayment handling.
func FetchPayment(params o.GetPaymentsIDParams) middleware.Responder {

	s := storage.New()
	p, err := s.FindPayment(params.ID.String())
	if err != nil {
		return o.NewGetPaymentsIDNotFound().WithPayload(&models.APIError{
			ErrorCode:    string(http.StatusNotFound),
			ErrorMessage: "not found",
		})
	}

	selfUrl := new(o.GetPaymentsIDURL)
	selfUrl.ID = strfmt.UUID(p.ID)

	pd := new(models.PaymentDetailsResponse)
	pd.Data = p
	pd.Links = &models.Links{Self: "/v1/payments/" + p.ID}
	pd.Links = &models.Links{Self: selfUrl.String()}

	return o.NewGetPaymentsIDOK().WithPayload(pd)
}

// CreatePayment handling.
func CreatePayment(params o.PostPaymentsParams) middleware.Responder {

	pd := params.PaymentCreationRequest
	if pd.Data == nil {
		return o.NewPostPaymentsBadRequest().WithPayload(&models.APIError{
			ErrorCode: string(http.StatusBadRequest), ErrorMessage: "bad request"})
	}

	v := int64(1)
	id := generateUUIDString()
	t := time.Now().Unix()

	pd.Data.ID = id
	pd.Data.Version = &v
	pd.Data.CreatedOn = &t
	pd.Data.ModifiedOn = &t

	s := storage.New()
	if err := s.InsertPayment(pd.Data); err != nil {
		log.Print("could not insert resource, db query problems", err)
	}

	// Now get it...
	p, err := s.FindPayment(id)
	if err != nil {
		return o.NewGetPaymentsIDNotFound().WithPayload(&models.APIError{
			ErrorCode:    string(http.StatusNotFound),
			ErrorMessage: "not found",
		})
	}

	selfUrl := new(o.GetPaymentsIDURL)
	selfUrl.ID = strfmt.UUID(p.ID)

	pdr := new(models.PaymentCreationResponse)
	pdr.Data = p
	pdr.Links = &models.Links{Self: selfUrl.String()}

	return o.NewPostPaymentsCreated().WithPayload(pdr)
}

// UpdatePayment handling.
func UpdatePayment(params o.PatchPaymentsParams) middleware.Responder {

	src := params.PaymentUpdateRequest.Data

	id := src.ID
	s := storage.New()
	dst, err := s.FindPayment(id)
	if err != nil {
		return o.NewPatchPaymentsNotFound().WithPayload(&models.APIError{
			ErrorCode: string(http.StatusNotFound), ErrorMessage: "not found"})
	}

	t := time.Now().Unix()
	dst.ModifiedOn = &t

	if err := mergo.Merge(dst, src, mergo.WithOverride); err != nil {
		log.Print("resource not merged", src, dst, err.Error())
	}

	if err := s.UpdatePayment(id, dst); err != nil {
		log.Print("resource not updated", dst, err.Error())
	}

	return o.NewPatchPaymentsOK()
}

// DeletePayment handling.
func DeletePayment(params o.DeletePaymentsIDParams) middleware.Responder {

	s := storage.New()
	if err := s.RemovePayment(params.ID.String()); err != nil {
		return o.NewDeletePaymentsIDNotFound().WithPayload(&models.APIError{
			ErrorCode: string(http.StatusNotFound), ErrorMessage: "not found"})
	}

	return o.NewDeletePaymentsIDNoContent()
}
