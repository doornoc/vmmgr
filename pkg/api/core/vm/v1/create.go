package v1

import (
	"fmt"
	"github.com/vmmgr/controller/pkg/api/core/tool/config"
	"github.com/vmmgr/controller/pkg/api/core/tool/gen"
	"github.com/vmmgr/controller/pkg/api/core/tool/template"
	"github.com/vmmgr/controller/pkg/api/core/vm"
	nic "github.com/vmmgr/controller/pkg/api/core/vm/nic/v3"
	storage "github.com/vmmgr/controller/pkg/api/core/vm/storage/v3"
	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
	"log"
	"time"
)

func (b *Base) createVM(hostname string, data vm.VMInput) {
	log.Println("WebSocket VM Create " + "(" + hostname + ")")
	sshHost, err := config.CollectConfig(&hostname)
	if err != nil {
		b.Error("failed to config collect config: " + err.Error())
		return
	}

	conn, err := ConnectLibvirt(sshHost[0].HostName, sshHost[0].User)
	if err != nil {
		b.Error("failed to connect to qemu: " + err.Error())
		return
	}
	defer conn.Close()
	doms, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE | libvirt.CONNECT_LIST_DOMAINS_INACTIVE)
	if err != nil {
		b.Error(fmt.Sprintf("ListAllDomains error: %s", err.Error()))
		return
	}
	// 同じVM名がないか確認
	sameName := false
	for _, dom := range doms {
		vmName, _ := dom.GetName()
		if vmName == data.Name {
			sameName = true
			break
		}
	}
	if sameName {
		b.Error("failed to same vm name")
		return
	}

	vncPort, webSocketPort, err := gen.GenerateVNCAndWebSocketPort(doms)
	if err != nil {
		log.Printf("ListAllDomains error: %s", err)
		return
	}

	// Get Template
	tpl, err := template.Get(sshHost[0])
	if err != nil {
		log.Printf("Get Template error: %s", err)
		return
	}

	// addr
	addr := vm.Address{
		PCICount:  0,
		DiskCount: 0,
	}

	var domCfg *libvirtxml.Domain
	if data.CloudInit == nil {
		domCfg = &libvirtxml.Domain{
			Type: "kvm",
			Memory: &libvirtxml.DomainMemory{
				Value:    data.Memory,
				Unit:     "MB",
				DumpCore: "on",
			},
			VCPU:  &libvirtxml.DomainVCPU{Value: data.CPU},
			Name:  data.Name,
			Title: data.Name,
			Features: &libvirtxml.DomainFeatureList{
				ACPI: &libvirtxml.DomainFeature{},
				APIC: &libvirtxml.DomainFeatureAPIC{},
			},
			OS: &libvirtxml.DomainOS{
				BootDevices: []libvirtxml.DomainBootDevice{{Dev: data.Boot}},
				//Kernel:      "",
				//Initrd:  "/home/markus/workspace/worker-management/centos/kvm-centos.ks",
				//Cmdline: "ks=file:/home/markus/workspace/worker-management/centos/kvm-centos.ks method=http://repo02.agfa.be/CentOS/7/os/x86_64/",
				Type: &libvirtxml.DomainOSType{
					Arch:    template.GetArchStr(data.Arch),
					Machine: "pc", //kvm -machine help
					Type:    "hvm",
				},
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
							//Keymap:    h.VM.KeyMap,
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

		// storage
		sto := storage.NewStorageHandler(storage.StorageHandler{
			SSHHost:  sshHost[0],
			Template: tpl,
			Address:  &addr,
			Xml:      domCfg,
			Chan:     &vm.ClientBroadcast,
		})
		err = sto.Add(data.Name, data.Disk)
		if err != nil {
			b.Error("failed to storage apply: " + err.Error())
			return
		}

		// nic
		ni := nic.NewNICHandler(nic.NICHandler{
			Conn:     conn,
			Doms:     doms,
			SSHHost:  sshHost[0],
			Template: tpl,
			Address:  &addr,
			Xml:      domCfg,
			Chan:     &vm.ClientBroadcast,
		})
		err = ni.Add(data.NIC)
		if err != nil {
			b.Error("failed to nic apply: " + err.Error())
			return
		}
	}

	xml, err := domCfg.Marshal()
	if err != nil {
		b.Error("domCfg marshal error: " + err.Error())
		return
	}
	log.Println("xml", xml)
	dom, err := conn.DomainDefineXML(xml)
	if err != nil {
		b.Error("vm domainDefineXML error: " + err.Error())
		return
	}

	err = dom.Create()
	if err != nil {
		b.Error("vm create error: " + err.Error())
		return
	}
	vm.ClientBroadcast <- vm.WebSocketResult{
		Type:      b.Type,
		CreatedAt: time.Now(),
		VMDetail: []vm.VMDetail{{
			Node: hostname,
		}},
	}
}
