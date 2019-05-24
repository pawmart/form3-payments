package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/imdario/mergo"
	"github.com/satori/go.uuid"
	"gopkg.in/mgo.v2/bson"
)

const collectionName = "payments"

// HealthCheck handling.
func HealthCheck(c *gin.Context) {
	db := connectToStorage(c)
	err := db.Session.Ping()
	if err != nil {
		c.JSON(http.StatusOK, &Health{Status: "down"})
		return
	}
	c.JSON(http.StatusOK, &Health{Status: "up"})
}

// GetPayments handling.
func GetPayments(c *gin.Context) {
	var err error
	var result []*Payment
	db := connectToStorage(c)
	filterMap, ok := c.GetQueryMap("filter")
	if ok {
		var mQueries []bson.M
		for k, v := range filterMap {
			k := strings.Replace(k, "_", "", -1)
			mQueries = append(mQueries, bson.M{k: v})
		}
		err = db.C(collectionName).Find(bson.D{{"$and", mQueries}}).All(&result)
	} else {
		err = db.C(collectionName).Find(nil).All(&result)
	}
	if err != nil {
		log.Print("could not list resources, db query problems")
		returnFatal(c)
		return
	}
	pdc := new(PaymentDataList)
	pdc.Data = result
	//TODO: Links
	c.JSON(http.StatusOK, pdc)
}

// GetSinglePayment handling.
func GetSinglePayment(c *gin.Context) {
	p := new(Payment)
	db := connectToStorage(c)
	err := db.C(collectionName).Find(bson.M{"id": c.Param("id")}).One(p)
	if err != nil || p.ID == "" {
		returnError(c, http.StatusNotFound, "not found")
		return
	}
	pd := new(PaymentData)
	pd.Data = p
	pd.Links = &Links{Self: "/v1/payments/" + p.ID}
	c.JSON(http.StatusOK, pd)
}

// CreatePayment handling.
func CreatePayment(c *gin.Context) {

	pd := new(PaymentData)
	if err := c.ShouldBindJSON(&pd); err != nil {
		log.Print(err.Error())
		returnFatal(c)
		return
	}
	if pd.Data == nil {
		returnError(c, http.StatusBadRequest, "wrong request input")
		return
	}

	pd.Data.ID = uuid.NewV4().String()
	pd.Data.Type = "Payment"
	pd.Data.Version = 1

	db := connectToStorage(c)
	err := db.C(collectionName).Insert(pd.Data)
	if err != nil {
		log.Print("could not insert resource, db query problems")
		returnError(c, http.StatusInternalServerError, "contact administrator - db query")
		return
	}
	c.JSON(http.StatusCreated, pd)
}

// UpdatePayment handling.
func UpdatePayment(c *gin.Context) {
	source := new(PaymentData)
	if err := c.ShouldBindJSON(&source); err != nil {
		log.Print(err.Error())
		returnFatal(c)
		return
	}
	id := source.Data.ID
	dst := new(Payment)
	db := connectToStorage(c)
	err := db.C(collectionName).Find(bson.M{"id": id}).One(dst)
	if err != nil || dst.ID == "" {
		returnError(c, http.StatusNotFound, "not found")
		return
	}
	if err := mergo.Merge(dst, source.Data, mergo.WithOverride); err != nil {
		log.Print(err.Error())
		returnFatal(c)
		return
	}
	err = db.C(collectionName).Update(bson.M{"id": id}, dst)
	if err != nil {
		log.Print("resource not updated")
		returnFatal(c)
		return
	}
	c.AbortWithStatus(http.StatusAccepted)
}

// DeletePayment handling.
func DeletePayment(c *gin.Context) {
	db := connectToStorage(c)
	err := db.C(collectionName).Remove(bson.M{"id": c.Param("id")})
	if err != nil {
		log.Print("could not remove resource")
		returnFatal(c)
		return
	}
	c.AbortWithStatus(http.StatusNoContent)
}
