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
		net.BroadcastInRing(s.Ctx, common.NewSetLeaderCommand(s.Ctx.Leader))
		time.Sleep(time.Second / 4)
		s.EmitDistanceBroadcast()
	}(s)

}
func (s *LeaderState) NewPeer(sender common.PeerAddr, addr common.PeerAddr,
	name string, shouldReply bool) {
	s.Ctx.AddNewPeer(name, addr)
	logger.Info("Added new peer: %v(%v)", name, addr)
	logger.Info("Linked peers: %v", s.Ctx.LinkedPeers)
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

func (s *LeaderState) LeaderChanged(sender common.PeerAddr, leader common.PeerAddr) {
	if leader != s.Ctx.Leader {
		logger.Warn("Leader transition is unsupported", leader, s.Ctx.Leader)
		//TODO Set Linked state
		//TODO Send current value
	}
}

func (s *LeaderState)PeerReported(reportedPeer common.PeerAddr) {
	logger.Warn("Peer %v removed", reportedPeer)
	s.Ctx.RemovePeer(reportedPeer)
	logger.Info("Linked peers after remove: %v", s.Ctx.LinkedPeers) //TODO Remove
	go net.SendToHi(s.Ctx, common.NewRemovePeerCommand(reportedPeer), true)
	go net.SendToLo(s.Ctx, common.NewRemovePeerCommand(reportedPeer), true)
	time.Sleep(time.Second)
	s.EmitDistanceBroadcast()
}
func (s *LeaderState)PeerRemoved(sender common.PeerAddr, reportedPeer common.PeerAddr) {}

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
	logger.Warn("Leader transfer is not supported")
	return false
}
