package gen

import (
	"encoding/xml"
	"fmt"
	"github.com/google/uuid"
	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
	"sort"
)

var portCount = 500
var vncPortStart = 5910
var webSocketPortStart = 5310

func GenerateUUID() string {
	u, err := uuid.NewRandom()
	if err != nil {
		fmt.Println(err)
		//return ""
	}
	uu := u.String()

	return uu
}

// Generate vnc and websocket port
func GenerateVNCAndWebSocketPort(doms []libvirt.Domain) (int, int, error) {
	//5900-6400
	var vncPort []int
	//5300-5800
	var webSocketPort []int

	for _, dom := range doms {
		t := libvirtxml.Domain{}
		xmlString, _ := dom.GetXMLDesc(libvirt.DOMAIN_XML_SECURE)
		xml.Unmarshal([]byte(xmlString), &t)
		if t.Devices == nil || len(t.Devices.Graphics) == 0 {
			continue
		}
		for _, gra := range t.Devices.Graphics {
			if gra.VNC == nil {
				continue
			}
			vncPort = append(vncPort, gra.VNC.Port)
			webSocketPort = append(webSocketPort, gra.VNC.WebSocket)
		}
	}

	//昇順に並び替える
	sort.Ints(vncPort)
	sort.Ints(webSocketPort)

	//ポート番号
	vncPortCount := vncPortStart
	webSocketPortCount := webSocketPortStart

	for _, port := range vncPort {
		//Port番号が上限に達する場合、エラーを返す
		if vncPortStart+portCount <= vncPortCount {
			return 0, 0, fmt.Errorf("Error: max port ")
		}

		if vncPortCount < port {
			break
		}
		vncPortCount++
	}

	for _, port := range webSocketPort {
		//Port番号が上限に達する場合、エラーを返す
		if webSocketPortCount+portCount <= webSocketPortCount {
			return 0, 0, fmt.Errorf("Error: max port ")
		}

		if webSocketPortCount < port {
			break
		}
		webSocketPortCount++
	}
	return vncPortCount, webSocketPortCount, nil
}
