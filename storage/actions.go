package storage

import (
	"github.com/pawmart/form3-payments/models"
	"github.com/pawmart/form3-payments/restapi/operations"
	"gopkg.in/mgo.v2/bson"
)

const collectionName = "payments"

// InsertPayment handling.
func (s *Storage) InsertPayment(payment *models.Payment) error {
	return s.getDB().C(collectionName).Insert(payment)
}

// FindPayment handling.
func (s *Storage) FindPayment(id string) (payment *models.Payment, err error) {
	payment = new(models.Payment)
	err = s.getDB().C(collectionName).Find(bson.M{"id": id}).One(payment)
	if err != nil {
		return nil, err
	}
	return payment, nil
}

// UpdatePayment handling.
func (s *Storage) UpdatePayment(id string, p *models.Payment) error {
	return s.getDB().C(collectionName).Update(bson.M{"id": id}, p)
}

// RemovePayment handling.
func (s *Storage) RemovePayment(id string) error {
	return s.getDB().C(collectionName).Remove(bson.M{"id": id})
}

// FindPayments handling.
func (s *Storage) FindPayments(params operations.GetPaymentsParams) []*models.Payment {
	var q interface{}
	shouldFilter := false
	for _, k := range params.FilterOrganisationID {
		k.String()
		shouldFilter = true
	}
	if shouldFilter {
		var mQueries []bson.M
		for _, v := range params.FilterOrganisationID {
			mQueries = append(mQueries, bson.M{"organisationid": v.String()})
		}
		q = bson.D{{"$and", mQueries}}
	}

	var result []*models.Payment
	s.getDB().C(collectionName).Find(q).All(&result)
	return result
}

// Ping database.
func (s *Storage) Ping() error {
	return s.getDB().Session.Ping()
}

// Drop database.
func (s *Storage) Drop() {
	s.getDB().DropDatabase()
}
