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

func (h *DiscoveryState) NewPeer(sender common.PeerAddr, addr common.PeerAddr,
	name string, shouldReply bool) {
	h.Ctx.AddNewPeer(name, addr)
	logger.Info("Added new peer: %v(%v)", name, addr)
	logger.Info("Linked peers: %v", h.Ctx.LinkedPeers)
	if shouldReply {
		net.SendToDirectly(h.Ctx, addr,
			common.NewPeerInfoResponseCommand(h.Ctx.Name, h.Ctx.Leader))
	}
}

func (h *DiscoveryState) LeaderChanged(sender common.PeerAddr, leader common.PeerAddr) {
	h.Ctx.Leader = leader
	ls := &LinkedState{*h}
	if h.Ctx.Server.CasCommandHandler(h, ls) {
		logger.Info("%v is ring leader", leader)
		ls.Start()
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
		ls := &LeaderState{DiscoveryState:*h}
		if h.Ctx.Server.CasCommandHandler(h, ls) {
			logger.Info("This peer is ring leader")
			h.Ctx.Leader = h.Ctx.ServerAddr
			ls.Start()
		}
	}
}
