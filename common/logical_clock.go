package common

import (
	"fmt"
	"sync/atomic"
)

type LamportClock struct {
	Value uint64
}

func (l LamportClock) String() string {
	return fmt.Sprintf("[%4d]", uint64(l.Value))
}

func (l *LamportClock) SyncAfter(remoteClock LamportClock, delta uint64) {
	remote := remoteClock.Value + delta
	for {
		old := l.Value
		local := old + delta
		if remote > local {
			local = remote
		}
		if atomic.CompareAndSwapUint64(&(l.Value), old, local) {
			break
		}
	}
}
