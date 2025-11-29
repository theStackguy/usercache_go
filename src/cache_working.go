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
	Id       string
	Mu       sync.RWMutex
	Sessions map[string]*Session
	memory   map[string]*memorylimit
}

type Session struct {
	SessionId     string
	SessionToken  string
	RefreshToken  string
	IsActive      bool
	SessionExpiry time.Time
	RefreshExpiry time.Time
	LastAccessed  time.Time
	Cache         map[string]CacheItem
	Mu            sync.RWMutex
	Err           error
}

type CacheItem struct {
	Value        any
	ExpiryTime   time.Time
	LastAccessed time.Time
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
	var userId string = tokens[1]
	var sessionuser Session
	sessionToken, refreshToken, err := generateSessionRefreshToken(sessionTokenExpiryTime, refreshTokenExpiryTime)
	if err != nil {
		return nil, errTokenGen
	}
	// sessionuser := &Session{
	// 	SessionId:     tokens[ZERO],
	// 	SessionToken:  sessionToken,
	// 	RefreshToken:  refreshToken,
	// 	Expiry:        time.Now().Add(sessionTokenExpiryTime),
	// 	RefreshExpiry: time.Now().Add(refreshTokenExpiryTime),
	// 	LastAccessed:  time.Now(),
	// 	IsActive:      true,
	// 	Cache:         make(map[string]CacheItem),
	// }

	wg.Wait()

	if osmemerror := <-osmemorychannel; osmemerror != nil {
		return nil, osmemerror
	}
	if !compareConfigOsMem(osavailableMemory.Load(), <-memorychannel) {
		return nil, errMemExceeded
	}

	memory := &memorylimit{
		configured:     <-memorychannel,
		remainingSpace: <-memorychannel,
	}

	um.Mu.Lock()
	um.Users[userId] = &User{
		Id:       userId,
		Sessions: map[string]*Session{tokens[0]: sessionuser},
		memory:   map[string]*memorylimit{userId: memory},
	}
	um.Mu.Unlock()
	return um.Users[userId], nil
}

func (um *UserManager) AddSessionToUser(userId string, sessionTokenExpiryTime time.Duration, refreshTokenExpiryTime time.Duration) (*Session, error) {
	um.Mu.RLock()
	user, exists := um.Users[userId]
	um.Mu.RUnlock()
	if !exists {
		return nil, errUser
	}
	sessionId, err := newString()
	if err != nil {
		return nil, errGuid
	}
	sessionToken, refreshToken, err := generateSessionRefreshToken(sessionTokenExpiryTime, refreshTokenExpiryTime)
	if err != nil {
		return nil, errTokenGen
	}
	session := &Session{
		SessionId:     sessionId,
		SessionToken:  sessionToken,
		RefreshToken:  refreshToken,
		SessionExpiry:        time.Now().Add(sessionTokenExpiryTime),
		RefreshExpiry: time.Now().Add(refreshTokenExpiryTime),
		LastAccessed:  time.Now(),
		Cache:         make(map[string]CacheItem),
	}
	user.Mu.Lock()
	user.Sessions[sessionId] = session
	user.Mu.Unlock()
	return session, nil
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
