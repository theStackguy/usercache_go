package src

import (
	"sync"
	"time"
)

type UserManager struct {
	Users map[string]*User
	Mu    sync.RWMutex
}

type User struct {
	Id               string
	Mu               sync.RWMutex
	Sessions         map[string]*Session
	CurrentSessionId string
	SharedCache      *Cache[string, any]
	memory           *memorylimit
	IsActive         bool
}

type Session struct {
	SessionId     string
	SessionToken  string
	RefreshToken  string
	IsActive      bool
	SessionExpiry time.Time
	RefreshExpiry time.Time
	LastAccessed  time.Time
	Cache         *Cache[string, any]
	Mu            sync.RWMutex
	Err           error
}

type cacheItem[T any] struct {
	Value        T
	ExpiryTime   time.Time
	LastAccessed time.Time
}

type Cache[K comparable, V any] struct {
	Mu    sync.Mutex
	Store map[K]cacheItem[V]
}

func newCache[K comparable, V any]() *Cache[K, V] {
	return &Cache[K, V]{
		Store: make(map[K]cacheItem[V]),
	}
}

type userPayload struct {
	Id    string
	Key   string
	Value any
}

func NewUserManager() *UserManager {
	um := &UserManager{
		Users: make(map[string]*User),
	}
	// um.userCacheCleanup(4 * time.Hour)
	return um
}

func (um *UserManager) AddNewUser(sessionTokenExpiryTime time.Duration, refreshTokenExpiryTime time.Duration, memorylimitInMB float64) (*User, error) {

	if memorylimitInMB <= 0 {
		return nil, errMemoryLimit
	}

	var wg sync.WaitGroup
	var osmemorychannel chan error
	var memorychannel chan uint64

	wg.Add(2)

	go operatingSystemAvailableMemory(osmemorychannel, &wg)
	go mbSizeToUINT(memorylimitInMB, memorychannel, &wg)

	tokens, err := newTokenStrings(2)
	if err != nil {
		return nil, errGuid
	}

	var sessionuser *Session
	var userId *string = &tokens[1]
	sessionuser.SessionId = tokens[0]

	(sessionuser).generateSessionRefreshToken(sessionTokenExpiryTime, refreshTokenExpiryTime)
	if sessionuser.Err != nil {
		return nil, sessionuser.Err
	}
	if err != nil {
		return nil, errTokenGen
	}

	sessionuser.LastAccessed = time.Now()
	sessionuser.IsActive = true
	sessionuser.Cache = newCache[string, any]()

	wg.Wait()

	if osmemerror := <-osmemorychannel; osmemerror != nil {
		return nil, osmemerror
	}
	if !compareConfigOsMem(osavailableMemory.Load(), <-memorychannel) {
		return nil, errMemExceeded
	}

	usermemory := &memorylimit{
		configured:     <-memorychannel,
		remainingSpace: osavailableMemory.Load(),
	}

	um.Mu.Lock()
	um.Users[*userId] = &User{
		Id:               *userId,
		Sessions:         map[string]*Session{tokens[0]: sessionuser},
		CurrentSessionId: tokens[0],
		SharedCache:      newCache[string, any](),
		memory:           usermemory,
		IsActive:         true,
	}
	um.Mu.Unlock()
	return um.Users[*userId], nil
}

func (um *UserManager) AddNewSessionToUser(userId string, sessionTokenExpiryTime time.Duration, refreshTokenExpiryTime time.Duration) (*Session, error) {
	um.Mu.RLock()
	user, exists := um.Users[userId]
	um.Mu.RUnlock()
	if !exists {
		return nil, errUser
	}
	sessionId, err := newTokenString()
	if err != nil {
		return nil, errGuid
	}

	var newsession *Session
	newsession.SessionId = sessionId
	(newsession).generateSessionRefreshToken(sessionTokenExpiryTime, refreshTokenExpiryTime)
	newsession.LastAccessed = time.Now()
	newsession.Cache = newCache[string, any]()

	user.Mu.Lock()
	user.Sessions[sessionId] = newsession
	user.CurrentSessionId = sessionId
	user.Mu.Unlock()
	return newsession, nil
}

func (u *User) AddorUpdateUserCache(sessionid, sessionToken, key string, value any) (*Session, error) {

	u.Mu.RLock()
	session, exists := u.Sessions[sessionid]
	u.Mu.RUnlock()
	if !exists {
		return nil, errSession
	}
	(session).checkTokenExpired()
	if session.Err == errAuth {
		RetryAuthentication(session)
	}
	if sessionToken != session.SessionToken {

	}

	updatedsession := s.checkTokenExpired(sessionToken)
	switch {
	case updatedsession == nil:
		return updatedsession, errAddorUpdateCache
	case updatedsession.Err != nil:
		return updatedsession, updatedsession.Err
	}

}

func AddorUpdateSessionCache() {

}

func (c userPayload) hasAllNeededData(flag bool) bool {
	switch {
	case c.Id == "":
	case c.Key == "":
	case flag && c.Value == nil:
	default:
		return true
	}
	return false
}
