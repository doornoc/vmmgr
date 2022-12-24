package config

import (
	"fmt"
	"github.com/kevinburke/ssh_config"
	"os"
	"path/filepath"
	"strings"
)

func CollectConfig(hostname *string) ([]SSHHost, error) {
	var sshHosts []SSHHost
	f, err := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "config"))
	if err != nil {
		return sshHosts, err
	}
	cfg, err := ssh_config.Decode(f)
	if err != nil {
		return sshHosts, err
	}
	for _, host := range cfg.Hosts {
		fmt.Println("patterns:", host.Patterns)
		var sshHost SSHHost
		if len(host.Nodes) < 2 {
			continue
		}
		for _, node := range host.Nodes {
			fmt.Println(node.String())
			nodeStringSplit := strings.Split(strings.TrimSpace(node.String()), " ")

			switch nodeStringSplit[0] {
			case "User":
				sshHost.User = nodeStringSplit[1]
			case "HostName":
				sshHost.HostName = nodeStringSplit[1]
			case "IdentityFile":
				sshHost.KeyPath = nodeStringSplit[1]
			}
		}
		if hostname != nil && *hostname != sshHost.HostName {
			continue
		}
		sshHosts = append(sshHosts, sshHost)
	}

	return sshHosts, nil
}
