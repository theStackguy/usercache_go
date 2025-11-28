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

var osavailableMemory atomic.Uint64

func CalculateInputBytes(value any, ch chan<- bool) {
	bytevalue := size.Of(value)

}

func mbSizeToUINT(value float64, c chan<- uint64, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	c <- uint64(value * MBTOUINTSIZE)
	close(c)
}

func compareConfigOsMem(osmem uint64, configmem uint64) bool {
	osmem -= osmem * MEMORY_CUTOFF / 100
	if osmem > configmem {
		return true
	}
	return false
}
