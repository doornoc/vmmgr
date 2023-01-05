package v1

import (
	"fmt"
	"github.com/vmmgr/controller/pkg/api/core/vm"
	"github.com/vmmgr/controller/pkg/api/core/vm/cloudinit"
)

func GenNetworkConfig(nics []vm.VMNIC) (cloudinit.NetworkCon, error) {
	networkConf := cloudinit.NetworkCon{
		Version: 1,
	}

	nicName := "eth"
	count := 0

	for _, nic := range nics {
		if nic.CloudInit.Address == "" {
			continue
		}
		networkConf.Config = append(networkConf.Config, cloudinit.NetworkConfig{
			Type:       "physical",
			Name:       fmt.Sprintf("%s%d", nicName, count),
			MacAddress: nic.MAC,
			Subnets: []cloudinit.NetworkConfigSubnet{
				{
					Type:    "static",
					Address: nic.CloudInit.Address,
					Netmask: nic.CloudInit.Netmask,
					Gateway: nic.CloudInit.Gateway,
					DNS:     nic.CloudInit.DNS,
				},
			},
		})
	}

	return networkConf, nil
}
