package mock

import (
	"github.com/ganitzsh/f3-te/api"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/google/uuid"
)

type PaymentQuery struct {
	*api.MgoWrapQuery
	id    uuid.UUID
	query bson.M
	fail  bool
	limit int
	skip  int
	store *api.PaymentInMemStore
}

func NewPaymentQuery() *PaymentQuery {
	return &PaymentQuery{}
}

func (q *PaymentQuery) All(results interface{}) error {
	ret := results.(*[]*api.Payment)
	filters := []*api.PaymentStoreFilter{}
	if val, ok := q.query["scheme"].(string); ok {
		filters = append(filters, api.PaymentStoreFilterIsScheme(val))
	}
	tmp, err := q.store.GetMany(q.limit, q.skip, filters...)
	if err != nil {
		return err
	}
	for _, d := range tmp {
		*ret = append(*ret, d)
	}
	return nil
}

func (q *PaymentQuery) One(result interface{}) error {
	tmp, err := q.store.GetByID(q.id)
	if err != nil {
		if err == api.ErrNotFound {
			return mgo.ErrNotFound
		}
		return err
	}
	*result.(*api.Payment) = *tmp
	return nil
}

func (q *PaymentQuery) Limit(n int) api.MongoQuery {
	q.limit = n
	return q
}

func (q *PaymentQuery) Skip(n int) api.MongoQuery {
	q.skip = n
	return q
}

type PaymentCollection struct {
	*mgo.Collection
	Data *api.PaymentInMemStore
}

func NewCollection(s *api.PaymentInMemStore) *PaymentCollection {
	return &PaymentCollection{
		Data: s,
	}
}

func (c *PaymentCollection) FindId(id interface{}) api.MongoQuery {
	return &PaymentQuery{
		store: c.Data,
		id:    id.(uuid.UUID),
	}
}

func (c *PaymentCollection) Insert(docs ...interface{}) error {
	return nil
}

func (c *PaymentCollection) Find(query interface{}) api.MongoQuery {
	return &PaymentQuery{
		store: c.Data,
		query: query.(bson.M),
	}
}

func (c *PaymentCollection) RemoveId(id interface{}) error {
	return c.Data.Delete(id.(uuid.UUID))
}

func (c *PaymentCollection) Count() (int, error) {
	return c.Data.Total(), nil
}
