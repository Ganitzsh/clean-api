package api

import (
	"reflect"
	"strings"

	"github.com/google/uuid"
)

type APIVersion int32

type PaymentStore interface {
	Total() int
	GetMany(limit, offset int, filters ...*PaymentStoreFilter) ([]*Payment, *APIError)
	GetByID(id uuid.UUID) (*Payment, *APIError)
	Save(p *Payment) *APIError
	Delete(id uuid.UUID) *APIError
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
) ([]*Payment, *APIError) {
	to := limit + offset
	if to > len(store.Database) {
		return []*Payment{}, nil
	}
	subset := []*Payment{}
	if filters != nil {
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
	if offset == 0 && limit == 0 {
		return subset, nil
	}
	if offset > 0 && limit == 0 {
		return subset[offset:], nil
	}
	return subset[offset:to], nil
}

func (store *PaymentInMemStore) GetByID(id uuid.UUID) (*Payment, *APIError) {
	for _, p := range store.Database {
		if p.ID.String() == id.String() {
			return p, nil
		}
	}
	return nil, ErrNotFound
}

func (store *PaymentInMemStore) Save(d *Payment) *APIError {
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

func (store *PaymentInMemStore) Delete(id uuid.UUID) *APIError {
	for i, storedDocument := range store.Database {
		if storedDocument.ID.String() == id.String() {
			store.Database = store.Database[:i+copy(store.Database[i:], store.Database[i+1:])]
			return nil
		}
	}
	return ErrNotFound
}