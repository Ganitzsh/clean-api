package api

import (
	"reflect"

	"github.com/google/uuid"
)

type APIVersion int32

// PaymentStore defines what a PaymentStore should be able to do
type PaymentStore interface {
	// Total should return the total of payments in the data store and return 0
	// on error
	Total() int

	// GetMany will take different parameters and should return a list of Payments
	// accordingly
	GetMany(limit, offset int, filters ...*PaymentStoreFilter) ([]*Payment, error)

	// GetByID should return a single payment corresponding to the given ID
	GetByID(id uuid.UUID) (*Payment, error)

	// Save should create or update a Payment
	Save(p *Payment) error

	// Delete should remove a Payment from the data source
	Delete(id uuid.UUID) error
}

type PaymentStoreFilterType uint

const (
	PaymentStoreFilterTypeEqual = iota
	PaymentStoreFilterTypeIn
)

// PaymentStoreFilter defines a filter that can be applied to a store query
type PaymentStoreFilter struct {
	// Field is the litteral name of the field in the Payment
	Field string

	// Want is the value that is wanted to match the filter
	Want interface{}
	Type PaymentStoreFilterType
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

// Match will take a field's value as parameter and compare it to the Wanted one
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
