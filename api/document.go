package api

import (
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"
)

type APIVersion int32

type Model interface {
	Validate() error
}

type Document struct {
	ID         uuid.UUID   `json:"id" validate:"required"`
	Type       string      `json:"type" validate:"required"`
	CreatedAt  time.Time   `json:"createdAt" validate:"required"`
	APIVersion APIVersion  `json:"apiVersion" validate:"required"`
	Attributes interface{} `json:"attributes"`
}

func NewDocument() *Document {
	return &Document{
		CreatedAt:  time.Now(),
		ID:         uuid.New(),
		APIVersion: CurrentAPIVersion,
	}
}

func (d *Document) SetAttributes(attr Model) *Document {
	d.Attributes = attr
	return d
}

type DocumentStore interface {
	Total() int
	GetMany(limit, offset int, filters ...*DocumentStoreFilter) ([]*Document, error)
	GetByID(id uuid.UUID) (*Document, error)
	Save(p *Document) error
	Delete(id uuid.UUID) error
}

type DocumentStoreFilterType uint

const (
	DocumentStoreFilterTypeEqual = iota
	DocumentStoreFilterTypeIn
)

type DocumentStoreFilter struct {
	Field string
	Want  interface{}
	Type  DocumentStoreFilterType
}

func NewDocumentStoreFilter() *DocumentStoreFilter {
	return &DocumentStoreFilter{}
}

func (sf *DocumentStoreFilter) SetWant(value interface{}) *DocumentStoreFilter {
	sf.Want = value
	return sf
}

func (sf *DocumentStoreFilter) SetType(value DocumentStoreFilterType) *DocumentStoreFilter {
	sf.Type = value
	return sf
}

func (f DocumentStoreFilter) Match(has interface{}) (bool, error) {
	if f.Want == nil && has == nil {
		return false, ErrSomethingWentWrong(nil)
	}
	switch f.Type {
	case DocumentStoreFilterTypeEqual:
		return reflect.DeepEqual(has, f.Want), nil
	case DocumentStoreFilterTypeIn:
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

func DocumentStoreFilterIsType(typ string) *DocumentStoreFilter {
	return &DocumentStoreFilter{
		Field: "Type",
		Want:  typ,
		Type:  DocumentStoreFilterTypeEqual,
	}
}

func (store *DocumentInMemStore) GetMany(
	limit, offset int,
	filters ...*DocumentStoreFilter,
) ([]*Document, error) {
	to := limit + offset
	if to > len(store.Database) {
		return []*Document{}, nil
	}
	subset := []*Document{}
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

func (store *DocumentInMemStore) GetByID(id uuid.UUID) (*Document, error) {
	for _, p := range store.Database {
		if p.ID.String() == id.String() {
			return p, nil
		}
	}
	return nil, ErrNotFound
}

func (store *DocumentInMemStore) Save(d *Document) error {
	if d == nil {
		return ErrNilValue
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

func (store *DocumentInMemStore) Delete(id uuid.UUID) error {
	for i, storedDocument := range store.Database {
		if storedDocument.ID.String() == id.String() {
			store.Database = store.Database[:i+copy(store.Database[i:], store.Database[i+1:])]
			return nil
		}
	}
	return ErrNotFound
}
