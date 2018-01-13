package state

import (
	"fmt"
	"time"

	"github.com/metopa/distributed_variable/common"
	"github.com/metopa/distributed_variable/logger"
	"github.com/metopa/distributed_variable/net"
)

type LeaderState struct {
	DiscoveryState
	Value int
}

func (s *LeaderState) Init() {
	logger.Info("Current state: LEADER")
	go func(s *LeaderState) {
		net.BroadcastInRing(s.Ctx, common.NewSetLeaderCommand(s.Ctx.ServerAddr, 0))
		time.Sleep(time.Second / 4)
		s.EmitDistanceBroadcast()
	}(s)

}
func (s *LeaderState) NewPeer(sender common.PeerAddr, addr common.PeerAddr,
	name string, shouldReply bool) {
	s.Ctx.AddNewPeer(name, addr)
	logger.Info("Added new peer: %v(%v), linked peers: %v", name, addr, s.Ctx.LinkedPeers)
	if shouldReply {
		net.SendToDirectly(s.Ctx, addr,
			common.NewPeerInfoResponseCommand(s.Ctx.Name))
	}
	time.Sleep(time.Second / 4)
	s.EmitDistanceBroadcast()
}

func (s *LeaderState) ValueGetRequested(sender common.PeerAddr, source common.PeerAddr) {
	cmd := common.NewGetResponseCommand(s.Value)
	cmd.Destination = source
	net.ReplyInRing(s.Ctx, sender, cmd)
}

func (s *LeaderState) ValueSetRequested(sender common.PeerAddr,
	source common.PeerAddr, value int) {
	s.ActionSetValue(value)
	cmd := common.NewSetResponseCommand()
	cmd.Destination = source
	net.ReplyInRing(s.Ctx, sender, cmd)
}

func (s *LeaderState) LeaderChanged(sender common.PeerAddr, leader common.PeerAddr, value int) {}

func (s *LeaderState) PeerReported(reportedPeer common.PeerAddr) {
	logger.Warn("Peer %v removed", reportedPeer)
	s.Ctx.RemovePeer(reportedPeer)
	go net.SendToHi(s.Ctx, common.NewRemovePeerCommand(reportedPeer), true)
	go net.SendToLo(s.Ctx, common.NewRemovePeerCommand(reportedPeer), true)
	time.Sleep(time.Second / 2)
	s.EmitDistanceBroadcast()
}
func (s *LeaderState) PeerRemoved(sender common.PeerAddr, reportedPeer common.PeerAddr) {}

func (s *LeaderState) DistanceRequested(sender common.PeerAddr, source common.PeerAddr) {
	s.EmitDistanceBroadcast()
}

func (s *LeaderState) DistanceReceived(sender common.PeerAddr, distance int) {}

func (s *LeaderState) EmitDistanceBroadcast() {
	go net.SendToHi(s.Ctx, common.NewLeaderDistanceResponseCommand(0), true)
	net.SendToLo(s.Ctx, common.NewLeaderDistanceResponseCommand(0), true)
}

func (s *LeaderState) Name() string {
	return "Leader state"
}

func (s *LeaderState) ActionStartChRo() {
	fmt.Println("Peer has already joined the ring")
}

func (s *LeaderState) ActionSetValue(value int) {
	s.Value = value
	fmt.Printf("Value updated to %v\n", value)
}

func (s *LeaderState) ActionGetValue() {
	fmt.Printf("Value = %v\n", s.Value)
}
func (s *LeaderState) ActionReportPeer(addr common.PeerAddr) {
	s.PeerReported(addr)
}
func (s *LeaderState) ActionLeave() bool {
	if len(s.Ctx.KnownPeers) == 0 {
		logger.Info("No peers in the ring. Just leaving")
		return true
	}
	if len(s.Ctx.LinkedPeers[1]) == 0 {
		logger.Warn("Hi peer is unknown, the ring is in inconsistent state")
		return true
	}
	s.Ctx.Leader = s.Ctx.LinkedPeers[1]

	logger.Info("Transferring leadership to %v", s.Ctx.Leader)
	net.SendToHi(s.Ctx, common.NewSetLeaderCommand(s.Ctx.Leader, s.Value), false)
	time.Sleep(time.Second / 4)
	net.SendToRingLeader(s.Ctx, common.NewReportPeerCommand(s.Ctx.ServerAddr))
	return true
}
