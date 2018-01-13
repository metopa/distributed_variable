package common

type ActionHandler interface {
	ActionSetValue(value int)
	ActionGetValue()
	ActionStartChRo()
	ActionLeave() bool
	ActionDisconnect()
	ActionSync()
	ActionReconnect()
	ActionReportPeer(addr PeerAddr)
}

type CommandHandler interface {
	NewPeer(sender PeerAddr, addr PeerAddr, name string, shouldReply bool)
	LeaderChanged(sender PeerAddr, leader PeerAddr, value int)

	PeerReported(reportedPeer PeerAddr)
	PeerRemoved(sender PeerAddr, reportedPeer PeerAddr)
	DistanceRequested(sender PeerAddr, source PeerAddr)
	DistanceReceived(sender PeerAddr, distance int)
	SyncPeers(sender PeerAddr, values []string)
	Ping(sender PeerAddr, source PeerAddr)
	Pong(sender PeerAddr, source PeerAddr)
	ChRoIdReceived(sender PeerAddr, id int)

	ValueGetRequested(sender PeerAddr, source PeerAddr)
	GotValue(sender PeerAddr, value int)
	ValueSetRequested(sender PeerAddr, source PeerAddr, value int)
	ValueSetConfirmed(sender PeerAddr)
}

type State interface {
	CommandHandler
	ActionHandler
	Name() string
	Init()
}
