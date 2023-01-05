package v3

import (
	"github.com/vmmgr/controller/pkg/api/core/tool/config"
	"github.com/vmmgr/controller/pkg/api/core/tool/remote"
	"log"
	"strings"
	"testing"
)

func Test1(t *testing.T) {
	sh := NewStorageHandler(StorageHandler{SSHHost: config.SSHHost{
		User: "opsadm", HostName: "10.100.1.180", KeyPath: "/Users/y-yoneda/.ssh/doornoc_id_rsa",
	}})
	re := remote.Auth{Config: sh.SSHHost}
	//qemu-img create -f qcow2 file.qcow2 100M
	file := "/etc/resolv.conf"
	command := "FILE=" + file
	command += `
if [ -f "$FILE" ]; then
    echo true
else 
    echo false
fi
`
	log.Println(command)
	result, err := re.SSHClientExecCmd(command)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
	if strings.Contains(result, "true") {
		t.Log("Contain!!")
	}
}

//func Test1(t *testing.T) {
//	sh := connect.Auth{
//		IP:   "192.168.22.132",
//		Port: 22,
//		User: "yonedayuto",
//	}
//	//qemu-img create -f qcow2 file.qcow2 100M
//	command := "/home/yonedayuto/imacon copy --uuid 2e684c62-d680-40f9-818b-2919ca02507e --url http://localhost:8081/api/v1/controller --src /home/yonedayuto/image/focal-server-cloudimg-amd64-disk-kvm.img --dst /home/yonedayuto/Documents/vmmgr/vm-image/test1.img --addr 192.168.22.1:22 --user yonedayuto --config /home/yonedayuto/config.yaml"
//	log.Println(command)
//	result, err := sh.SSHClientExecCmd(command)
//	if err != nil {
//		t.Fatal(err)
//	}
//	t.Log(result)
//}
