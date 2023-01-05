package v1

import (
	"fmt"
	"github.com/vmmgr/controller/pkg/api/core/tool/gen"
	"github.com/vmmgr/controller/pkg/api/core/tool/template"
	"github.com/vmmgr/controller/pkg/api/core/vm"
	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
)

func initGenerateXml(input vm.VMInput, doms []libvirt.Domain) (*libvirtxml.Domain, error) {
	var domCfg *libvirtxml.Domain

	vncPort, webSocketPort, err := gen.GenerateVNCAndWebSocketPort(doms)
	if err != nil {
		return domCfg, fmt.Errorf("ListAllDomains error: %s", err)
	}

	domCfg = &libvirtxml.Domain{
		Type:  "kvm",
		Name:  input.Name,
		Title: input.Name,
		Features: &libvirtxml.DomainFeatureList{
			ACPI: &libvirtxml.DomainFeature{},
			APIC: &libvirtxml.DomainFeatureAPIC{},
		},
		OS: &libvirtxml.DomainOS{
			BootDevices: []libvirtxml.DomainBootDevice{{Dev: input.Boot}},
			//Kernel:      "",
			//Initrd:  "/home/markus/workspace/worker-management/centos/kvm-centos.ks",
			//Cmdline: "ks=file:/home/markus/workspace/worker-management/centos/kvm-centos.ks method=http://repo02.agfa.be/CentOS/7/os/x86_64/",
			Type: &libvirtxml.DomainOSType{
				Arch:    template.GetArchStr(input.Arch),
				Machine: "pc", //kvm -machine help
				Type:    "hvm",
			},
		},
		VCPU: &libvirtxml.DomainVCPU{
			Value: input.CPU,
		},
		Memory: &libvirtxml.DomainMemory{
			Value:    input.Memory,
			Unit:     "MB",
			DumpCore: "on",
		},
		Devices: &libvirtxml.DomainDeviceList{
			//Emulator: h.Node.Emulator,
			Inputs: []libvirtxml.DomainInput{
				{Type: "mouse", Bus: "ps2"},
				{Type: "keyboard", Bus: "ps2"},
			},
			Graphics: []libvirtxml.DomainGraphic{
				{
					VNC: &libvirtxml.DomainGraphicVNC{
						Port:      vncPort,
						WebSocket: webSocketPort,
						//Keymap:    h.Store.KeyMap,
						Listen: "0.0.0.0",
					},
				},
			},
			Videos: []libvirtxml.DomainVideo{
				{
					Model: libvirtxml.DomainVideoModel{
						Type:    "qxl",
						Heads:   1,
						Ram:     65536,
						VRam:    65536,
						VGAMem:  16384,
						Primary: "yes",
					},
					Alias: &libvirtxml.DomainAlias{Name: "video0"},
				},
			},
		},
	}

	return domCfg, nil
}
