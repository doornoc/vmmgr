package template

type Template struct {
	Storage        []Storage `json:"storage"`
	NIC            []NIC     `json:"nic"`
	BaseMacAddress string    `json:"base_mac_address"`
}

type Storage struct {
	Name    string        `json:"name"`
	Comment string        `json:"comment"`
	Path    string        `json:"path"`
	Option  StorageOption `json:"option"`
}

type StorageOption struct {
	IsIso      bool `json:"is_iso"`
	IsCloudimg bool `json:"is_cloudimg"`
}

type NIC struct {
	Name      string `json:"name"`
	Comment   string `json:"comment"`
	Interface string `json:"interface"`
}

type ImageList struct {
	Name     string   `json:"name"`
	BasePath string   `json:"base_path"`
	Path     []string `json:"path"`
}
