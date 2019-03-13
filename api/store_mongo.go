package api

import (
	"reflect"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// MongoQuery interfaces the *mgo.Query type
type MongoQuery interface {
	All(result interface{}) error
	One(result interface{}) error
	Skip(n int) MongoQuery
	Limit(n int) MongoQuery
}

// MongoCollection interfaces the *mgo.Collection type
type MongoCollection interface {
	FindId(id interface{}) MongoQuery
	RemoveId(id interface{}) error
	Insert(docs ...interface{}) error
	Find(query interface{}) MongoQuery
	Count() (int, error)
}

type MgoWrapQuery struct {
	*mgo.Query
}

func (q *MgoWrapQuery) Skip(n int) MongoQuery {
	return &MgoWrapQuery{}
}

func (q *MgoWrapQuery) Limit(n int) MongoQuery {
	return &MgoWrapQuery{}
}

type PaymentMongoStore struct {
	MongoCollection
}

func NewPaymentMongoStore(c MongoCollection) *PaymentMongoStore {
	return &PaymentMongoStore{c}
}

func (store *PaymentMongoStore) Total() int {
	n, _ := store.Count()
	return n
}

func (store *PaymentMongoStore) GetMany(
	limit, offset int,
	filters ...*PaymentStoreFilter,
) ([]*Payment, error) {
	ret := []*Payment{}
	query := bson.M{}
	if filters != nil {
		t := reflect.TypeOf(Payment{})
		for _, f := range filters {
			for i := 0; i < t.NumField(); i++ {
				if t.Field(i).Name == f.Field {
					tagValue := t.Field(i).Tag.Get("bson")
					if tagValue == "" {
						tagValue = t.Field(i).Tag.Get("json")
					} else if tagValue == "" {
						tagValue = f.Field
					}
					switch f.Type {
					case PaymentStoreFilterTypeEqual:
						query[tagValue] = f.Want
					case PaymentStoreFilterTypeIn:
						if t.Kind() != reflect.Slice {
							logrus.Warn("This filter is expecting a slice")
						} else {
							query[tagValue] = bson.M{
								"$in": f.Want,
							}
						}
					}
				}
			}
		}
	}
	q := store.Find(query).Skip(offset)
	if limit > 0 {
		q = q.Limit(limit)
	}
	if err := q.All(&ret); err != nil {
		return nil, ErrSomethingWentWrong(err)
	}
	return ret, nil
}

func (store *PaymentMongoStore) GetByID(id uuid.UUID) (*Payment, error) {
	ret := Payment{}
	if err := store.FindId(id).One(&ret); err != nil {
		if err != mgo.ErrNotFound {
			return nil, ErrSomethingWentWrong(err)
		} else {
			return nil, ErrNotFound
		}
	}
	return &ret, nil
}

func (store *PaymentMongoStore) Save(p *Payment) error {
	if len(p.ID) == 0 {
		p.ID = uuid.New()
	}
	if err := store.Insert(p); err != nil {
		return ErrSomethingWentWrong(err)
	}
	return nil
}

func (store *PaymentMongoStore) Delete(id uuid.UUID) error {
	if err := store.RemoveId(id); err != nil {
		if err != mgo.ErrNotFound {
			return ErrSomethingWentWrong(err)
		}
	}
	return nil
}