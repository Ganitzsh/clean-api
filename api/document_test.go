package api_test

import (
	"testing"

	"github.com/ganitzsh/f3-te/api"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockDB struct {
	Total    int
	ID1      uuid.UUID
	ID2      uuid.UUID
	ID3      uuid.UUID
	Document1 *api.Document
	Document2 *api.Document
	Document3 *api.Document
	Store    api.DocumentStore
}

func newTestDBInMem() *mockDB {
	store := api.NewDocumentInMemStore()
	store.Database = append(store.Database, []*api.Document{
		&api.Document{ID: uuid.New()},
		&api.Document{ID: uuid.New()},
		&api.Document{ID: uuid.New()},
	}...)
	return &mockDB{
		Total:    3,
		ID1:      store.Database[0].ID,
		Document1: store.Database[0],
		ID2:      store.Database[1].ID,
		Document2: store.Database[1],
		ID3:      store.Database[2].ID,
		Document3: store.Database[2],
		Store:    store,
	}
}

func testDocumentStoreGetMany(db *mockDB) func(*testing.T) {
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
		assert.Equal(t, db.ID1, db.Document1.ID)

		payments, err = store.GetMany(1, 1)
		assert.NoError(t, err)
		assert.Len(t, payments, 1)
		assert.Equal(t, db.ID2, db.Document2.ID)

		payments, err = store.GetMany(1, 2)
		assert.NoError(t, err)
		assert.Len(t, payments, 1)
		assert.Equal(t, db.ID3, db.Document3.ID)
	}
}

func testDocumentStoreTotal(db *mockDB) func(*testing.T) {
	return func(t *testing.T) {
		store := db.Store
		total := store.Total()
		assert.Equal(t, 3, total)
	}
}

func testDocumentStoreCreate(db *mockDB) func(*testing.T) {
	return func(t *testing.T) {
		store := db.Store
		newDocument := api.NewDocument()
		assert.NoError(t, store.Create(newDocument))

		fromDB, err := store.GetByID(newDocument.ID)
		assert.NoError(t, err)
		assert.Equal(t, newDocument.ID, fromDB.ID)

		assert.Error(t, store.Create(nil))
		assert.NoError(t, store.Create(&api.Document{}))
	}
}

func testDocumentStoreDelete(db *mockDB) func(*testing.T) {
	return func(t *testing.T) {
		store := db.Store
		newDocument := api.NewDocument()
		assert.NoError(t, store.Create(newDocument))

		fromDB, err := store.GetByID(newDocument.ID)
		assert.NoError(t, err)
		assert.Equal(t, newDocument.ID, fromDB.ID)

		assert.Error(t, store.Create(nil))
		assert.NoError(t, store.Create(&api.Document{}))
	}
}

func TestDocumentInMemStoreTotal(t *testing.T) {
	testDocumentStoreTotal(newTestDBInMem())(t)
}

func TestDocumentInMemStoreGetMany(t *testing.T) {
	testDocumentStoreGetMany(newTestDBInMem())(t)
}

func TestDocumentInMemStoreCreate(t *testing.T) {
	testDocumentStoreCreate(newTestDBInMem())(t)
}
