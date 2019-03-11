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
	Total     int
	ID1       uuid.UUID
	ID2       uuid.UUID
	ID3       uuid.UUID
	ID4       uuid.UUID
	Document1 *api.Document
	Document2 *api.Document
	Document3 *api.Document
	Document4 *api.Document
	Store     api.DocumentStore
}

func newMockPayment() *api.Payment {
	return &api.Payment{
		Amount: fake.Digits(),
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
		PaymentID:            fake.Digits(),
		PaymentPurpose:       fake.Sentence(),
		PaymentScheme:        fake.Product(),
		PaymentType:          fake.Brand(),
		ProcessingDate:       "01-01-2019",
		Reference:            fake.Digits(),
		SchemePaymentType:    fake.Brand(),
		SchemePaymentSubType: fake.Brand(),
	}
}

func newTestDBInMem() *mockDB {
	store := api.NewDocumentInMemStore()
	now := time.Now()
	store.Database = append(store.Database, []*api.Document{
		&api.Document{ID: uuid.New(), APIVersion: api.CurrentAPIVersion, CreatedAt: now, Type: "Payment", Attributes: newMockPayment()},
		&api.Document{ID: uuid.New(), APIVersion: api.CurrentAPIVersion, CreatedAt: now, Type: "Payment", Attributes: newMockPayment()},
		&api.Document{ID: uuid.New(), APIVersion: api.CurrentAPIVersion, CreatedAt: now, Type: "Payment", Attributes: newMockPayment()},
		&api.Document{ID: uuid.New(), APIVersion: api.CurrentAPIVersion, CreatedAt: now, Type: "Other"},
	}...)
	return &mockDB{
		Total:     4,
		ID1:       store.Database[0].ID,
		Document1: store.Database[0],
		ID2:       store.Database[1].ID,
		Document2: store.Database[1],
		ID3:       store.Database[2].ID,
		Document3: store.Database[2],
		ID4:       store.Database[3].ID,
		Document4: store.Database[3],
		Store:     store,
	}
}

func testDocumentStoreGetMany(db *mockDB) func(*testing.T) {
	return func(t *testing.T) {
		store := db.Store

		docs, err := store.GetMany(0, 0)
		assert.NoError(t, err)
		assert.Len(t, docs, db.Total)

		docs, err = store.GetMany(0, db.Total)
		assert.NoError(t, err)
		assert.Len(t, docs, 0)

		docs, err = store.GetMany(1, 0)
		assert.NoError(t, err)
		assert.Len(t, docs, 1)
		assert.Equal(t, db.ID1, db.Document1.ID)

		docs, err = store.GetMany(1, 1)
		assert.NoError(t, err)
		assert.Len(t, docs, 1)
		assert.Equal(t, db.ID2, db.Document2.ID)

		docs, err = store.GetMany(1, 2)
		assert.NoError(t, err)
		assert.Len(t, docs, 1)
		assert.Equal(t, db.ID3, db.Document3.ID)

		docs, err = store.GetMany(0, 0, api.DocumentStoreFilterIsType("Other"))
		assert.NoError(t, err)
		assert.Len(t, docs, 1)
		assert.NotNil(t, docs[0])
		assert.Equal(t, "Other", docs[0].Type)

		docs, err = store.GetMany(0, 0, api.DocumentStoreFilterIsType("Payment"))
		assert.NoError(t, err)
		assert.Len(t, docs, 3)
		assert.Equal(t, "Payment", docs[0].Type)
		assert.Equal(t, "Payment", docs[1].Type)
		assert.Equal(t, "Payment", docs[2].Type)
	}
}

func testDocumentStoreTotal(db *mockDB) func(*testing.T) {
	return func(t *testing.T) {
		store := db.Store
		total := store.Total()
		assert.Equal(t, 4, total)
	}
}

func testDocumentStoreSave(db *mockDB) func(*testing.T) {
	return func(t *testing.T) {
		store := db.Store
		newDocument := api.NewDocument().SetAttributes(newMockPayment())
		assert.NoError(t, store.Save(newDocument))

		fromDB, err := store.GetByID(newDocument.ID)
		assert.NoError(t, err)
		assert.Equal(t, newDocument.ID, fromDB.ID)

		assert.Error(t, store.Save(nil))
		assert.NoError(t, store.Save(&api.Document{}))
	}
}

func testDocumentStoreDelete(db *mockDB) func(*testing.T) {
	return func(t *testing.T) {
		store := db.Store
		assert.NoError(t, store.Delete(db.ID1))
		assert.EqualError(t, store.Delete(uuid.New()), api.ErrNotFound.Error())
		_, err := store.GetByID(db.ID1)
		assert.EqualError(t, err, api.ErrNotFound.Error())
	}
}

func TestDocumentInMemStoreTotal(t *testing.T) {
	testDocumentStoreTotal(newTestDBInMem())(t)
}

func TestDocumentInMemStoreGetMany(t *testing.T) {
	testDocumentStoreGetMany(newTestDBInMem())(t)
}

func TestDocumentInMemStoreSave(t *testing.T) {
	testDocumentStoreSave(newTestDBInMem())(t)
}
