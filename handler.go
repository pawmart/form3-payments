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
	"github.com/pawmart/form3-payments/utils"
)

// Handler responsible for payments endpoints.
type PaymentsHandler struct {
	storage *storage.Storage
}

// GetHealth handling.
func (h *PaymentsHandler) GetHealth(params o.GetHealthParams) middleware.Responder {

	err := h.storage.Ping()

	result := new(models.Health)
	result.Status = "up"

	if err != nil {
		result.Status = "down"
	}

	return o.NewGetHealthOK().WithPayload(result)
}

// GetPayments handling.
func (h *PaymentsHandler) GetPayments(params o.GetPaymentsParams) middleware.Responder {

	var err error

	result := h.storage.FindPayments(params)
	if err != nil {
		log.Print("could not list resources, db query problems", params)
		return o.NewGetPaymentsInternalServerError()
	}

	resp := new(models.PaymentDetailsListResponse)
	resp.Data = result
	resp.Links = &models.Links{Self: new(o.GetPaymentsURL).String()}

	return o.NewGetPaymentsOK().WithPayload(resp)
}

// FetchPayment handling.
func (h *PaymentsHandler) FetchPayment(params o.GetPaymentsIDParams) middleware.Responder {

	p, err := h.storage.FindPayment(params.ID.String())
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
func (h *PaymentsHandler) CreatePayment(params o.PostPaymentsParams) middleware.Responder {

	pd := params.PaymentCreationRequest
	if pd.Data == nil {
		return o.NewPostPaymentsBadRequest().WithPayload(&models.APIError{
			ErrorCode: string(http.StatusBadRequest), ErrorMessage: "bad request"})
	}

	v := int64(1)
	id := utils.GenerateUUIDString()
	t := time.Now().Unix()

	pd.Data.ID = id
	pd.Data.Version = &v
	pd.Data.CreatedOn = &t
	pd.Data.ModifiedOn = &t

	s := h.storage
	if err := s.InsertPayment(pd.Data); err != nil {
		log.Print("could not insert resource, db problems", err)
		return o.NewPostPaymentsInternalServerError()
	}

	// Now get it...
	p, err := s.FindPayment(id)
	if err != nil {
		log.Print("freshly created entity could not be fetched", id)
		return o.NewPostPaymentsInternalServerError()
	}

	selfUrl := new(o.GetPaymentsIDURL)
	selfUrl.ID = strfmt.UUID(p.ID)

	pdr := new(models.PaymentCreationResponse)
	pdr.Data = p
	pdr.Links = &models.Links{Self: selfUrl.String()}

	return o.NewPostPaymentsCreated().WithPayload(pdr)
}

// UpdatePayment handling.
func (h *PaymentsHandler) UpdatePayment(params o.PatchPaymentsParams) middleware.Responder {

	src := params.PaymentUpdateRequest.Data

	id := src.ID
	s := h.storage
	dst, err := s.FindPayment(id)
	if err != nil {
		return o.NewPatchPaymentsNotFound().WithPayload(&models.APIError{
			ErrorCode: string(http.StatusNotFound), ErrorMessage: "not found"})
	}

	t := time.Now().Unix()
	dst.ModifiedOn = &t

	if err := mergo.Merge(dst, src, mergo.WithOverride); err != nil {
		log.Print("resource not merged", src, dst, err.Error())
		return o.NewPatchPaymentsInternalServerError()
	}

	if err := s.UpdatePayment(id, dst); err != nil {
		log.Print("resource not updated", dst, err.Error())
		return o.NewPatchPaymentsInternalServerError()
	}

	return o.NewPatchPaymentsOK()
}

// DeletePayment handling.
func (h *PaymentsHandler) DeletePayment(params o.DeletePaymentsIDParams) middleware.Responder {

	if err := h.storage.RemovePayment(params.ID.String()); err != nil {
		return o.NewDeletePaymentsIDNotFound().WithPayload(&models.APIError{
			ErrorCode: string(http.StatusNotFound), ErrorMessage: "not found"})
	}

	return o.NewDeletePaymentsIDNoContent()
}
