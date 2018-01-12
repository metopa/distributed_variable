package console

import (
	"github.com/metopa/distributed_variable/common"
	"github.com/metopa/distributed_variable/net"
)

type DefaultActionHandler struct {
	Ctx *common.Context
}

func (h *DefaultActionHandler) SetValue(value int) {

}
func (h *DefaultActionHandler) GetValue() {

}
func (h *DefaultActionHandler) StartChRo() {
	//TODO Lock peer list
	cmd := common.NewSyncPeersCmd(h.Ctx)
	h.Ctx.Sync.Lock()
	for addr := range h.Ctx.KnownPeers {
		net.SendToDirectly(h.Ctx, addr, cmd)
	}
	h.Ctx.Sync.Unlock()
	net.StartChRoTimer(h.Ctx)
}

func (h *DefaultActionHandler) Leave() {

}
func (h *DefaultActionHandler) Disconnect() {

}
func (h *DefaultActionHandler) Reconnect() {

}
