package v3

import (
	"encoding/xml"
	"fmt"
	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
	"log"
	"sort"
	"strconv"
	"strings"
)

var maxMAC = 65535
var startMAC = 10

func (h *NICHandler) generateMac(usedMAC []string) (string, error) {
	log.Println("Generate MAC")
	//log.Println(h.BaseMAC)

	var macs []int
	//startMACを定義
	macIndex := startMAC

	var doms []libvirt.Domain

	doms, err := h.Conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE | libvirt.CONNECT_LIST_DOMAINS_INACTIVE)
	if err != nil {
		log.Println(err)
	}

	// Todo:
	if len(doms) != 0 {
		for _, dom := range doms {
			data := libvirtxml.Domain{}
			xmlString, _ := dom.GetXMLDesc(libvirt.DOMAIN_XML_SECURE)
			xml.Unmarshal([]byte(xmlString), &data)

			if len(data.Devices.Interfaces) != 0 {
				for _, tmp := range data.Devices.Interfaces {
					mac := strings.Split(tmp.MAC.Address, ":")
					if (mac[0] + mac[1]) == "5254" {
						v, _ := strconv.ParseInt(mac[4]+mac[5], 16, 0)
						macs = append(macs, int(v))
					}
				}
			}
		}
		//割当済みMACアドレスを検索して、macsに値を代入
		for _, tmp := range usedMAC {
			mac := strings.Split(tmp, ":")
			if (mac[0] + mac[1]) == "5254" {
				v, _ := strconv.ParseInt(mac[4]+mac[5], 16, 0)
				macs = append(macs, int(v))
			}
		}

		//昇順に並び替える
		sort.Ints(macs)

		for _, m := range macs {
			//Port番号が上限に達する場合、エラーを返す
			if maxMAC <= macIndex {
				return "", fmt.Errorf("Error: max mac address ")
			}
			if macIndex < m {
				break
			}
			macIndex++
		}
	}

	macIndex1 := macIndex / 256
	macIndex2 := macIndex % 256

	// macアドレスを10進数から16進数に変換し、結合
	mac := fmt.Sprintf("52:54:%s:%.2x:%.2x", h.Template.BaseMacAddress, macIndex1, macIndex2)

	return mac, nil
}
