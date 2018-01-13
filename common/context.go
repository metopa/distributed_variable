package common

import (
	"math/rand"
	"sync"
	"time"
)

type Context struct {
	Name             string
	ServerAddr       PeerAddr
	PeerId           int
	LinkedPeers      [2]PeerAddr
	Leader           PeerAddr
	LeaderDistance   [2]int
	KnownPeers       map[PeerAddr]PeerInfo
	SendNumRetries   int
	SendRetryPause   time.Duration
	ChRoTimerDur     time.Duration
	state            State
	StateSync        sync.Mutex
	Sync             sync.Mutex
	StartedChRoTimer int32
	Clock            LamportClock
}

func NewContext(name string, sendNumRetries int, sendRetryPause, chRoTimerDur time.Duration) *Context {
	return &Context{
		Name:           name,
		SendNumRetries: sendNumRetries,
		SendRetryPause: sendRetryPause,
		KnownPeers:     make(map[PeerAddr]PeerInfo),
		PeerId:         rand.Int(),
		ChRoTimerDur:   chRoTimerDur}
}

func (ctx *Context) AddNewPeer(name string, addr PeerAddr) {
	if addr == ctx.ServerAddr {
		return
	}
	ctx.Sync.Lock()
	defer ctx.Sync.Unlock()
	if _, ok := ctx.KnownPeers[addr]; !ok {
		ctx.KnownPeers[addr] = PeerInfo{Addr: addr, Name: name}

		if len(ctx.LinkedPeers[0]) == 0 ||
			ctx.LinkedPeers[0] < addr && addr < ctx.ServerAddr ||
			ctx.ServerAddr < ctx.LinkedPeers[0] &&
				(ctx.LinkedPeers[0] < addr || addr < ctx.ServerAddr) {
			ctx.LinkedPeers[0] = addr
		}
		if len(ctx.LinkedPeers[1]) == 0 ||
			ctx.LinkedPeers[1] > addr && addr > ctx.ServerAddr ||
			ctx.ServerAddr > ctx.LinkedPeers[1] &&
				(ctx.LinkedPeers[1] > addr || addr > ctx.ServerAddr) {
			ctx.LinkedPeers[1] = addr
		}
	}
}

func (ctx *Context) ResolvePeerName(addr PeerAddr) string {
	ctx.Sync.Lock()
	defer ctx.Sync.Unlock()
	res := ""
	if addr == ctx.Leader {
		res += "[L]"
	}
	if info, ok := ctx.KnownPeers[addr]; ok {
		return info.Name + res
	}
	return string(addr) + res
}

func (ctx *Context) CASState(current, new State) bool {
	ctx.StateSync.Lock()
	defer ctx.StateSync.Unlock()
	if ctx.state == current {
		ctx.state = new
		ctx.state.Init()
		return true
	}
	return false
}

func (ctx *Context) GetState() State {
	return ctx.state
}

func (ctx *Context) SetState(new State) {
	ctx.state = new
	ctx.state.Init()
}
