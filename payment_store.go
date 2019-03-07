package main

import (
	"errors"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("Record not found")

type Payment struct {
	ID    uuid.UUID
	Label string
}

type PaymentStore interface {
	Total() int
	GetMany(limit, offset int) ([]*Payment, error)
	GetOne(id uuid.UUID) (*Payment, error)
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

func (store *PaymentInMemStore) GetMany(limit, offset int) ([]*Payment, error) {
	to := limit + offset
	if to > len(store.Database) {
		return []*Payment{}, nil
	}
	return store.Database[offset:to], nil
}

func (store *PaymentInMemStore) GetOne(id uuid.UUID) (*Payment, error) {
	for _, p := range store.Database {
		if p.ID.String() == id.String() {
			return p, nil
		}
	}
	return nil, ErrNotFound
}
