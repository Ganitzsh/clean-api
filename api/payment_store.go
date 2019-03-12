package api

import (
	"reflect"
	"strings"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type APIVersion int32

type PaymentStore interface {
	Total() int
	GetMany(limit, offset int, filters ...*PaymentStoreFilter) ([]*Payment, error)
	GetByID(id uuid.UUID) (*Payment, error)
	Save(p *Payment) error
	Delete(id uuid.UUID) error
}

type PaymentStoreFilterType uint

const (
	PaymentStoreFilterTypeEqual = iota
	PaymentStoreFilterTypeIn
)

type PaymentStoreFilter struct {
	Field string
	Want  interface{}
	Type  PaymentStoreFilterType
}

func NewPaymentStoreFilter() *PaymentStoreFilter {
	return &PaymentStoreFilter{}
}

func (sf *PaymentStoreFilter) SetWant(value interface{}) *PaymentStoreFilter {
	sf.Want = value
	return sf
}

func (sf *PaymentStoreFilter) SetType(value PaymentStoreFilterType) *PaymentStoreFilter {
	sf.Type = value
	return sf
}

func (f PaymentStoreFilter) Match(has interface{}) (bool, error) {
	if f.Want == nil && has == nil {
		return false, ErrSomethingWentWrong(nil)
	}
	switch f.Type {
	case PaymentStoreFilterTypeEqual:
		return reflect.DeepEqual(has, f.Want), nil
	case PaymentStoreFilterTypeIn:
		s := reflect.ValueOf(f.Want)
		if s.Kind() != reflect.Slice {
			return false, ErrUnsupportedFilterValue
		}
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(s.Index(i).Interface(), has) {
				return true, nil
			}
		}
		return false, nil
	default:
		return false, ErrUnknownFilterType
	}
}

type PaymentInMemStore struct {
	Database []*Payment
}

func NewPaymentInMemStore() *PaymentInMemStore {
	return &PaymentInMemStore{
		Database: []*Payment{},
	}
}

func (store *PaymentInMemStore) Total() int {
	return len(store.Database)
}

func PaymentStoreFilterIsScheme(typ string) *PaymentStoreFilter {
	return &PaymentStoreFilter{
		Field: "Scheme",
		Want:  typ,
		Type:  PaymentStoreFilterTypeEqual,
	}
}

func (store *PaymentInMemStore) GetMany(
	limit, offset int,
	filters ...*PaymentStoreFilter,
) ([]*Payment, error) {
	to := limit + offset
	if to > len(store.Database) {
		return []*Payment{}, nil
	}
	subset := []*Payment{}
	if filters != nil && len(filters) > 0 {
		for _, filter := range filters {
			for _, d := range store.Database {
				e := reflect.ValueOf(d).Elem()
				for i := 0; i < e.NumField(); i++ {
					varName := e.Type().Field(i).Name
					varValue := e.Field(i).Interface()
					if strings.ToLower(varName) == strings.ToLower(filter.Field) {
						isMatching, _ := filter.Match(varValue)
						if isMatching {
							subset = append(subset, d)
						}
					}
				}
			}
		}
	} else {
		subset = append(subset, store.Database...)
	}
	if offset > len(subset) {
		return []*Payment{}, nil
	}
	if offset == 0 && limit == 0 {
		return subset, nil
	}
	if offset > 0 && limit == 0 {
		return subset[offset:], nil
	}
	return subset[offset:to], nil
}

func (store *PaymentInMemStore) GetByID(id uuid.UUID) (*Payment, error) {
	for _, p := range store.Database {
		if p.ID.String() == id.String() {
			return p, nil
		}
	}
	return nil, ErrNotFound
}

func (store *PaymentInMemStore) Save(d *Payment) error {
	if d == nil {
		return ErrSomethingWentWrong(ErrNilValue)
	}
	if d.ID.String() != "" {
		for i, storedDocument := range store.Database {
			if d.ID.String() == storedDocument.ID.String() {
				store.Database[i] = d
				return nil
			}
		}
	}
	store.Database = append(store.Database, d)
	return nil
}

func (store *PaymentInMemStore) Delete(id uuid.UUID) error {
	for i, storedDocument := range store.Database {
		if storedDocument.ID.String() == id.String() {
			store.Database = store.Database[:i+copy(store.Database[i:], store.Database[i+1:])]
			return nil
		}
	}
	return nil
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
