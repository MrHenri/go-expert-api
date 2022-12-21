package database

import (
	"testing"

	"github.com/MrHenri/go-api/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestUserCreate(t *testing.T) {
	user, err := entity.NewUser("John", "john@gmail.com", "12345")
	assert.Nil(t, err)
	assert.NotNil(t, user)

	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	assert.NotNil(t, db)
	assert.Nil(t, err)

	err = db.AutoMigrate(&entity.User{})
	assert.Nil(t, err)

	userDb := NewUser(db)
	assert.NotNil(t, userDb)

	err = userDb.Create(user)
	assert.Nil(t, err)

	var resultUser entity.User
	result := db.First(&resultUser, "id = ?", user.ID)
	assert.Nil(t, result.Error)
	assert.Greater(t, result.RowsAffected, int64(0))
	assert.Equal(t, user.Name, resultUser.Name)
	assert.Equal(t, user.Email, resultUser.Email)
	assert.NotEmpty(t, resultUser.Password)
}

func TestFindByEmail(t *testing.T) {
	user, _ := entity.NewUser("John", "john@gmail.com", "12345")

	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	assert.Nil(t, err)

	db.AutoMigrate(&entity.User{})

	userDb := NewUser(db)
	userDb.Create(user)

	resultUser, err := userDb.FindByEmail(user.Email)
	assert.Nil(t, err)
	assert.NotNil(t, resultUser)
	assert.Equal(t, user.Name, resultUser.Name)
	assert.Equal(t, user.Email, resultUser.Email)
	assert.NotEmpty(t, resultUser.Password)

	failResultUser, err := userDb.FindByEmail("inexistent@user.com")
	assert.Nil(t, failResultUser)
	assert.NotNil(t, err)
}
