package api

import (
	"context"
)

type APIContext struct {
	context.Context
	Error error
	Data  interface{}
	Code  int
}

const (
	APIContextKeyError = iota
	APIContextKeyCode
	APIContextKeyData
)

func NewAPIContext(ctx context.Context) *APIContext {
	return &APIContext{
		Context: ctx,
	}
}

func (c *APIContext) SetError(val error) *APIContext {
	c.Error = val
	if c.Context != nil {
		c.Context = context.WithValue(c.Context, APIContextKeyError, val)
	}
	return c
}

func (c *APIContext) SetData(val interface{}) *APIContext {
	c.Data = val
	if c.Context != nil {
		c.Context = context.WithValue(c.Context, APIContextKeyData, val)
	}
	return c
}

func (c *APIContext) SetCode(val int) *APIContext {
	c.Code = val
	if c.Context != nil {
		c.Context = context.WithValue(c.Context, APIContextKeyCode, val)
	}
	return c
}
