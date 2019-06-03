package storage

import (
	"log"
	"strings"

	"github.com/pawmart/form3-payments/config"
	"gopkg.in/mgo.v2"
)

// Storage struct.
type Storage struct {
	Config *config.DbConfig
	db     *mgo.Database
}

func (s *Storage) getDB() (db *mgo.Database) {
	if s.db != nil {
		return s.db
	}

	conf := s.Config
	var session *mgo.Session

	session, err := s.dial(s.Config)
	if err != nil {
		log.Fatalln(err)
	}

	if conf.User != "" && conf.Password != "" {
		err := session.Login(&mgo.Credential{Username: conf.User, Password: conf.Password, Source: conf.Auth})
		if err != nil {
			log.Fatalln(err)
		}
	}

	session.SetSafe(&mgo.Safe{})

	s.db = session.DB(conf.Database)
	return s.db
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
