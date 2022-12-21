package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProduct(t *testing.T) {
	p, err := NewProduct("computer", 3499.00)
	assert.NotNil(t, p)
	assert.Nil(t, err)
	assert.Equal(t, "computer", p.Name)
	assert.Equal(t, 3499.00, p.Price)
	assert.NotEmpty(t, p.ID)
	assert.NotEmpty(t, p.CreatedAt)
}

func TestValidate(t *testing.T) {
	p, err := NewProduct("computer", 3499.00)
	assert.NotNil(t, p)
	assert.Nil(t, err)
	assert.Nil(t, p.Validate())
}

func TestProductWhenNameIsRequired(t *testing.T) {
	p, err := NewProduct("", 3499.00)
	assert.Nil(t, p)
	assert.NotNil(t, err)
	assert.Equal(t, ErrNameIsRequired, err)
}

func TestProductWhenPriceIsRequired(t *testing.T) {
	p, err := NewProduct("computer", 0.00)
	assert.Nil(t, p)
	assert.NotNil(t, err)
	assert.Equal(t, ErrPriceIsRequired, err)
}

func TestProductWhenPriceIsInvalid(t *testing.T) {
	p, err := NewProduct("computer", -3499.00)
	assert.Nil(t, p)
	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidPrice, err)
}
