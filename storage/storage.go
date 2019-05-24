package storage

import (
	"log"
	"strings"

	"github.com/pawmart/form3-payments/config"
	"gopkg.in/mgo.v2"
)

// Storage struct.
type Storage struct {
	config *config.DbConfig
	db     *mgo.Database
}

// New storage instance creation.
func New() (s *Storage) {
	s = new(Storage)
	s.init()
	return s
}

// Db dial and return.
func (s *Storage) Db() (db *mgo.Database, err error) {

	conf := s.config
	var session *mgo.Session

	session, err = s.dial(s.config)
	if err != nil {
		return nil, err
	}

	if conf.User != "" && conf.Password != "" {
		err := session.Login(&mgo.Credential{Username: conf.User, Password: conf.Password, Source: conf.Auth})
		if err != nil {
			return nil, err
		}
	}

	session.SetSafe(&mgo.Safe{})

	s.db = session.DB(conf.Database)
	return s.db, nil
}

func (s *Storage) init() {
	s.config = new(config.Config).LoadConfiguration().Db
	return
}

func (s *Storage) dial(conf *config.DbConfig) (*mgo.Session, error) {

	// TODO: Enforce secure connection for production.
	host := s.formatHost(conf)
	sess, err := mgo.Dial(host)
	if err != nil {
		log.Print("Connection to mongo error: " + err.Error())
		return nil, err
	}
	cloned := sess.Clone()
	defer sess.Close()

	return cloned, err
}

func (s *Storage) formatHost(conf *config.DbConfig) string {
	h := strings.TrimSpace(conf.Host)
	h = strings.TrimPrefix(h, "mongodb://")
	return h
}
