package cloudinit

// humstackのcloud-init部分を参考

type MetaData struct {
	InstanceID    string `yaml:"instance-id"`
	LocalHostName string `yaml:"local-hostname"`
}

type UserData struct {
	PackagesUpdate    bool             `json:"packages_update" yaml:"packages_update"`
	PackagesUpgrade   bool             `json:"packages_upgrade" yaml:"packages_upgrade"`
	Packages          []string         `json:"packages" yaml:"packages"`
	User              string           `json:"user" yaml:"user"`
	Password          string           `json:"password" yaml:"password"`
	ChPasswd          UserDataChPasswd `json:"ch_passwd" yaml:"chpasswd"`
	SSHPwAuth         bool             `json:"ssh_pwauth" yaml:"ssh_pwauth"`
	SSHAuthorizedKeys []string         `json:"ssh_authorized_keys" yaml:"ssh_authorized_keys"`
	Users             []UsersData      `json:"users" yaml:"users"`
}

type UserDataChPasswd struct {
	Expire bool `json:"expire" yaml:"expire"`
}

// TODO: keyにdefaultがないと動かない
// https://cloudinit.readthedocs.io/en/latest/topics/modules.html#users-and-groups

type UsersData struct {
	Name     string `json:"name" yaml:"name"`
	Password string `json:"password" yaml:"passwd"`
	//Groups            string  `json:"groups" yaml:"groups"`
	//Shell             string  `json:"shell" yaml:"shell"`
	Sudo              []string `json:"sudo" yaml:"sudo"`
	SSHAuthorizedKeys []string `json:"ssh_authorized_keys" yaml:"ssh_authorized_keys"`
	SSHPWAuth         bool     `json:"ssh_pwauth" yaml:"ssh_pwauth"`
	LockPasswd        bool     `json:"lock_passwd" yaml:"lock_passwd"`
}

type NetworkCon struct {
	Version int32           `json:"version" yaml:"version"`
	Config  []NetworkConfig `json:"config" yaml:"config"`
}

type NetworkConfigType string

const (
	NetworkConfigTypePhysical NetworkConfigType = "physical"
)

type NetworkConfig struct {
	Type       NetworkConfigType     `json:"type" yaml:"type"`
	Name       string                `json:"name" yaml:"name"`
	MacAddress string                `json:"mac_address" yaml:"mac_address"`
	Subnets    []NetworkConfigSubnet `json:"subnets" yaml:"subnets"`
}

type NetworkConfigSubnet struct {
	Type    NetworkConfigSubnetType `json:"type" yaml:"type"`
	Address string                  `json:"address" yaml:"address"`
	Netmask string                  `json:"netmask" yaml:"netmask"`
	Gateway string                  `json:"gateway" yaml:"gateway"`
	DNS     []string                `json:"dns" yaml:"dns_nameservers"`
}

type NetworkConfigSubnetType string

const (
	NetworkConfigSubnetTypeStatic NetworkConfigSubnetType = "static"
)
