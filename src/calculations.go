package src

import (
	"sync"
	"sync/atomic"

	"github.com/streamonkey/size"
)

type memorylimit struct {
	configured   uint64
	remainingSpace uint64
}


type sessionPool struct {
	userSessions map[string]
} 
type activeUserSession struct {
	userid string
	currentSessionPool uint8
	
}


type ActiveSessionsRegistry struct {
    mu          sync.RWMutex
    maxSessions uint8
    users       map[string]*ActiveUserSession
}

type ActiveUserSession struct {
    SessionIDs []string
}


var osavailableMemory atomic.Uint64

func CalculateInputBytes(value any, ch chan<- bool) {
	bytevalue := size.Of(value)

}

func mbSizeToUINT(value float64, c chan<- uint64, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	c <- uint64(value * mbtouintsize)
	close(c)
}

func compareConfigOsMem(osmem uint64, configmem uint64) bool {
	osmem -= osmem * Memory_cutoff / 100
	if osmem > configmem {
		return true
	}
	return false
}

func sessionPoolChecker(userdto *userDTO, c chan <- error, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	defer close(c)

	switch {
	case userdto.isNew == false:
        var u *UserManager
		u.Mu.RLock()
		user,exist := u.Users[userdto.userid]
		u.Mu.RUnlock()
		if exist && user.isActive == true {
            var pool *activeUserSessions
			pool.mu.RLock()
			userinpool,exist := pool.userid[userdto.userid]
			pool.mu.RUnlock()
		}
		c <- errUser
		return

	case userdto.isNew == true:

	}
}

