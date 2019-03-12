package api

import (
	"reflect"

	"github.com/google/uuid"
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
