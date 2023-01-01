package template

type Template struct {
	Storage        []Storage `yaml:"storage"`
	NIC            []NIC     `yaml:"nic"`
	BaseMacAddress string    `yaml:"base_mac_address"`
}

type Storage struct {
	Name    string        `yaml:"name"`
	Comment string        `yaml:"comment"`
	Path    string        `yaml:"path"`
	Option  StorageOption `yaml:"option"`
}

type StorageOption struct {
	IsIso      bool `yaml:"is_iso"`
	IsCloudimg bool `yaml:"is_cloudimg"`
}

type NIC struct {
	Name      string `yaml:"name"`
	Comment   string `yaml:"comment"`
	Interface string `yaml:"interface"`
}

type ImageList struct {
	Name     string   `yaml:"name"`
	BasePath string   `yaml:"base_path"`
	Path     []string `yaml:"path"`
}
