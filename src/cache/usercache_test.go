package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAddNewUser(t *testing.T) {
	um := NewUserManager()
	userID, err := um.AddNewUser(2 * time.Second)
	assert.NoError(t, err)
	assert.NotEmpty(t, userID)
}

func TestAddOrUpdateUserCache(t *testing.T) {
	um := NewUserManager()
	userID, _ := um.AddNewUser(2 * time.Second)

	fruit := map[string]string{
		"name":  "Apple",
		"color": "red",
	}

	err := um.AddOrUpdateUserCache(userID, "Fruits", fruit, 1*time.Second)
	assert.NoError(t, err)

	val, err := um.ReadDataFromCache(userID, "Fruits")
	assert.NoError(t, err)
	assert.Equal(t, fruit, val)
}

func TestReadUser(t *testing.T) {
	um := NewUserManager()
	userID, _ := um.AddNewUser(2 * time.Second)

	_, err := um.ReadUser(userID)
	assert.NoError(t, err)
}

func TestReadDataFromCacheExpired(t *testing.T) {
	um := NewUserManager()
	userID, _ := um.AddNewUser(1 * time.Second)

	um.AddOrUpdateUserCache(userID, "key1", "value1", 1*time.Second)
	time.Sleep(2 * time.Second)

	val, err := um.ReadDataFromCache(userID, "key1")
	assert.Error(t, err)
	assert.Nil(t, val)
}

func TestUserFlush(t *testing.T) {
	um := NewUserManager()
	userID, _ := um.AddNewUser(1 * time.Second)

	um.UserFlush(500 * time.Millisecond)
	time.Sleep(2 * time.Second)

	_, err := um.ReadUser(userID)
	assert.Error(t, err)
}

func TestCacheFlush(t *testing.T) {
	um := NewUserManager()
	userID, _ := um.AddNewUser(2 * time.Second)

	um.AddOrUpdateUserCache(userID, "key1]", "value1", 1*time.Second)
	um.CacheFlush(500 * time.Millisecond)
	time.Sleep(2 * time.Second)

	val, err := um.ReadDataFromCache(userID, "key1")
	assert.Error(t, err)
	assert.Nil(t, val)
}
