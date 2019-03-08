package api

import (
	"time"

	"github.com/google/uuid"
)

type APIVersion int32

type Document struct {
	ID         uuid.UUID   `json:"id"`
	Type       string      `json:"type"`
	CreatedAt  time.Time   `json:"createdAt"`
	APIVersion APIVersion  `json:"apiVersion"`
	Attributes interface{} `json:"attributes"`
}

func NewDocument() *Document {
	return &Document{
		CreatedAt:  time.Now(),
		ID:         uuid.New(),
		APIVersion: CurrentAPIVersion,
	}
}

type DocumentStore interface {
	Total() int
	GetMany(limit, offset int) ([]*Document, error)
	GetByID(id uuid.UUID) (*Document, error)
	Create(p *Document) error
}

type DocumentInMemStore struct {
	Database []*Document
}

func NewDocumentInMemStore() *DocumentInMemStore {
	return &DocumentInMemStore{
		Database: []*Document{},
	}
}

func (store *DocumentInMemStore) Total() int {
	return len(store.Database)
}

func (store *DocumentInMemStore) GetMany(limit, offset int) ([]*Document, error) {
	to := limit + offset
	if to > len(store.Database) {
		return []*Document{}, nil
	}
	if offset == 0 && limit == 0 {
		return store.Database, nil
	}
	if offset > 0 && limit == 0 {
		return store.Database[offset:], nil
	}
	return store.Database[offset:to], nil
}

func (store *DocumentInMemStore) GetByID(id uuid.UUID) (*Document, error) {
	for _, p := range store.Database {
		if p.ID.String() == id.String() {
			return p, nil
		}
	}
	return nil, ErrNotFound
}

func (store *DocumentInMemStore) Create(p *Document) error {
	if p == nil {
		return ErrNilValue
	}
	store.Database = append(store.Database, p)
	return nil
}
