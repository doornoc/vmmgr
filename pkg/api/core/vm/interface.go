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
const MessageTypeGetTemplate = 8
const MessageTypeGetHostName = 9
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

// websocket request
type WebSocketInput struct {
	Type    uint              `json:"type"`
	Data    map[string]string `json:"data"`
	VMInput VMInput           `json:"vm_input"`
}

type VMInput struct {
	Name        string       `json:"name"`
	IsCloudInit bool         `json:"is_cloud_init"`
	Arch        uint         `json:"arch"` //0: x86_64 1:x86
	CPU         uint         `json:"cpu"`
	Memory      uint         `json:"memory"`
	Boot        string       `json:"boot"` //hd
	Disk        []VMDisk     `json:"disk"`
	NIC         []VMNIC      `json:"nic"`
	CloudInit   *VMCloudInit `json:"cloud_init"`
}

type VMDisk struct {
	Type     uint   `json:"type"`      //0:BootDev(VirtIO) 1: CDROM 2:Floppy (no support) 11: BootDev(SATA) 12: BootDev(IDE)
	FileType uint   `json:"file_type"` //0:qcow2 1:raw
	Path     string `json:"path"`      //node側のパス or storage type(hdd1,hdd2,ssd1,ssd2,nvme1,nvme2)
	ReadOnly bool   `json:"readonly"`  //Readonlyであるか
	Size     uint   `json:"size"`
}

type VMNIC struct {
	Type   uint   `json:"type"`   //0: Bridge 1: NAT 2:macvtap
	Driver uint   `json:"driver"` // 0: virtio 1:e1000e 2:rtl8139
	Mode   uint   `json:"mode"`   //0: Bridge 1: vpea 2: private 3: passthrough
	MAC    string `json:"mac"`
	Device string `json:"device"`
}

type VMCloudInit struct {
	ImageCopy         string   `json:"os"` // http: wget, scp: scp, local: local copy
	Name              string   `yaml:"name"`
	Password          string   `yaml:"password"`
	Groups            string   `yaml:"groups"`
	Shell             string   `yaml:"shell"`
	Sudo              []string `yaml:"sudo"`
	SSHAuthorizedKeys []string `yaml:"ssh-authorized-keys"`
	SSHPWAuth         bool     `yaml:"ssh_pwauth"`
	LockPasswd        bool     `yaml:"lock_passwd"`
}

type Address struct {
	PCICount  uint
	DiskCount uint
}

// websocket response
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
