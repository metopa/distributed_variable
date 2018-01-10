package common

import (
	"time"
)

//TODO Add synchronization
type Context struct {
	Name           string
	ServerAddr     PeerAddr
	LinkedPeers    [2]PeerAddr
	Leader         PeerAddr
	LeaderDistance [2]int
	KnownPeers     map[PeerAddr]PeerInfo
	SendNumRetries int
	SendRetryPause time.Duration
}

func NewContext(name string, sendNumRetries int, sendRetryPause time.Duration) *Context {
	return &Context{
		Name:           name,
		SendNumRetries: sendNumRetries,
		SendRetryPause: sendRetryPause,
		KnownPeers:     make(map[PeerAddr]PeerInfo)}
}

func (ctx *Context) AddNewPeer(name string, addr PeerAddr) {
	if _, ok := ctx.KnownPeers[addr]; !ok {
		ctx.KnownPeers[addr] = PeerInfo{Addr: addr, Name: name}
	}
}

func (ctx *Context) ResolvePeerName(addr PeerAddr) string {
	res := ""
	if addr == ctx.Leader {
		res += "[L]"
	}
	if info, ok := ctx.KnownPeers[addr]; ok {
		return info.Name + res
	}
	return string(addr) + res
}
