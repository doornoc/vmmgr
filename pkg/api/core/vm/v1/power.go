package v1

import (
	"encoding/xml"
	"github.com/vmmgr/controller/pkg/api/core/vm"
	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
	"log"
	"time"
)

func (b *Base) startVM(hostname, uuid string) {
	log.Println("WebSocket Store Start " + uuid + "(" + hostname + ")")
	dom, err := getOneVM(hostname, uuid)
	if err != nil {
		b.Error("failed to getting state: " + err.Error())
		return
	}

	stat, _, err := dom.GetState()
	if err != nil {
		b.Error("failed to getting state: " + err.Error())
		return
	}

	if stat != libvirt.DOMAIN_RUNNING {
		if err = dom.Create(); err != nil {
			b.Error("failed to dom create: " + err.Error())
			return
		}
	}

	t := libvirtxml.Domain{}
	stat, _, _ = dom.GetState()
	xmlString, _ := dom.GetXMLDesc(libvirt.DOMAIN_XML_SECURE)
	xml.Unmarshal([]byte(xmlString), &t)

	err = dom.Free()
	if err != nil {
		b.Error("failed to dom free: " + err.Error())
		return
	}

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

func (b *Base) shutdownVM(force bool, hostname, uuid string) {
	log.Println("WebSocket Store Shutdown/Force Shutdown " + uuid + "(" + hostname + ")")
	dom, err := getOneVM(hostname, uuid)
	if err != nil {
		b.Error("failed to getting state: " + err.Error())
		return
	}

	stat, _, err := dom.GetState()
	if err != nil {
		b.Error("failed to getting state: " + err.Error())
		return
	}

	if stat != libvirt.DOMAIN_SHUTOFF {
		// Forceがtrueである場合、強制終了
		if force {
			if err = dom.Destroy(); err != nil {
				log.Println(err)
				b.Error("failed to force shutdown state: " + err.Error())
				return
			}
		} else {
			if err = dom.Shutdown(); err != nil {
				log.Println(err)
				b.Error("failed to shutdown state: " + err.Error())
				return
			}
		}
	}

	t := libvirtxml.Domain{}
	stat, _, _ = dom.GetState()
	xmlString, _ := dom.GetXMLDesc(libvirt.DOMAIN_XML_SECURE)
	xml.Unmarshal([]byte(xmlString), &t)

	err = dom.Free()
	if err != nil {
		b.Error("failed to dom free: " + err.Error())
		return
	}

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

func (b *Base) resetVM(hostname, uuid string) {
	log.Println("WebSocket Store Reset " + uuid + "(" + hostname + ")")
	dom, err := getOneVM(hostname, uuid)
	if err != nil {
		b.Error("failed to getting state: " + err.Error())
		return
	}

	stat, _, err := dom.GetState()
	if err != nil {
		b.Error("failed to getting state: " + err.Error())
		return
	}

	if err = dom.Reset(0); err != nil {
		b.Error("failed to resetting state: " + err.Error())
		return
	}

	t := libvirtxml.Domain{}
	stat, _, _ = dom.GetState()
	xmlString, _ := dom.GetXMLDesc(libvirt.DOMAIN_XML_SECURE)
	xml.Unmarshal([]byte(xmlString), &t)

	err = dom.Free()
	if err != nil {
		b.Error("failed to dom free: " + err.Error())
		return
	}

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
