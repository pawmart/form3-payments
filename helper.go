package main

import (
	"log"

	"github.com/gofrs/uuid"
	"github.com/pawmart/form3-payments/storage"
	mgo "gopkg.in/mgo.v2"
)

func connectToStorage() *mgo.Database {
	db, err := storage.New().Db()
	if err != nil {
		log.Print("connection to db problems", err.Error())
	}
	return db
}

func generateUUIDString() string {
	uuID, err := uuid.NewV4()
	if err != nil {
		log.Print("could not generate UUID", err.Error())
	}
	return uuID.String()
}
