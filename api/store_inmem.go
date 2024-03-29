package api

import (
	"reflect"
	"strings"

	"github.com/google/uuid"
)

// This is an implementation of PaymentStore with temporary in memory storage

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
) (*PaginatedList, error) {
	to := limit + offset
	if to > len(store.Database) {
		return &PaginatedList{Results: []*Payment{}}, nil
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
	total := len(subset)
	if offset > len(subset) {
		return &PaginatedList{Results: []*Payment{}}, nil
	}
	if offset == 0 && limit == 0 {
		return &PaginatedList{
			Total:    total,
			SubTotal: len(subset),
			Results:  subset,
		}, nil
	}
	if offset > 0 && limit == 0 {
		ret := subset[offset:]
		return &PaginatedList{
			Total:    total,
			SubTotal: len(ret),
			Results:  ret,
		}, nil
	}
	ret := subset[offset:to]
	return &PaginatedList{
		Total:    total,
		SubTotal: len(ret),
		Results:  ret,
	}, nil
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
				d.UpdatedAt = Now()
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
