package state

import (
	"fmt"

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

func (s *LinkedState) LeaderChanged(sender common.PeerAddr, leader common.PeerAddr) {
	logger.Info("New leader: %v, prev: %v", leader, s.Ctx.Leader)
	s.Ctx.Leader = leader
	//TODO Check we're not the leader
}

func (s *LinkedState) DistanceReceived(sender common.PeerAddr, distance int) {
	distance++
	if sender == s.Ctx.LinkedPeers[0] {
		s.Ctx.LeaderDistance[0] = distance
		net.SendToHi(s.Ctx, common.NewLeaderDistanceResponseCommand(distance))
		logger.Info("Leader distance updated: %v", s.Ctx.LeaderDistance)
	} else if sender == s.Ctx.LinkedPeers[1] {
		s.Ctx.LeaderDistance[1] = distance
		net.SendToLo(s.Ctx, common.NewLeaderDistanceResponseCommand(distance))
		logger.Info("Leader distance updated: %v", s.Ctx.LeaderDistance)
	} else {
		logger.Warn("Distance received from %v, but it's not in linked peers", sender)
	}
}


func (s *LinkedState) Name() string {
	return "Linked state state"
}

func (s *LinkedState) ActionStartChRo() {
	fmt.Println("Peer has already joined the ring")
}
func (s *LinkedState)  ActionSetValue(value int) {
	fmt.Println("Set requested")
	s.RequestDistancesIfMissing()
	net.SendToRingLeader(s.Ctx, common.NewSetRequestCommand(value))
}
func (s *LinkedState)  ActionGetValue() {
	fmt.Println("Get requested")
	s.RequestDistancesIfMissing()
	net.SendToRingLeader(s.Ctx, common.NewGetRequestCommand())
}

func (s *LinkedState) RequestDistancesIfMissing() {
	if s.Ctx.LeaderDistance[0] == -1 || s.Ctx.LeaderDistance[1] == -1 {
		logger.Info("Requested leader distance broadcast")
		go net.SendToRingLeader(s.Ctx, common.NewLeaderDistanceRequestCommand())
	}
}
