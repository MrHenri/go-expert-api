package database

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/MrHenri/go-api/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestProductCreate(t *testing.T) {
	assert := assert.New(t)

	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	assert.Nil(err)
	assert.NotNil(db)

	err = db.AutoMigrate(&entity.Product{})
	assert.Nil(err)

	productDB := NewProduct(db)
	assert.NotNil(productDB)

	product, err := entity.NewProduct("computer", 3499.0)
	if err != nil {
		t.Error(err)
	}

	err = productDB.Create(product)
	assert.Nil(err)

	var productResult entity.Product

	result := db.Find(&productResult, "id = ?", product.ID)
	assert.Greater(result.RowsAffected, int64(0))
	assert.Nil(result.Error)
	assert.Equal(product.Name, productResult.Name)
	assert.Equal(product.Price, productResult.Price)
	assert.Equal(product.CreatedAt.Day(), productResult.CreatedAt.Day())
}

func TestFindByID(t *testing.T) {
	assert := assert.New(t)

	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	assert.Nil(err)
	assert.NotNil(db)

	err = db.AutoMigrate(&entity.Product{})
	assert.Nil(err)

	productDB := NewProduct(db)
	assert.NotNil(productDB)

	product, err := entity.NewProduct("computer", 3499.0)
	if err != nil {
		t.Error(err)
	}

	err = productDB.Create(product)
	assert.Nil(err)

	productResult, err := productDB.FindById(product.ID.String())
	assert.Nil(err)
	assert.NotNil(productResult)
	assert.Equal(product.Name, productResult.Name)
	assert.Equal(product.Price, productResult.Price)
	assert.Equal(product.CreatedAt.Day(), productResult.CreatedAt.Day())

	productFailResult, err := productDB.FindById("failID")
	assert.NotNil(err)
	assert.Nil(productFailResult)
}

func TestUpdate(t *testing.T) {
	assert := assert.New(t)

	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	assert.Nil(err)
	assert.NotNil(db)

	err = db.AutoMigrate(&entity.Product{})
	assert.Nil(err)

	productDB := NewProduct(db)
	assert.NotNil(productDB)

	product, err := entity.NewProduct("computer", 3499.0)
	if err != nil {
		t.Error(err)
	}

	err = productDB.Create(product)
	assert.Nil(err)

	product.Name = "Iphone"
	err = productDB.Update(product)
	assert.Nil(err)

	productResult, err := productDB.FindById(product.ID.String())
	assert.Nil(err)
	assert.NotNil(productResult)
	assert.Equal(product.Name, productResult.Name)

	failProduct, err := entity.NewProduct("failProduct", 1.0)
	assert.Nil(err)
	assert.NotNil(failProduct)
	assert.NotEmpty(failProduct.ID.String())

	err = productDB.Update(failProduct)
	assert.NotNil(err)
}

func TestDelete(t *testing.T) {
	assert := assert.New(t)

	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	assert.Nil(err)
	assert.NotNil(db)

	err = db.AutoMigrate(&entity.Product{})
	assert.Nil(err)

	productDB := NewProduct(db)
	assert.NotNil(productDB)

	product, err := entity.NewProduct("computer", 3499.0)
	if err != nil {
		t.Error(err)
	}

	err = productDB.Create(product)
	assert.Nil(err)

	err = productDB.Delete(product.ID.String())
	assert.Nil(err)
}

func TestFinalAllProducts(t *testing.T) {
	assert := assert.New(t)

	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Product{})
	for i := 1; i < 24; i++ {
		product, err := entity.NewProduct(fmt.Sprintf("Product %d", i), rand.Float64()*100)
		assert.NoError(err)
		db.Create(product)
	}
	productDB := NewProduct(db)
	products, err := productDB.FindAll(1, 10, "asc")
	assert.NoError(err)
	assert.Len(products, 10)
	assert.Equal("Product 1", products[0].Name)
	assert.Equal("Product 10", products[9].Name)

	products, err = productDB.FindAll(2, 10, "asc")
	assert.NoError(err)
	assert.Len(products, 10)
	assert.Equal("Product 11", products[0].Name)
	assert.Equal("Product 20", products[9].Name)

	products, err = productDB.FindAll(3, 10, "asc")
	assert.NoError(err)
	assert.Len(products, 3)
	assert.Equal("Product 21", products[0].Name)
	assert.Equal("Product 23", products[2].Name)

	products, err = productDB.FindAll(0, 0, "x")
	assert.NoError(err)
	assert.Len(products, 23)
	assert.Equal("Product 1", products[0].Name)
	assert.Equal("Product 23", products[22].Name)
}
