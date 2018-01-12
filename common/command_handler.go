package common

type CommandHandler interface {
	NewPeer(sender PeerAddr, addr PeerAddr, name string, shouldReply bool)
	LeaderChanged(sender PeerAddr, leader PeerAddr)
	LinkedPeersChanged(sender PeerAddr, loPeer PeerAddr, hiPeer PeerAddr)
	PeerReported(sender PeerAddr, reportedPeer PeerAddr)
	DistanceRequested(sender PeerAddr, source PeerAddr)
	DistanceReceived(sender PeerAddr, distance int)
	SyncPeers(sender PeerAddr, values []string)
	ChRoIdReceived(sender PeerAddr, id int)
	RingJoinRequested(sender PeerAddr, source PeerAddr)
	RingLeaveAnnounced(sender PeerAddr, source PeerAddr)
	ValueGetRequested(sender PeerAddr, source PeerAddr)
	GotValue(sender PeerAddr, value int)
	ValueSetRequested(sender PeerAddr, source PeerAddr, value int)
	ValueSetConfirmed(sender PeerAddr)
}

type HandlerOwner interface {
	SetCommandHandler(handler CommandHandler)
	CasCommandHandler(current CommandHandler, new CommandHandler) bool
	GetCommandHandler() CommandHandler
}
