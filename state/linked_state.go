package state

import (
	"fmt"
	"time"

	"github.com/metopa/distributed_variable/common"
	"github.com/metopa/distributed_variable/logger"
	"github.com/metopa/distributed_variable/net"
)

type LinkedState struct {
	DiscoveryState
}

func (s *LinkedState) Init() {
	logger.Info("Current state: LINKED")
}
func (s *LinkedState) GotValue(sender common.PeerAddr, value int) {
	fmt.Printf("Value = %v\n", value)
}
func (s *LinkedState) ValueSetConfirmed(sender common.PeerAddr) {
	fmt.Printf("Value is updated\n")
}

func (s *LinkedState) LeaderChanged(sender common.PeerAddr, leader common.PeerAddr, value int) {
	logger.Info("New leader: %v, prev: %v", leader, s.Ctx.Leader)
	s.Ctx.Leader = leader
	if leader == s.Ctx.ServerAddr {
		ls := &LeaderState{DiscoveryState: DiscoveryState{Ctx: s.Ctx}, Value: value}
		s.Ctx.CASState(s, ls)
	}
}

func (s *LinkedState) DistanceReceived(sender common.PeerAddr, distance int) {
	distance++
	updated := false
	prevDistances := s.Ctx.LeaderDistance
	if sender == s.Ctx.LinkedPeers[0] {
		s.Ctx.LeaderDistance[0] = distance
		go net.SendToHi(s.Ctx, common.NewLeaderDistanceResponseCommand(distance), true)
		updated = true
	}
	if sender == s.Ctx.LinkedPeers[1] {
		s.Ctx.LeaderDistance[1] = distance
		go net.SendToLo(s.Ctx, common.NewLeaderDistanceResponseCommand(distance), true)
		updated = true
	}

	if updated && prevDistances != s.Ctx.LeaderDistance {
		logger.Info("Leader distance updated: %v", s.Ctx.LeaderDistance)
	} else {
		logger.Warn("Distance received from %v, but it's not in linked peers", sender)
	}
}

func (s *LinkedState) PeerRemoved(sender common.PeerAddr, removedPeer common.PeerAddr) {
	logger.Warn("Peer %v removed", removedPeer)
	fromLo := sender == s.Ctx.LinkedPeers[0]
	s.Ctx.RemovePeer(removedPeer)

	if fromLo {
		net.SendToHi(s.Ctx, common.NewRemovePeerCommand(removedPeer), true)
	} else {
		net.SendToLo(s.Ctx, common.NewRemovePeerCommand(removedPeer), true)
	}
}

func (s *LinkedState) Name() string {
	return "Linked state state"
}

func (s *LinkedState) ActionStartChRo() {
	fmt.Println("Peer has already joined the ring")
}
func (s *LinkedState) ActionSetValue(value int) {
	s.RequestDistancesIfMissing()
	net.SendToRingLeader(s.Ctx, common.NewSetRequestCommand(value))
}
func (s *LinkedState) ActionGetValue() {
	s.RequestDistancesIfMissing()
	net.SendToRingLeader(s.Ctx, common.NewGetRequestCommand())
}

func (s *LinkedState) RequestDistancesIfMissing() {
	if s.Ctx.LeaderDistance[0] == -1 || s.Ctx.LeaderDistance[1] == -1 {
		logger.Info("Requested leader distance broadcast")
		go net.SendToRingLeader(s.Ctx, common.NewLeaderDistanceRequestCommand())
	}
}
func (s *LinkedState) ActionReportPeer(addr common.PeerAddr) {
	if addr == s.Ctx.Leader {
		logger.Warn("Wanted to report peer %v, but it's leader", addr)
		return
	}
	logger.Warn("Reporting peer %v", addr)
	cmd := common.NewReportPeerCommand(addr)
	cmd.Destination = s.Ctx.Leader
	if addr == s.Ctx.LinkedPeers[0] {
		net.SendToHi(s.Ctx, cmd, false)
	} else {
		net.SendToLo(s.Ctx, cmd, false)
	}
}

func (s *LinkedState) ActionLeave() bool {
	net.SendToRingLeader(s.Ctx, common.NewReportPeerCommand(s.Ctx.ServerAddr))
	time.Sleep(time.Second / 10)
	return true
}