package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/vmmgr/controller/pkg/api/core/tool/gen"
	"github.com/vmmgr/controller/pkg/api/core/vm"
	"log"
)

func GetWebSocketAdmin(c *gin.Context) {
	conn, err := vm.WsUpgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %+v\n", err)
		return
	}

	defer conn.Close()

	uuid := gen.GenerateUUID()
	b := Base{UUID: uuid}

	// WebSocket送信
	vm.Clients[&vm.WebSocket{
		UUID:   uuid,
		Socket: conn,
	}] = true

	//WebSocket受信
	for {
		var msg vm.WebSocketInput
		err = conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(vm.Clients, &vm.WebSocket{
				UUID:   uuid,
				Socket: conn,
			})
			break
		}

		b.Type = msg.Type

		switch msg.Type {
		case vm.MessageTypeGetVM:
			// Get
			b.getVM(msg.Data["hostname"], msg.Data["uuid"])
		case vm.MessageTypeGetVMAll:
			b.getVMAll()
		case vm.MessageTypeCreateVM:
			break
		case vm.MessageTypeDeleteVM:
			break
		case vm.MessageTypeStartVM:
			// Start
			b.startVM(msg.Data["hostname"], msg.Data["uuid"])
		case vm.MessageTypeForceShutdownVM:
			// Force Shutdown
			b.shutdownVM(true, msg.Data["hostname"], msg.Data["uuid"])
		case vm.MessageTypeShutdownVM:
			// Shutdown
			b.shutdownVM(false, msg.Data["hostname"], msg.Data["uuid"])
		case vm.MessageTypeResetVM:
			// Reset
			b.resetVM(msg.Data["hostname"], msg.Data["uuid"])
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
