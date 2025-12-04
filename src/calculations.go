package src

import (
	"sync"
	"sync/atomic"

	"github.com/streamonkey/size"
)

type memorylimit struct {
	configured     uint64
	remainingSpace uint64
}

type activeSessionsRegistry struct {
	mu          sync.RWMutex
	maxSessions uint8
	users       map[string]*activeUserSession
}

type activeUserSession struct {
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

func sessionPoolChecker(userdto *userDTO, c chan<- error, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	defer close(c)

	var u *UserManager
	u.Mu.RLock()
	user, exist := u.Users[userdto.userid]
	u.Mu.RUnlock()

	if exist && user.isActive == true {
		var pool *activeSessionsRegistry
		pool.mu.RLock()
		userinpool, userinpoolexist := pool.users[userdto.userid]
		pool.mu.RUnlock()

		if !userdto.isNew {
			if userinpoolexist {
				if len(userinpool.SessionIDs) < Allowed_Sessions {
					pool.mu.Lock()
					userinpool.SessionIDs = append(userinpool.SessionIDs, user.CurrentSessionId)
					pool.mu.Unlock()
					c <- nil
					return
				}
				c <- errSessionLimit
				return
			}
			newRegistryAssigner(userdto, user, pool)
			c <- nil
			return
		} else if userdto.isNew {
			newRegistryAssigner(userdto, user, pool)
			c <- nil
			return
		}
		c <- errUserDto
		return
	}
	c <- errUser
	return

}

func newRegistryAssigner(userdto *userDTO, user *User, pool *activeSessionsRegistry) {
	newActiveUserSession := &activeUserSession{
		SessionIDs: []string{user.CurrentSessionId},
	}
	pool.mu.Lock()
	pool.users[userdto.userid] = newActiveUserSession
	pool.mu.Unlock()
}
