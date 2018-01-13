package state

import (
	"github.com/metopa/distributed_variable/common"
	"github.com/metopa/distributed_variable/logger"
	"github.com/metopa/distributed_variable/net"
)

type DiscoveryState struct {
	NullState
	Ctx *common.Context
}

func (s *DiscoveryState) Init() {
	logger.Info("Current state: DISCOVERY")
}

func (h *DiscoveryState) NewPeer(sender common.PeerAddr, addr common.PeerAddr,
	name string, shouldReply bool) {
	h.Ctx.AddNewPeer(name, addr)
	logger.Info("Added new peer: %v(%v)", name, addr)
	logger.Info("Linked peers: %v", h.Ctx.LinkedPeers)
	if shouldReply {
		net.SendToDirectly(h.Ctx, addr,
			common.NewPeerInfoResponseCommand(h.Ctx.Name))
	}
}

func (h *DiscoveryState) LeaderChanged(sender common.PeerAddr, leader common.PeerAddr) {
	if leader != h.Ctx.Leader {
		logger.Info("New leader: %v, prev: %v", leader, h.Ctx.Leader)
		h.Ctx.Leader = leader
		ls := &LinkedState{*h}
		h.Ctx.CASState(h, ls)
	}
}

func (h *DiscoveryState) SyncPeers(sender common.PeerAddr, values []string) {
	//TODO Lock peers?
	for i := 0; i < len(values); i += 2 {
		h.Ctx.AddNewPeer(values[i+1], common.PeerAddr(values[i]))
	}
	net.StartChRoTimer(h.Ctx)
}

func (h *DiscoveryState) ChRoIdReceived(sender common.PeerAddr, id int) {
	if id > h.Ctx.PeerId {
		net.SendToHi(h.Ctx, common.NewChangRobertIdCmd(id))
	} else if id == h.Ctx.PeerId {
		ls := &LeaderState{DiscoveryState: *h}
		if h.Ctx.CASState(h, ls) {
			logger.Info("This peer is ring leader")
			h.Ctx.Leader = h.Ctx.ServerAddr
		}
	}
}

func (h *DiscoveryState) ActionStartChRo() {
	//TODO Lock peer list
	if len(h.Ctx.KnownPeers) == 0 {
		logger.Warn("No other peers connected, can't build peer ring")
		return
	}
	if h.Ctx.StartedChRoTimer == 1 {
		logger.Info("ChRo has already been started")
		return
	}
	cmd := common.NewSyncPeersCmd(h.Ctx)
	h.Ctx.Sync.Lock()
	for addr := range h.Ctx.KnownPeers {
		net.SendToDirectly(h.Ctx, addr, cmd)
	}
	h.Ctx.Sync.Unlock()
	net.StartChRoTimer(h.Ctx)
}

func (h *DiscoveryState) Name() string {
	return "Discovery state"
}
