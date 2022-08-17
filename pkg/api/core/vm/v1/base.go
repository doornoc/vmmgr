package v1

import (
	"fmt"
	"github.com/doornoc/vmmgr/pkg/api/core/tool/config"
	"libvirt.org/go/libvirt"
	"log"
)

type Base struct {
	UUID string
	Type uint
}

func ConnectLibvirt(hostName, user string) (*libvirt.Connect, error) {
	log.Println("qemu+ssh://" + user + "@" + hostName + "/system")
	return libvirt.NewConnect("qemu+ssh://" + user + "@" + hostName + "/system")
}

func getOneVM(hostname, uuid string) (*libvirt.Domain, error) {
	sshHost, err := config.CollectConfig(&hostname)
	if err != nil {
		return nil, err
	}

	conn, err := ConnectLibvirt(sshHost[0].HostName, sshHost[0].User)
	if err != nil {
		log.Println("failed to connect to qemu: " + err.Error())
		return nil, fmt.Errorf("failed to connect to qemu: " + err.Error())
	}
	defer conn.Close()

	dom, err := conn.LookupDomainByUUIDString(uuid)
	if err != nil {
		return nil, fmt.Errorf("failed to looking domain: " + err.Error())
	}

	return dom, nil
}
