package main

import (
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	// mongo db name
	mdb string
	// mongo active session
	msess *mgo.Session
)

// define doc interface
type Document interface {
	Collection() string
	Index() []mgo.Index
}

// mongo initialize connection pool to db
func mongoInit(host, db string, docs ...Document) (err error) {
	// start pool of connectin to mongo db with 1 second timeout
	msess, err = mgo.DialWithTimeout(host, time.Second)
	if err != nil {
		return
	}
	// set db to work with
	mdb = db

	// init indexes set for collections
	for _, doc := range docs {
		log.Println("ensure index for ", doc)
		if err = EnsureIndex(doc); err != nil {
			return
		}
	}

	return
}

// EnsureIndex sets indexes for collections
func EnsureIndex(doc Document) error {
	var err error
	s := msess.Copy()
	defer s.Close()

	for _, v := range doc.Index() {
		err = coll(s, doc).EnsureIndex(v)
		if err != nil {
			break
		}
	}
	return err
}

// Save stores document to mongo db
func Save(doc Document) error {
	s := msess.Copy()
	defer s.Close()

	return coll(s, doc).Insert(doc)
}

// Update change document in mongo
func UpdateByPid(doc Document, pid int64) error {
	s := msess.Copy()
	defer s.Close()

	return coll(s, doc).Update(bson.M{"pid": pid}, doc)
}

// return collection of document
func coll(s *mgo.Session, doc Document) *mgo.Collection {
	return s.DB(mdb).C(doc.Collection())
}
