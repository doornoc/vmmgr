package v1

import (
	"github.com/vmmgr/controller/pkg/api/core/vm/cloudinit"
	"gopkg.in/yaml.v2"
	"testing"
)

func Test(t *testing.T) {
	tmpCloudInit := Handler{
		DirPath:  "",
		MetaData: cloudinit.MetaData{},
		UserData: cloudinit.UserData{
			Password:  "ubuntu",
			ChPasswd:  "{ expire: False }",
			SSHPwAuth: true,
		},
		NetworkConfig: cloudinit.NetworkCon{
			Version: 0,
			Config: []cloudinit.NetworkConfig{
				{
					Type:       cloudinit.NetworkConfigTypePhysical,
					Name:       "test",
					MacAddress: "00:11:22:33:44:55",
					//Subnets: []NetworkConfigSubnet{
					//	{
					//
					//	},
					//},
				},
			},
		},
	}

	yaml, _ := yaml.Marshal(tmpCloudInit)
	t.Log(string(yaml))
}
