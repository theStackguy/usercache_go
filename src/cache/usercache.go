package cache

import (
	"fmt"
	"sync"
	"time"

	"github.com/theStackguy/usercache_go/src/logs"
	"github.com/theStackguy/usercache_go/src/util"
)

type userManager struct {
	users map[string]*user
	mu    sync.RWMutex
}

type userPayload struct {
	id    string
	key   string
	value any
}

type user struct {
	id        string
	cache     map[string]cacheItem
	mu        sync.RWMutex
	expiresAt time.Time
}

type cacheItem struct {
	value      any
	expiryTime time.Time
}

func NewUserManager() *userManager {
	um := &userManager{
		users: make(map[string]*user),
	}
	um.userCacheCleanup(4 * time.Hour)
	return um
}

func (um *userManager) AddNewUser(userExpirationTime time.Duration) (string, error) {
	userId, err := util.NewString()
	if err != nil {
		return "", logs.ErrGuid
	}
	um.mu.Lock()
	um.users[userId] = &user{
		id:        userId,
		cache:     make(map[string]cacheItem),
		expiresAt: time.Now().Add(userExpirationTime),
	}
	um.mu.Unlock()
	return userId, nil
}

func (um *userManager) AddOrUpdateUserCache(usertoken string, key string, value any, expirationTime time.Duration) error {
	up := userPayload{
		id:    usertoken,
		key:   key,
		value: value,
	}
	if up.HasAllNeededData(true) {
		um.mu.RLock()
		user, ok := um.users[usertoken]
		um.mu.RUnlock()
		if !ok {
			fmt.Println("User not found:", usertoken)
			return logs.ErrUser
		}
		if hasExpired(user.expiresAt) {
			return logs.ErrUserExpired
		}
		done := user.setCache(key, value, expirationTime)
		if !done {
			return logs.ErrCacheUpdate
		}

		return nil
	}
	return logs.ErrPayLoad
}

func (u *user) setCache(key string, value any, ttl time.Duration) bool {
	var expiresAt time.Time
	u.mu.RLock()
	responseData, found := u.cache[key]
	u.mu.RUnlock()
	if found {
		if hasExpired(responseData.expiryTime) {
			return false
		}
		expiresAt = responseData.expiryTime.Add(ttl)
	} else {
		expiresAt = time.Now().Add(ttl)
	}
	if ttl == 0 {
		expiresAt = time.Time{}
	}
	u.mu.Lock()
	u.cache[key] = cacheItem{
		value:      value,
		expiryTime: expiresAt,
	}
	if ttl == 0 {
		u.expiresAt = time.Time{}
	} else {
		u.expiresAt = u.expiresAt.Add(ttl)
	}
	u.mu.Unlock()

	return true
}

func (um *userManager) ReadUser(usertoken string) (any, error) {

	if usertoken != "" {
		um.mu.RLock()
		user, ok := um.users[usertoken]
		um.mu.RUnlock()
		if !ok {
			return nil, logs.ErrUser
		}
		if hasExpired(user.expiresAt) {
			return nil, logs.ErrReadUser
		}
		for key, cache := range user.cache {
			if time.Now().After(cache.expiryTime) && !cache.expiryTime.IsZero() {
				user.mu.Lock()
				delete(user.cache, key)
				user.mu.Unlock()
			}
		}
		return user, nil
	}
	return nil, logs.ErrReadUserToken
}

func (um *userManager) ReadDataFromCache(usertoken string, key string) (any, error) {

	if usertoken != "" {
		um.mu.RLock()
		user, ok := um.users[usertoken]
		um.mu.RUnlock()
		if !ok {
			return nil, logs.ErrUser
		}
		if hasExpired(user.expiresAt) {
			um.mu.Lock()
			delete(um.users, usertoken)
			um.mu.Unlock()
			return nil, logs.ErrReadUser
		}
		value, done := user.get(key)

		if done != nil {
			return nil, done
		}
		return value, nil

	}
	return nil, logs.ErrReadUserToken
}

func (u *user) get(key string) (any, error) {
	u.mu.RLock()
	data, ok := u.cache[key]
	u.mu.RUnlock()
	if !ok {
		return nil, logs.ErrReadCacheKey
	}
	if hasExpired(data.expiryTime) {
		u.mu.Lock()
		delete(u.cache, key)
		u.mu.Unlock()
		return nil, logs.ErrCacheExpired
	}
	return data.value, nil
}

func (um *userManager) UserFlush(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for range ticker.C {
			now := time.Now()
			um.mu.Lock()
			for id, user := range um.users {
				if now.After(user.expiresAt) {
					delete(um.users, id)
				}
			}
			um.mu.Unlock()
		}
	}()
}

func (um *userManager) CacheFlush(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for range ticker.C {
			now := time.Now()
			um.mu.RLock()
			for _, user := range um.users {
				user.mu.Lock()
				for key, value := range user.cache {
					if !value.expiryTime.IsZero() && now.After(value.expiryTime) {
						delete(user.cache, key)
					}
				}
				user.mu.Unlock()
			}
			um.mu.RUnlock()
		}
	}()
}

func (um *userManager) userCacheCleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for range ticker.C {
			now := time.Now()
			um.mu.Lock()
			for id, user := range um.users {
				if now.After(user.expiresAt) {
					delete(um.users, id)
				}
			}
			um.mu.Unlock()
			um.mu.RLock()
			for _, user := range um.users {
				user.mu.Lock()
				for key, value := range user.cache {
					if !value.expiryTime.IsZero() && now.After(value.expiryTime) {
						delete(user.cache, key)
					}
				}
				user.mu.Unlock()
			}
			um.mu.RUnlock()
		}
	}()
}

func (um *userManager) FlushAll() {
	um.mu.Lock()
	for id := range um.users {
		delete(um.users, id)
	}
	um.mu.Unlock()
}

func (um *userManager) RemoveUser(usertoken string) bool {
	um.mu.RLock()
	_, ok := um.users[usertoken]
	um.mu.RUnlock()
	if !ok {
		return false
	}
	um.mu.Lock()
	delete(um.users, usertoken)
	um.mu.Unlock()
	return true
}

func (um *userManager) RemoveUserCache(usertoken string, key string) bool {
	um.mu.RLock()
	user, ok := um.users[usertoken]
	um.mu.RUnlock()
	if !ok {
		return false
	}
	user.mu.RLock()
	value, done := user.cache[key]
	user.mu.RUnlock()
	if !done {
		return false
	}
	user.mu.Lock()
	difference := time.Until(value.expiryTime)
	user.expiresAt = user.expiresAt.Add(-difference)
	delete(user.cache, key)
	user.mu.Unlock()
	return true
}

func hasExpired(ttl time.Time) bool {
	if ttl.IsZero() {
		return false
	}
	return time.Now().After(ttl)
}

func (c userPayload) HasAllNeededData(flag bool) bool {
	switch {
	case c.id == "":
	case c.key == "":
	case flag && c.value == nil:
	default:
		return true
	}
	return false
}
