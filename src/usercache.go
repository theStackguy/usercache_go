//package cache

// import (
// 	"fmt"
// 	"sync"
// 	"time"

// 	"github.com/theStackguy/usercache_go/src/logs"
// 	"github.com/theStackguy/usercache_go/src/util"
// )

// type UserManager struct {
// 	Users map[string]*User
// 	Mu    sync.RWMutex
// }

// type userPayload struct {
// 	Id    string
// 	Key   string
// 	Value any
// }

// type User struct {
// 	Id        string
// 	Cache     map[string]CacheItem
// 	Mu        sync.RWMutex
// 	ExpiresAt time.Time
// }

// type CacheItem struct {
// 	Value      any
// 	ExpiryTime time.Time
// }

// func NewUserManager() *UserManager {
// 	um := &UserManager{
// 		Users: make(map[string]*User),
// 	}
// 	um.userCacheCleanup(4 * time.Hour)
// 	return um
// }

// func (um *UserManager) AddNewUser(userExpirationTime time.Duration) (string, error) {
// 	userId, err := util.NewString()
// 	if err != nil {
// 		return "", logs.ErrGuid
// 	}
// 	um.Mu.Lock()
// 	um.Users[userId] = &User{
// 		Id:        userId,
// 		Cache:     make(map[string]CacheItem),
// 		ExpiresAt: time.Now().Add(userExpirationTime),
// 	}
// 	um.Mu.Unlock()
// 	return userId, nil
// }

// func (um *UserManager) AddOrUpdateUserCache(usertoken string, key string, value any, expirationTime time.Duration) error {
// 	up := userPayload{
// 		Id:    usertoken,
// 		Key:   key,
// 		Value: value,
// 	}
// 	if up.HasAllNeededData(true) {
// 		um.Mu.RLock()
// 		user, ok := um.Users[usertoken]
// 		um.Mu.RUnlock()
// 		if !ok {
// 			fmt.Println("User not found:", usertoken)
// 			return logs.ErrUser
// 		}
// 		if hasExpired(user.ExpiresAt) {
// 			return logs.ErrUserExpired
// 		}
// 		done := user.setCache(key, value, expirationTime)
// 		if !done {
// 			return logs.ErrCacheUpdate
// 		}

// 		return nil
// 	}
// 	return logs.ErrPayLoad
// }

// func (u *User) setCache(key string, value any, ttl time.Duration) bool {
// 	var expiresAt time.Time
// 	u.Mu.RLock()
// 	responseData, found := u.Cache[key]
// 	u.Mu.RUnlock()
// 	if found {
// 		if hasExpired(responseData.ExpiryTime) {
// 			return false
// 		}
// 		expiresAt = responseData.ExpiryTime.Add(ttl)
// 	} else {
// 		expiresAt = time.Now().Add(ttl)
// 	}
// 	if ttl == 0 {
// 		expiresAt = time.Time{}
// 	}
// 	u.Mu.Lock()
// 	u.Cache[key] = CacheItem{
// 		Value:      value,
// 		ExpiryTime: expiresAt,
// 	}
// 	if ttl == 0 {
// 		u.ExpiresAt = time.Time{}
// 	} else {
// 		u.ExpiresAt = u.ExpiresAt.Add(ttl)
// 	}
// 	u.Mu.Unlock()

// 	return true
// }

// func (um *UserManager) ReadUser(usertoken string) (any, error) {

// 	if usertoken != "" {
// 		um.Mu.RLock()
// 		user, ok := um.Users[usertoken]
// 		um.Mu.RUnlock()
// 		if !ok {
// 			return nil, logs.ErrUser
// 		}
// 		if hasExpired(user.ExpiresAt) {
// 			return nil, logs.ErrReadUser
// 		}
// 		for key, cache := range user.Cache {
// 			if time.Now().After(cache.ExpiryTime) && !cache.ExpiryTime.IsZero() {
// 				user.Mu.Lock()
// 				delete(user.Cache, key)
// 				user.Mu.Unlock()
// 			}
// 		}
// 		return user, nil
// 	}
// 	return nil, logs.ErrReadUserToken
// }

// func (um *UserManager) ReadDataFromCache(usertoken string, key string) (any, error) {

// 	if usertoken != "" {
// 		um.Mu.RLock()
// 		user, ok := um.Users[usertoken]
// 		um.Mu.RUnlock()
// 		if !ok {
// 			return nil, logs.ErrUser
// 		}
// 		if hasExpired(user.ExpiresAt) {
// 			um.Mu.Lock()
// 			delete(um.Users, usertoken)
// 			um.Mu.Unlock()
// 			return nil, logs.ErrReadUser
// 		}
// 		value, done := user.get(key)

// 		if done != nil {
// 			return nil, done
// 		}
// 		return value, nil

// 	}
// 	return nil, logs.ErrReadUserToken
// }

// func (u *User) get(key string) (any, error) {
// 	u.Mu.RLock()
// 	data, ok := u.Cache[key]
// 	u.Mu.RUnlock()
// 	if !ok {
// 		return nil, logs.ErrReadCacheKey
// 	}
// 	if hasExpired(data.ExpiryTime) {
// 		u.Mu.Lock()
// 		delete(u.Cache, key)
// 		u.Mu.Unlock()
// 		return nil, logs.ErrCacheExpired
// 	}
// 	return data.Value, nil
// }

// func (um *UserManager) UserFlush(interval time.Duration) {
// 	ticker := time.NewTicker(interval)
// 	go func() {
// 		defer ticker.Stop()
// 		for range ticker.C {
// 			now := time.Now()
// 			um.Mu.Lock()
// 			for id, user := range um.Users {
// 				if now.After(user.ExpiresAt) {
// 					delete(um.Users, id)
// 				}
// 			}
// 			um.Mu.Unlock()
// 		}
// 	}()
// }

// func (um *UserManager) CacheFlush(interval time.Duration) {
// 	ticker := time.NewTicker(interval)
// 	go func() {
// 		defer ticker.Stop()
// 		for range ticker.C {
// 			now := time.Now()
// 			um.Mu.RLock()
// 			for _, user := range um.Users {
// 				user.Mu.Lock()
// 				for key, value := range user.Cache {
// 					if !value.ExpiryTime.IsZero() && now.After(value.ExpiryTime) {
// 						delete(user.Cache, key)
// 					}
// 				}
// 				user.Mu.Unlock()
// 			}
// 			um.Mu.RUnlock()
// 		}
// 	}()
// }

// func (um *UserManager) userCacheCleanup(interval time.Duration) {
// 	ticker := time.NewTicker(interval)
// 	go func() {
// 		defer ticker.Stop()
// 		for range ticker.C {
// 			now := time.Now()
// 			um.Mu.Lock()
// 			for id, user := range um.Users {
// 				if now.After(user.ExpiresAt) {
// 					delete(um.Users, id)
// 				}
// 			}
// 			um.Mu.Unlock()
// 			um.Mu.RLock()
// 			for _, user := range um.Users {
// 				user.Mu.Lock()
// 				for key, value := range user.Cache {
// 					if !value.ExpiryTime.IsZero() && now.After(value.ExpiryTime) {
// 						delete(user.Cache, key)
// 					}
// 				}
// 				user.Mu.Unlock()
// 			}
// 			um.Mu.RUnlock()
// 		}
// 	}()
// }

// func (um *UserManager) FlushAll() {
// 	um.Mu.Lock()
// 	for id := range um.Users {
// 		delete(um.Users, id)
// 	}
// 	um.Mu.Unlock()
// }

// func (um *UserManager) RemoveUser(usertoken string) bool {
// 	um.Mu.RLock()
// 	_, ok := um.Users[usertoken]
// 	um.Mu.RUnlock()
// 	if !ok {
// 		return false
// 	}
// 	um.Mu.Lock()
// 	delete(um.Users, usertoken)
// 	um.Mu.Unlock()
// 	return true
// }

// func (um *UserManager) RemoveUserCache(usertoken string, key string) bool {
// 	um.Mu.RLock()
// 	user, ok := um.Users[usertoken]
// 	um.Mu.RUnlock()
// 	if !ok {
// 		return false
// 	}
// 	user.Mu.RLock()
// 	value, done := user.Cache[key]
// 	user.Mu.RUnlock()
// 	if !done {
// 		return false
// 	}
// 	user.Mu.Lock()
// 	difference := time.Until(value.ExpiryTime)
// 	user.ExpiresAt = user.ExpiresAt.Add(-difference)
// 	delete(user.Cache, key)
// 	user.Mu.Unlock()
// 	return true
// }

// func hasExpired(ttl time.Time) bool {
// 	if ttl.IsZero() {
// 		return false
// 	}
// 	return time.Now().After(ttl)
// }

// func (c userPayload) HasAllNeededData(flag bool) bool {
// 	switch {
// 	case c.Id == "":
// 	case c.Key == "":
// 	case flag && c.Value == nil:
// 	default:
// 		return true
// 	}
// 	return false
// }
