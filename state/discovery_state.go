package state

import (
	"time"

	"github.com/metopa/distributed_variable/common"
	"github.com/metopa/distributed_variable/logger"
	"github.com/metopa/distributed_variable/net"
)

type DiscoveryState struct {
	NullState
	Ctx        *common.Context
	alivePeers []common.PeerAddr
}

func (s *DiscoveryState) Init() {
	logger.Info("Current state: DISCOVERY")
}

func (h *DiscoveryState) NewPeer(sender common.PeerAddr, addr common.PeerAddr,
	name string, shouldReply bool) {
	h.Ctx.AddNewPeer(name, addr)
	logger.Info("Added new peer: %v(%v), linked peers: %v", name, addr, h.Ctx.LinkedPeers)
	if shouldReply {
		net.SendToDirectly(h.Ctx, addr,
			common.NewPeerInfoResponseCommand(h.Ctx.Name))
	}
}

func (h *DiscoveryState) LeaderChanged(sender common.PeerAddr, leader common.PeerAddr, value int) {
	if leader != h.Ctx.Leader {
		logger.Info("New leader: %v", leader)
		h.Ctx.Leader = leader
		ls := &LinkedState{*h}
		h.Ctx.CASState(h, ls)
	}
}

func (h *DiscoveryState) SyncPeers(sender common.PeerAddr, values []string) {
	peers := make(map[common.PeerAddr]common.PeerInfo)
	for i := 0; i < len(values); i += 2 {
		peers[common.PeerAddr(values[i])] = common.PeerInfo{
			Name: values[i+1], Addr: common.PeerAddr(values[i])}
	}
	h.Ctx.SetKnownPeers(peers)
	logger.Info("Ring peers: %v", h.Ctx.KnownPeers)
	logger.Info("Linked peers: %v", h.Ctx.LinkedPeers)
	net.StartChRoTimer(h.Ctx)
}

func (s *DiscoveryState) PeerRemoved(sender common.PeerAddr, removedPeer common.PeerAddr, direction int) {
	logger.Warn("Peer %v removed", removedPeer)

	s.Ctx.RemovePeer(removedPeer)
}

func (h *DiscoveryState) Ping(sender common.PeerAddr, source common.PeerAddr) {
	net.SendToDirectly(h.Ctx, source, common.NewPongCmd())
}

func (h *DiscoveryState) Pong(sender common.PeerAddr, source common.PeerAddr) {
	h.alivePeers = append(h.alivePeers, source)
}

func (h *DiscoveryState) ChRoIdReceived(sender common.PeerAddr, id int) {
	if id > h.Ctx.PeerId {
		net.SendToHi(h.Ctx, common.NewChangRobertIdCmd(id), false)
	} else if id == h.Ctx.PeerId {
		h.Ctx.Leader = h.Ctx.ServerAddr
		ls := &LeaderState{DiscoveryState: *h}
		h.Ctx.CASState(h, ls)
	}
}

func (h *DiscoveryState) ActionStartChRo() {
	if len(h.Ctx.KnownPeers) == 0 {
		logger.Warn("No other peers connected, can't build peer ring")
		return
	}
	if h.Ctx.StartedChRoTimer == 1 {
		logger.Info("ChRo has already been started")
		return
	}
	h.alivePeers = nil
	for addr := range h.Ctx.KnownPeers {
		go net.SendToDirectly(h.Ctx, addr, common.NewPingCmd())
	}

	time.Sleep(h.Ctx.ChRoTimerDur)

	alivePeers := make(map[common.PeerAddr]common.PeerInfo)
	h.Ctx.Sync.Lock()
	for _, addr := range h.alivePeers {
		info, ok := h.Ctx.KnownPeers[addr]
		if ok {
			alivePeers[addr] = info
		}
	}
	h.Ctx.Sync.Unlock()
	alivePeers[h.Ctx.ServerAddr] = common.PeerInfo{Name: h.Ctx.Name, Addr: h.Ctx.ServerAddr}
	h.Ctx.SetKnownPeers(alivePeers)

	if len(alivePeers) == 0 {
		logger.Warn("No other peers connected, can't build peer ring")
		return
	}
	logger.Info("Ring peers: %v", h.alivePeers)
	cmd := common.NewSyncPeersCmd(alivePeers)
	h.Ctx.Sync.Lock()
	for addr := range h.Ctx.KnownPeers {
		go net.SendToDirectly(h.Ctx, addr, cmd)
	}
	h.Ctx.Sync.Unlock()
	net.StartChRoTimer(h.Ctx)
}

func (s *DiscoveryState) ActionLeave() bool {
	cmd := common.NewRemovePeerCommand(s.Ctx.ServerAddr, 0)
	s.Ctx.Sync.Lock()
	for addr, _ := range s.Ctx.KnownPeers {
		go net.SendToDirectly(s.Ctx, addr, cmd)
	}
	s.Ctx.Sync.Unlock()
	time.Sleep(time.Second / 10)
	return true
}

func (h *DiscoveryState) ActionSync() {
	logger.Info("Peer is not in the ring, no synchronization needed")
}

func (h *DiscoveryState) Name() string {
	return "Discovery state"
}
