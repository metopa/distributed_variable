package state

import (
	"fmt"

	"github.com/metopa/distributed_variable/common"
	"github.com/metopa/distributed_variable/net"
)

type LeaderState struct {
	DiscoveryState
	Value int
}

func (s *LeaderState) Start() {
	net.BroadcastInRing(s.Ctx, common.NewSetLeaderCommand(s.Ctx.Leader))
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
