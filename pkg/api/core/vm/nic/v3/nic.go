package v3

import (
	"fmt"
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

func (h *NICHandler) GenerateOnlyMac(input []vm.VMNIC) ([]vm.VMNIC, error) {
	var usedMAC []string

	for idx, nicTmp := range input {
		if nicTmp.MAC == "" {
			mac, err := h.generateMac(usedMAC)
			if err != nil {
				return input, fmt.Errorf("MAC Address Generate Error: %s", err)
			}
			usedMAC = append(usedMAC, mac)
			input[idx].MAC = mac
		}
	}

	return input, nil
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
