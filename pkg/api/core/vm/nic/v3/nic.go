package v3

import (
	"github.com/vmmgr/controller/pkg/api/core/tool/config"
	"github.com/vmmgr/controller/pkg/api/core/tool/template"
	"github.com/vmmgr/controller/pkg/api/core/vm"
	"github.com/vmmgr/controller/pkg/api/core/vm/nic"
	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
)

type NICHandler struct {
	Conn     *libvirt.Connect
	Doms     []libvirt.Domain
	SSHHost  config.SSHHost
	Template template.Template
	Address  *vm.Address
	Xml      *libvirtxml.Domain
	Chan     *chan vm.WebSocketResult
}

func NewNICHandler(handler NICHandler) *NICHandler {
	return &handler
}

func (h *NICHandler) Add(input []vm.VMNIC) error {
	err := h.xmlGenerate(input)
	if err != nil {
		return err
	}

	return nil
}

func convertNIC(dev vm.VMNIC) nic.NIC {
	return nic.NIC{
		Type:   dev.Type,
		Driver: dev.Driver,
		Mode:   dev.Mode,
		MAC:    dev.MAC,
		Device: dev.Device,
	}
}
