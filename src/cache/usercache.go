package cache

import (
	"fmt"
	"sync"
	"time"

	"github.com/theStackguy/usercache_go/src/logs"
	"github.com/theStackguy/usercache_go/src/util"
)

type UserManager struct {
	Users map[string]*user
	Mu   sync.RWMutex
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

func NewUserManager() *UserManager {
	um := &UserManager{
		Users: make(map[string]*user),
	}
	um.userCacheCleanup(4 * time.Hour)
	return um
}

func (um *UserManager) AddNewUser(userExpirationTime time.Duration) (string, error) {
	userId, err := util.NewString()
	if err != nil {
		return "", logs.ErrGuid
	}
	um.Mu.Lock()
	um.Users[userId] = &user{
		id:        userId,
		cache:     make(map[string]cacheItem),
		expiresAt: time.Now().Add(userExpirationTime),
	}
	um.Mu.Unlock()
	return userId, nil
}

func (um *UserManager) AddOrUpdateUserCache(usertoken string, key string, value any, expirationTime time.Duration) error {
	up := userPayload{
		id:    usertoken,
		key:   key,
		value: value,
	}
	if up.HasAllNeededData(true) {
		um.Mu.RLock()
		user, ok := um.Users[usertoken]
		um.Mu.RUnlock()
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

func (um *UserManager) ReadUser(usertoken string) (any, error) {

	if usertoken != "" {
		um.Mu.RLock()
		user, ok := um.Users[usertoken]
		um.Mu.RUnlock()
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

func (um *UserManager) ReadDataFromCache(usertoken string, key string) (any, error) {

	if usertoken != "" {
		um.Mu.RLock()
		user, ok := um.Users[usertoken]
		um.Mu.RUnlock()
		if !ok {
			return nil, logs.ErrUser
		}
		if hasExpired(user.expiresAt) {
			um.Mu.Lock()
			delete(um.Users, usertoken)
			um.Mu.Unlock()
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

func (um *UserManager) UserFlush(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for range ticker.C {
			now := time.Now()
			um.Mu.Lock()
			for id, user := range um.Users {
				if now.After(user.expiresAt) {
					delete(um.Users, id)
				}
			}
			um.Mu.Unlock()
		}
	}()
}

func (um *UserManager) CacheFlush(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for range ticker.C {
			now := time.Now()
			um.Mu.RLock()
			for _, user := range um.Users {
				user.mu.Lock()
				for key, value := range user.cache {
					if !value.expiryTime.IsZero() && now.After(value.expiryTime) {
						delete(user.cache, key)
					}
				}
				user.mu.Unlock()
			}
			um.Mu.RUnlock()
		}
	}()
}

func (um *UserManager) userCacheCleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for range ticker.C {
			now := time.Now()
			um.Mu.Lock()
			for id, user := range um.Users {
				if now.After(user.expiresAt) {
					delete(um.Users, id)
				}
			}
			um.Mu.Unlock()
			um.Mu.RLock()
			for _, user := range um.Users {
				user.mu.Lock()
				for key, value := range user.cache {
					if !value.expiryTime.IsZero() && now.After(value.expiryTime) {
						delete(user.cache, key)
					}
				}
				user.mu.Unlock()
			}
			um.Mu.RUnlock()
		}
	}()
}

func (um *UserManager) FlushAll() {
	um.Mu.Lock()
	for id := range um.Users {
		delete(um.Users, id)
	}
	um.Mu.Unlock()
}

func (um *UserManager) RemoveUser(usertoken string) bool {
	um.Mu.RLock()
	_, ok := um.Users[usertoken]
	um.Mu.RUnlock()
	if !ok {
		return false
	}
	um.Mu.Lock()
	delete(um.Users, usertoken)
	um.Mu.Unlock()
	return true
}

func (um *UserManager) RemoveUserCache(usertoken string, key string) bool {
	um.Mu.RLock()
	user, ok := um.Users[usertoken]
	um.Mu.RUnlock()
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
