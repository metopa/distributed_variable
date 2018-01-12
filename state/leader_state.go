package state

import (
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
