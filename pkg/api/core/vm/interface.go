package vm

import (
	"github.com/gorilla/websocket"
	"libvirt.org/go/libvirtxml"
	"net/http"
	"time"
)

// 定数
const MessageTypeGetVM = 1
const MessageTypeGetVMAll = 2
const MessageTypeCreateVM = 51
const MessageTypeDeleteVM = 61
const MessageTypeStartVM = 101
const MessageTypeForceShutdownVM = 102
const MessageTypeShutdownVM = 103
const MessageTypeResetVM = 104

//const messageGetType = 31

// channel定義(websocketで使用)
var Clients = make(map[*WebSocket]bool)
var ClientBroadcast = make(chan WebSocketResult)

type WebSocket struct {
	UUID   string
	Socket *websocket.Conn
}

type WebSocketInput struct {
	Type uint              `json:"type"`
	Data map[string]string `json:"data"`
}

// websocket用
type WebSocketResult struct {
	ID        uint              `json:"id"`
	Err       string            `json:"error"`
	CreatedAt time.Time         `json:"created_at"`
	Type      uint              `json:"type"`
	UUID      string            `json:"uuid"`
	VMDetail  []VMDetail        `json:"vm_detail"`
	Data      map[string]string `json:"data"`
}

type VMDetail struct {
	VM   libvirtxml.Domain `json:"vm"`
	Stat uint              `json:"stat"`
	Node string            `json:"node"`
}

var WsUpgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
