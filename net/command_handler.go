package net

import (
	"github.com/metopa/distributed_variable/common"
	"github.com/metopa/distributed_variable/logger"
)

type CommandHandler interface {
	NewPeer(sender common.PeerAddr, addr common.PeerAddr, name string, shouldReply bool)
	LeaderChanged(sender common.PeerAddr, leader common.PeerAddr)
	LinkedPeersChanged(sender common.PeerAddr, loPeer common.PeerAddr, hiPeer common.PeerAddr)
	PeerReported(sender common.PeerAddr, reportedPeer common.PeerAddr)
	DistanceRequested(sender common.PeerAddr, source common.PeerAddr)
	DistanceReceived(sender common.PeerAddr, distance int)
	RingJoinRequested(sender common.PeerAddr, source common.PeerAddr)
	RingLeaveAnnounced(sender common.PeerAddr, source common.PeerAddr)
	ValueGetRequested(sender common.PeerAddr, source common.PeerAddr)
	GotValue(sender common.PeerAddr, value int)
	ValueSetRequested(sender common.PeerAddr, source common.PeerAddr, value int)
	ValueSetConfirmed(sender common.PeerAddr)
}

type InitialCommandHandler struct {
	Ctx *common.Context
}

func (h *InitialCommandHandler) NewPeer(sender common.PeerAddr, addr common.PeerAddr,
	name string, shouldReply bool) {
	h.Ctx.AddNewPeer(name, addr)
	logger.Info("Added new peer: %v(%v)", name, addr)
	if shouldReply {
		SendToDirectly(h.Ctx, addr, NewPeerInfoResponseCommand(name, h.Ctx.Leader))
	}
}
func (h *InitialCommandHandler) LeaderChanged(sender common.PeerAddr, leader common.PeerAddr) {
	logger.Warn("LeaderChanged")
}
func (h *InitialCommandHandler) LinkedPeersChanged(sender common.PeerAddr,
	loPeer common.PeerAddr, hiPeer common.PeerAddr) {
	logger.Warn("LinkedPeersChanged")
}
func (h *InitialCommandHandler) PeerReported(sender common.PeerAddr, reportedPeer common.PeerAddr) {
	logger.Warn("PeerReported")
}
func (h *InitialCommandHandler) DistanceRequested(sender common.PeerAddr, source common.PeerAddr) {
	logger.Warn("DistanceRequested")
}
func (h *InitialCommandHandler) DistanceReceived(sender common.PeerAddr, distance int) {
	logger.Warn("DistanceReceived")
}
func (h *InitialCommandHandler) RingJoinRequested(sender common.PeerAddr, source common.PeerAddr) {
	logger.Warn("RingJoinRequested")
}
func (h *InitialCommandHandler) RingLeaveAnnounced(sender common.PeerAddr, source common.PeerAddr) {
	logger.Warn("RingLeaveAnnounced")
}
func (h *InitialCommandHandler) ValueGetRequested(sender common.PeerAddr, source common.PeerAddr) {
	logger.Warn("ValueGetRequested")
}
func (h *InitialCommandHandler) GotValue(sender common.PeerAddr, value int) {
	logger.Warn("GotValue")
}
func (h *InitialCommandHandler) ValueSetRequested(sender common.PeerAddr,
	source common.PeerAddr, value int) {
	logger.Warn("ValueSetRequested")
}
func (h *InitialCommandHandler) ValueSetConfirmed(sender common.PeerAddr) {
	logger.Warn("ValueSetConfirmed")
}
