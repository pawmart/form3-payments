package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/imdario/mergo"
	"github.com/pawmart/form3-payments/models"
	"github.com/pawmart/form3-payments/restapi/operations"
	"gopkg.in/mgo.v2/bson"
)

const collectionName = "payments"

// GetHealth handling.
func GetHealth(params operations.GetHealthParams) middleware.Responder {

	db := connectToStorage()
	err := db.Session.Ping()

	result := new(models.Health)
	result.Status = "up"

	if err != nil {
		result.Status = "down"
	}

	return operations.NewGetHealthOK().WithPayload(result)
}

// GetPayments handling.
func GetPayments(params operations.GetPaymentsParams) middleware.Responder {

	var err error
	var result []*models.Payment
	db := connectToStorage()

	shouldFilter := false
	for _, k := range params.FilterOrganisationID {
		log.Print("dealing with organisation " + k.String())
		shouldFilter = true
	}
	if shouldFilter {
		var mQueries []bson.M
		for _, v := range params.FilterOrganisationID {
			mQueries = append(mQueries, bson.M{"organisationid": v.String()})
		}
		err = db.C(collectionName).Find(bson.D{{"$and", mQueries}}).All(&result)
	} else {
		err = db.C(collectionName).Find(nil).All(&result)
	}
	if err != nil {
		log.Print("could not list resources, db query problems", params)
	}

	resp := new(models.PaymentDetailsListResponse)
	resp.Data = result

	return operations.NewGetPaymentsOK().WithPayload(resp)
}

// FetchPayment handling.
func FetchPayment(params operations.GetPaymentsIDParams) middleware.Responder {

	p := new(models.Payment)
	db := connectToStorage()
	err := db.C(collectionName).Find(bson.M{"id": params.ID}).One(p)
	if err != nil || p.ID == "" {
		return operations.NewGetPaymentsIDNotFound().WithPayload(&models.APIError{
			ErrorCode:    string(http.StatusNotFound),
			ErrorMessage: "not found",
		})
	}

	pd := new(models.PaymentDetailsResponse)
	pd.Data = p
	pd.Links = &models.Links{Self: "/v1/payments/" + p.ID}

	return operations.NewGetPaymentsIDOK().WithPayload(pd)
}

// CreatePayment handling.
func CreatePayment(params operations.PostPaymentsParams) middleware.Responder {

	pd := params.PaymentCreationRequest
	if pd.Data == nil {
		return operations.NewPostPaymentsBadRequest().WithPayload(&models.APIError{
			ErrorCode: string(http.StatusBadRequest), ErrorMessage: "bad request"})
	}

	v := int64(1)
	id := generateUUIDString()
	t := time.Now().Unix()

	pd.Data.ID = id
	pd.Data.Version = &v
	pd.Data.CreatedOn = &t
	pd.Data.ModifiedOn = &t

	db := connectToStorage()
	if err := db.C(collectionName).Insert(pd.Data); err != nil {
		log.Print("could not insert resource, db query problems")
	}

	// Now get it...
	p := new(models.Payment)

	err := db.C(collectionName).Find(bson.M{"id": id}).One(p)
	if err != nil || p.ID == "" {

		if err != nil {
			log.Print("Failed to find resource ", err.Error())
		}
		return operations.NewGetPaymentsIDNotFound().WithPayload(&models.APIError{
			ErrorCode:    string(http.StatusNotFound),
			ErrorMessage: "not found",
		})
	}

	pdr := new(models.PaymentCreationResponse)
	pdr.Data = p
	pdr.Links = &models.Links{Self: "/v1/payments/" + id}

	return operations.NewPostPaymentsCreated().WithPayload(pdr)
}

// UpdatePayment handling.
func UpdatePayment(params operations.PatchPaymentsParams) middleware.Responder {

	log.Print("initial UpdatePayment")

	src := params.PaymentUpdateRequest.Data

	log.Print("SRC ", src)
	log.Print("REF to update ", src.Attributes.Reference)

	dst := new(models.Payment)
	id := src.ID
	db := connectToStorage()
	if err := db.C(collectionName).Find(bson.M{"id": id}).One(dst); err != nil || dst.ID == "" {
		return operations.NewPatchPaymentsNotFound().WithPayload(&models.APIError{
			ErrorCode: string(http.StatusNotFound), ErrorMessage: "not found"})
	}

	log.Print("REF fetched", dst.Attributes.Reference)
	log.Print("MOD fetched", dst.ModifiedOn)
	log.Print("CRE fetched", dst.CreatedOn)

	t := time.Now().Unix()
	dst.ModifiedOn = &t

	if err := mergo.Merge(dst, src, mergo.WithOverride); err != nil {
		log.Print(err.Error())
		log.Print("resource not merged", src, dst, err.Error())
	}

	log.Print("REF after merge", dst.Attributes.Reference)

	if err := db.C(collectionName).Update(bson.M{"id": id}, dst); err != nil {
		log.Print("resource not updated", dst, err.Error())
	}

	return operations.NewPatchPaymentsOK()
}

// DeletePayment handling.
func DeletePayment(params operations.DeletePaymentsIDParams) middleware.Responder {

	db := connectToStorage()
	if err := db.C(collectionName).Remove(bson.M{"id": params.ID}); err != nil {
		return operations.NewDeletePaymentsIDNotFound().WithPayload(&models.APIError{
			ErrorCode: string(http.StatusNotFound), ErrorMessage: "not found"})
	}

	return operations.NewDeletePaymentsIDNoContent()
}
