package v1

import (
	"encoding/xml"
	"github.com/doornoc/vmmgr/pkg/api/core/tool/config"
	"github.com/doornoc/vmmgr/pkg/api/core/tool/gen"
	"github.com/doornoc/vmmgr/pkg/api/core/vm"
	"github.com/gin-gonic/gin"
	"libvirt.org/go/libvirt"
	"log"
	"time"
)

func GetWebSocketAdmin(c *gin.Context) {
	conn, err := vm.WsUpgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %+v\n", err)
		return
	}

	defer conn.Close()

	uuid := gen.GenerateUUID()

	// WebSocket送信
	vm.Clients[&vm.WebSocket{
		UUID:  uuid,
		Admin: true,
		//GroupID: 0,
		Socket: conn,
	}] = true

	//WebSocket受信
	for {
		var msg vm.WebSocketInput
		err = conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(vm.Clients, &vm.WebSocket{UUID: uuid, Admin: true, GroupID: 0, Socket: conn})
			break
		}

		if msg.Type == vm.MessageTypeGetVM {
			// Get
			log.Println("WebSocket VM Get " + msg.UUID)
			_, conn, err := connectLibvirt(msg.NodeID)
			if err != nil {
				log.Println(err)
				continue
			}

			dom, err := conn.LookupDomainByUUIDString(msg.UUID)
			if err != nil {
				log.Println(err)
				continue
			}

			t := libVirtXml.Domain{}
			stat, _, _ := dom.GetState()
			xmlString, _ := dom.GetXMLDesc(libvirt.DOMAIN_XML_SECURE)
			xml.Unmarshal([]byte(xmlString), &t)

			vm.ClientBroadcast <- vm.WebSocketResult{
				UUID:      uuid,
				Type:      0,
				CreatedAt: time.Now(),
				Status:    true,
				Code:      0,
				VMDetail:  []vm.Detail{{VM: t, Stat: uint(stat)}},
			}

		} else if msg.Type == vm.MessageTypeGetVMAll {
			// Get All
			log.Println("WebSocket VM GetAll")

			sshHosts, err := config.CollectConfig(nil)
			if err != nil {
				vm.ClientBroadcast <- vm.WebSocketResult{
					UUID:      uuid,
					Type:      vm.MessageTypeGetVMAll,
					Err:       err.Error(),
					CreatedAt: time.Now(),
				}
			}

			resultNode := dbNode.GetAll()
			if resultNode.Err != nil {
				log.Println(resultNode.Err)
				vm.ClientBroadcast <- vm.WebSocketResult{
					UUID:      uuid,
					Type:      1,
					Err:       resultNode.Err.Error(),
					CreatedAt: time.Now(),
					Status:    false,
					Code:      0,
				}
				continue
			}

			var vms []vm.Detail

			for _, tmpNode := range resultNode.Node {
				log.Printf("[%s] %s\n", tmpNode.IP, tmpNode.User)
				conn, err := libvirt.NewConnect("qemu+ssh://" + tmpNode.User + "@" + tmpNode.HostName + "/system")
				if err != nil {
					log.Println("failed to connect to qemu: " + err.Error())
					//vm.ClientBroadcast <- vm.WebSocketResult{
					//	UUID:      uuid,
					//	Type:      1,
					//	Err:       err.Error(),
					//	CreatedAt: time.Now(),
					//	Status:    false,
					//	Code:      0,
					//}
					continue
				}
				defer conn.Close()

				net, _ := conn.ListNetworks()
				log.Println(net)
				doms, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE | libvirt.CONNECT_LIST_DOMAINS_INACTIVE)
				if err != nil {
					log.Printf("ListAllDomains error: %s", err)
					//vm.ClientBroadcast <- vm.WebSocketResult{
					//	UUID:      uuid,
					//	Type:      1,
					//	Err:       err.Error(),
					//	CreatedAt: time.Now(),
					//	Status:    false,
					//	Code:      0,
					//}
					continue
				} else {
					for _, dom := range doms {
						t := libVirtXml.Domain{}
						stat, _, _ := dom.GetState()
						xmlString, _ := dom.GetXMLDesc(libvirt.DOMAIN_XML_SECURE)
						xml.Unmarshal([]byte(xmlString), &t)

						vms = append(vms, vm.Detail{
							Node: tmpNode.ID,
							VM:   t,
							Stat: uint(stat),
						})
					}

					vm.ClientBroadcast <- vm.WebSocketResult{
						UUID:      uuid,
						Type:      1,
						CreatedAt: time.Now(),
						Status:    true,
						Code:      0,
						VMDetail:  vms,
					}
				}
			}
		} else if msg.Type == vm.MessageTypeCreateVM {
		} else if msg.Type == vm.MessageTypeDeleteVM {
			// Delete
		} else if msg.Type == vm.MessageTypeStartVM {
			// Start
			detail, err := Startup(msg.NodeID, msg.UUID)
			if err != nil {
				log.Println(err)
				continue
			}
			vm.ClientBroadcast <- vm.WebSocketResult{
				UUID:      msg.UUID,
				Type:      20,
				Err:       "",
				CreatedAt: time.Now(),
				Status:    true,
				VMDetail:  []vm.Detail{*detail},
			}
		} else if msg.Type == vm.MessageTypeForceShutdownVM {
			// Force Shutdown
			detail, err := Shutdown(msg.NodeID, msg.UUID, true)
			if err != nil {
				log.Println(err)
				continue
			}
			vm.ClientBroadcast <- vm.WebSocketResult{
				UUID:      msg.UUID,
				Type:      21,
				Err:       "",
				CreatedAt: time.Now(),
				Status:    true,
				VMDetail:  []vm.Detail{*detail},
			}
		} else if msg.Type == vm.MessageTypeShutdownVM {
			// Shutdown
			detail, err := Shutdown(msg.NodeID, msg.UUID, false)
			if err != nil {
				log.Println(err)
				continue
			}
			vm.ClientBroadcast <- vm.WebSocketResult{
				UUID:      msg.UUID,
				Type:      22,
				Err:       "",
				CreatedAt: time.Now(),
				Status:    true,
				VMDetail:  []vm.Detail{*detail},
			}
		} else if msg.Type == vm.MessageTypeResetVM {
			// Reset
			detail, err := Reset(msg.NodeID, msg.UUID)
			if err != nil {
				log.Println(err)
				continue
			}
			vm.ClientBroadcast <- vm.WebSocketResult{
				UUID:      msg.UUID,
				Type:      23,
				Err:       "",
				CreatedAt: time.Now(),
				Status:    true,
				VMDetail:  []vm.Detail{*detail},
			}
		}
	}
}

func HandleMessages() {
	for {
		msg := <-vm.ClientBroadcast

		//登録されているクライアント宛にデータ送信する
		//コントローラが管理者モードの場合
		for client := range vm.Clients {
			err := client.Socket.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Socket.Close()
				delete(vm.Clients, client)
			}
		}
	}
}
