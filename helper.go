package main

import (
	"log"

	"github.com/gofrs/uuid"
)

func generateUUIDString() string {
	uuID, err := uuid.NewV4()
	if err != nil {
		log.Print("could not generate UUID", err.Error())
	}
	return uuID.String()
}
