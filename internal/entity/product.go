package entity

import (
	"errors"
	"time"

	"github.com/MrHenri/go-api/pkg/entity"
)

var (
	ErrIDIsRequired    = errors.New("ID IS REQUIRED")
	ErrInvalidId       = errors.New("INVALID ID")
	ErrNameIsRequired  = errors.New("NAME IS REQUIRED")
	ErrPriceIsRequired = errors.New("PRICE IS REQUIRED")
	ErrInvalidPrice    = errors.New("INVALID PRICE")
)

type Product struct {
	ID        entity.ID
	Name      string
	Price     float64
	CreatedAt time.Time
}

func NewProduct(name string, price float64) (*Product, error) {
	product := &Product{
		ID:        entity.NewId(),
		Name:      name,
		Price:     price,
		CreatedAt: time.Now(),
	}
	err := product.Validate()
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (p *Product) Validate() error {
	if p.ID.String() == "" {
		return ErrIDIsRequired
	}
	if _, err := entity.ParseID(p.ID.String()); err != nil {
		return ErrInvalidId
	}
	if p.Name == "" {
		return ErrNameIsRequired
	}
	if p.Price == 0 {
		return ErrPriceIsRequired
	}
	if p.Price < 0 {
		return ErrInvalidPrice
	}
	return nil
}
