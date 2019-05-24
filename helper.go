package main

import (
	"log"
	"net/http"

	"github.com/pawmart/form3-payments/storage"
	mgo "gopkg.in/mgo.v2"

	"github.com/gin-gonic/gin"
)

func returnError(c *gin.Context, code int, message string) {
	error := APIError{
		Code:    string(code),
		Message: message}

	c.AbortWithStatusJSON(code, &error)
}

func returnFatal(c *gin.Context) {
	//c.JSON(http.StatusInternalServerError, "")
	c.AbortWithStatus(http.StatusInternalServerError)
}

func connectToStorage(c *gin.Context) *mgo.Database {
	db, err := storage.New().Db()
	if err != nil {
		log.Print("connection to db problems", err.Error())
		returnError(c, http.StatusInternalServerError, "")
		c.Done()
	}
	return db
}
