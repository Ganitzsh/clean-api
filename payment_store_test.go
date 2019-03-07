package main

import (
	"testing"

	"github.com/google/uuid"
	"github.com/icrowley/fake"
	"github.com/stretchr/testify/assert"
)

type mockDB struct {
	Total    int
	ID1      uuid.UUID
	ID2      uuid.UUID
	ID3      uuid.UUID
	Payment1 *Payment
	Payment2 *Payment
	Payment3 *Payment
	Store    PaymentStore
}

func newTestDBInMem() *mockDB {
	store := NewPaymentInMemStore()
	store.Database = append(store.Database, []*Payment{
		&Payment{ID: uuid.New(), Label: fake.Sentence()},
		&Payment{ID: uuid.New(), Label: fake.Sentence()},
		&Payment{ID: uuid.New(), Label: fake.Sentence()},
	}...)
	return &mockDB{
		Total:    3,
		ID1:      store.Database[0].ID,
		Payment1: store.Database[0],
		ID2:      store.Database[1].ID,
		Payment2: store.Database[1],
		ID3:      store.Database[2].ID,
		Payment3: store.Database[2],
		Store:    store,
	}
}

func TestPaymentInMemStoreTotal(t *testing.T) {
	db := newTestDBInMem()
	store := db.Store
	total := store.Total()
	assert.Equal(t, 3, total)
}

func TestPaymentInMemStoreGetMany(t *testing.T) {
	db := newTestDBInMem()
	store := db.Store

	payments, err := store.GetMany(0, 0)
	assert.NoError(t, err)
	assert.Len(t, payments, db.Total)

	payments, err = store.GetMany(0, db.Total)
	assert.NoError(t, err)
	assert.Len(t, payments, db.Total)
}

func TestPaymentInMemStoreGetOneNotFound(t *testing.T) {
	var store PaymentStore = NewPaymentInMemStore()
	id := uuid.New()
	if _, err := store.GetOne(id); err != nil && err != ErrNotFound {
		t.Fatalf("Unexpected error while retrieving payment: %v", err)
	}
}

func TestPaymentInMemStoreGetOne(t *testing.T) {
	var store PaymentStore = NewPaymentInMemStore()
	id := uuid.New()
	payment, err := store.GetOne(id)
	if err != nil {
		t.Fatalf("Unexpected error while retrieving payment: %v", err)
	}
	want := "Payment label for TestPaymentInMemStoreGetOne"
	have := payment.Label
	if have != want {
		t.Fatalf("Incorrect label: want %s have %s", want, have)
	}
}
