package v3

import (
	"github.com/vmmgr/controller/pkg/api/core/tool/config"
	"strings"
	"testing"
)

var confPath = "/home/yonedayuto/go/src/github.com/vmmgr/controller/cmd/backend/config.yaml"

func Test0(t *testing.T) {
	err := config.GetConfig(confPath)
	if err != nil {
		t.Fatalf("error config process |%v", err)
	}
}

func Test1(t *testing.T) {
	tmp := connect.Auth{
		IP:   "192.168.122.1",
		Port: 22,
		User: "yonedayuto",
		Pass: "",
	}
	result, err := tmp.SSHClientExecCmd("ls /sys/class/net")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
	resultArray := strings.Fields(result)
	t.Log(resultArray)
}
