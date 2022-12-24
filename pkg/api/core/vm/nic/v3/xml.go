package v3

import (
	"fmt"
	"github.com/vmmgr/controller/pkg/api/core/vm"
	"github.com/vmmgr/controller/pkg/api/core/vm/nic"
	"libvirt.org/go/libvirtxml"
)

func (h *NICHandler) xmlGenerate(input []vm.VMNIC) error {
	var nics []libvirtxml.DomainInterface

	var usedMAC []string
	var count = 0

	for _, nicTmp := range input {
		if nicTmp.MAC == "" {
			count++
		}
	}

	for _, nicTmp := range input {
		if nicTmp.MAC == "" {
			mac, err := h.generateMac(usedMAC)
			if err != nil {
				return fmt.Errorf("MAC Address Generate Error: %s", err)
			}
			usedMAC = append(usedMAC, mac)
			nicTmp.MAC = mac
		}

		h.Address.PCICount++

		nics = append(nics, *generateTemplate(nic.GenerateNICXml{
			NIC:           convertNIC(nicTmp),
			AddressNumber: h.Address.PCICount,
		}))
	}

	h.Xml.Devices.Interfaces = nics

	return nil
}

func generateTemplate(xmlStruct nic.GenerateNICXml) *libvirtxml.DomainInterface {
	//デフォルトはブートディスク(VirtIO)

	domNIC := libvirtxml.DomainInterface{}

	// Bridge
	if xmlStruct.NIC.Type == 0 {
		// defaultでもいけるかもしれない（要確認必要）
		domNIC.Source = &libvirtxml.DomainInterfaceSource{
			Bridge: &libvirtxml.DomainInterfaceSourceBridge{
				Bridge: xmlStruct.NIC.Device,
			},
		}
		// NAT
	} else if xmlStruct.NIC.Type == 1 {
		// defaultでもいけるかもしれない（要確認必要）
		domNIC.Source = &libvirtxml.DomainInterfaceSource{
			Network: &libvirtxml.DomainInterfaceSourceNetwork{
				Network: xmlStruct.NIC.Device,
			},
		}

		// macvtap
	} else if xmlStruct.NIC.Type == 2 {
		domNIC.Source = &libvirtxml.DomainInterfaceSource{
			Direct: &libvirtxml.DomainInterfaceSourceDirect{
				Dev:  xmlStruct.NIC.Device,
				Mode: nic.GetModeName(xmlStruct.NIC.Mode),
			},
		}
	}

	//Driver
	domNIC.Model = &libvirtxml.DomainInterfaceModel{
		Type: nic.GetDriverName(xmlStruct.NIC.Driver),
	}
	//MAC
	domNIC.MAC = &libvirtxml.DomainInterfaceMAC{
		Address: xmlStruct.NIC.MAC,
	}
	//PCI Address
	domNIC.Address = &libvirtxml.DomainAddress{
		PCI: &libvirtxml.DomainAddressPCI{
			Domain:   &[]uint{0}[0],
			Bus:      &[]uint{1}[0],
			Slot:     &[]uint{xmlStruct.AddressNumber}[0],
			Function: &[]uint{0}[0],
		},
	}

	return &domNIC
}
