package common

import (
	"fmt"
	"sync/atomic"
)

type LamportClock struct {
	value uint64
}

func (l LamportClock) String() string {
	return fmt.Sprintf("[%4d]", uint64(l.value))
}

func (l *LamportClock) SyncAfter(remoteClock LamportClock, delta uint64) {
	remote := remoteClock.value + delta
	for {
		old := l.value
		local := old + delta
		if remote > local {
			local = remote
		}
		if atomic.CompareAndSwapUint64(&(l.value), old, local) {
			break
		}
	}
}
