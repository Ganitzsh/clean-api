package api_test

import (
	"testing"
	"time"

	"github.com/ganitzsh/f3-te/api"
	"github.com/google/uuid"
	"github.com/icrowley/fake"
	"github.com/stretchr/testify/assert"
)

type mockDB struct {
	Total    int
	ID1      uuid.UUID
	ID2      uuid.UUID
	ID3      uuid.UUID
	Payment1 *api.Payment
	Payment2 *api.Payment
	Payment3 *api.Payment
	Store    api.PaymentStore
}

func newMockPayment() *api.Payment {
	now := time.Now()
	return &api.Payment{
		ID:        uuid.New(),
		CreatedAt: &now,
		UpdatedAt: &now,
		Purpose:   fake.Sentence(),
		Scheme:    fake.Product(),
		Type:      fake.Brand(),
		Amount:    fake.Digits(),
		Beneficiary: &api.PaymentParty{
			AccountName:       fake.FullName(),
			AccountNumber:     fake.Digits(),
			AccountNumberCode: fake.Brand(),
			BankID:            fake.Digits(),
			BankIDCode:        fake.Word(),
			Name:              fake.FullName(),
			Address:           fake.StreetAddress(),
		},
		Currency: fake.CurrencyCode(),
		DebitorParty: &api.PaymentParty{
			AccountName:       fake.FullName(),
			AccountNumber:     fake.Digits(),
			AccountNumberCode: fake.Brand(),
			BankID:            fake.Digits(),
			BankIDCode:        fake.Word(),
			Name:              fake.FullName(),
			Address:           fake.StreetAddress(),
		},
		EndToEndReference:    fake.Digits(),
		NumericReference:     fake.Digits(),
		ProcessingDate:       "01-01-2019",
		Reference:            fake.Digits(),
		SchemePaymentType:    fake.Brand(),
		SchemePaymentSubType: fake.Brand(),
	}
}

var (
	schemeA = "A"
	schemeB = "B"
)

func newTestDBInMem() *mockDB {
	store := api.NewPaymentInMemStore()
	store.Database = append(store.Database, []*api.Payment{
		newMockPayment().SetScheme(schemeA),
		newMockPayment().SetScheme(schemeA),
		newMockPayment().SetScheme(schemeB),
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

func testPaymentStoreGetMany(db *mockDB) func(*testing.T) {
	return func(t *testing.T) {
		store := db.Store

		payments, err := store.GetMany(0, 0)
		assert.NoError(t, err)
		assert.Len(t, payments, db.Total)

		payments, err = store.GetMany(0, db.Total)
		assert.NoError(t, err)
		assert.Len(t, payments, 0)

		payments, err = store.GetMany(1, 0)
		assert.NoError(t, err)
		assert.Len(t, payments, 1)
		assert.Equal(t, db.ID1, db.Payment1.ID)

		payments, err = store.GetMany(1, 1)
		assert.NoError(t, err)
		assert.Len(t, payments, 1)
		assert.Equal(t, db.ID2, db.Payment2.ID)

		payments, err = store.GetMany(1, 2)
		assert.NoError(t, err)
		assert.Len(t, payments, 1)
		assert.Equal(t, db.ID3, db.Payment3.ID)

		payments, err = store.GetMany(0, 0, api.PaymentStoreFilterIsScheme(schemeA))
		assert.NoError(t, err)
		assert.Len(t, payments, 2)
		assert.Equal(t, schemeA, payments[0].Scheme)
		assert.Equal(t, schemeA, payments[1].Scheme)
	}
}

func testPaymentStoreTotal(db *mockDB) func(*testing.T) {
	return func(t *testing.T) {
		store := db.Store
		total := store.Total()
		assert.Equal(t, 3, total)
	}
}

func testPaymentStoreSave(db *mockDB) func(*testing.T) {
	return func(t *testing.T) {
		store := db.Store
		newPayment := api.NewPayment()
		assert.NoError(t, store.Save(newPayment))

		fromDB, err := store.GetByID(newPayment.ID)
		assert.NoError(t, err)
		assert.Equal(t, newPayment.ID, fromDB.ID)

		assert.Error(t, store.Save(nil))
		assert.NoError(t, store.Save(&api.Payment{}))
	}
}

func testPaymentStoreDelete(db *mockDB) func(*testing.T) {
	return func(t *testing.T) {
		store := db.Store
		assert.NoError(t, store.Delete(db.ID1))
		assert.EqualError(t, store.Delete(uuid.New()), api.ErrNotFound.Error())
		_, err := store.GetByID(db.ID1)
		assert.EqualError(t, err, api.ErrNotFound.Error())
	}
}

func TestPaymentInMemStoreTotal(t *testing.T) {
	testPaymentStoreTotal(newTestDBInMem())(t)
}

func TestPaymentInMemStoreGetMany(t *testing.T) {
	testPaymentStoreGetMany(newTestDBInMem())(t)
}

func TestPaymentInMemStoreSave(t *testing.T) {
	testPaymentStoreSave(newTestDBInMem())(t)
}
