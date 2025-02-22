package pezdispenser

import (
	"gopkg.in/mgo.v2"
)

type (
	//MongoCollectionGetter - Getting collections in mongo
	MongoCollectionGetter interface {
		Collection() Persistence
	}

	//MongoCollection - interface to a collection in mongo
	MongoCollection interface {
		Remove(selector interface{}) error
		Find(query interface{}) *mgo.Query
		Upsert(selector interface{}, update interface{}) (info *mgo.ChangeInfo, err error)
	}

	//MongoCollectionWrapper - interface to wrap mongo collections with additional persistence functions
	MongoCollectionWrapper struct {
		Persistence
		col MongoCollection
	}

	//Persistence - interface to a persistence store of some kind
	Persistence interface {
		Remove(selector interface{}) error
		FindOne(query interface{}, result interface{}) (err error)
		Upsert(selector interface{}, update interface{}) (err error)
	}
)
