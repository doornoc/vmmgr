package v1

import (
	"encoding/xml"
	"github.com/doornoc/vmmgr/pkg/api/core/tool/config"
	"github.com/doornoc/vmmgr/pkg/api/core/vm"
	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
	"log"
	"time"
)

func (b *Base) getVM(hostname, uuid string) {
	log.Println("WebSocket VM Get " + uuid + "(" + hostname + ")")
	dom, err := getOneVM(hostname, uuid)
	if err != nil {
		b.Error("failed to getting state: " + err.Error())
		return
	}

	t := libvirtxml.Domain{}
	stat, _, _ := dom.GetState()
	xmlString, _ := dom.GetXMLDesc(libvirt.DOMAIN_XML_SECURE)
	xml.Unmarshal([]byte(xmlString), &t)

	vm.ClientBroadcast <- vm.WebSocketResult{
		UUID:      uuid,
		Type:      b.Type,
		CreatedAt: time.Now(),
		VMDetail: []vm.VMDetail{{
			Node: hostname,
			VM:   t,
			Stat: uint(stat),
		}},
	}

}

func (b *Base) getVMAll() {
	// Get All
	log.Println("WebSocket VM GetAll")

	sshHosts, err := config.CollectConfig(nil)
	if err != nil {
		b.Error(err.Error())
		return
	}

	var vms []vm.VMDetail

	for _, tmpHost := range sshHosts {
		log.Printf("[%s] %s\n", tmpHost.HostName, tmpHost.User)
		conn, err := ConnectLibvirt(tmpHost.HostName, tmpHost.User)
		if err != nil {
			log.Println("failed to connect to qemu: " + err.Error())
			continue
		}
		defer conn.Close()

		net, _ := conn.ListNetworks()
		log.Println(net)
		doms, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE | libvirt.CONNECT_LIST_DOMAINS_INACTIVE)
		if err != nil {
			log.Printf("ListAllDomains error: %s", err)
			continue
		}
		for _, dom := range doms {
			t := libvirtxml.Domain{}
			stat, _, _ := dom.GetState()
			xmlString, _ := dom.GetXMLDesc(libvirt.DOMAIN_XML_SECURE)
			xml.Unmarshal([]byte(xmlString), &t)

			vms = append(vms, vm.VMDetail{
				Node: tmpHost.HostName,
				VM:   t,
				Stat: uint(stat),
			})
		}

		vm.ClientBroadcast <- vm.WebSocketResult{
			UUID:      b.UUID,
			Type:      b.Type,
			CreatedAt: time.Now(),
			VMDetail:  vms,
		}
	}
}
