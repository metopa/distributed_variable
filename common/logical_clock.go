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

func (l *LamportClock) SyncAfter(remoteClock LamportClock) {
	new := remoteClock.value + 1

	for old := l.value; old < new; old = l.value {
		atomic.CompareAndSwapUint64(&(l.value), old, new)
	}
}
