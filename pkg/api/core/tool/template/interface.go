package template

type Template struct {
	Storage        []Storage       `json:"storage" yaml:"storage"`
	NIC            []NIC           `json:"nic" yaml:"nic"`
	ImageTemplate  []ImageTemplate `json:"image_template" yaml:"image_template"`
	BaseMacAddress string          `json:"base_mac_address" yaml:"base_mac_address"`
}

type Storage struct {
	Name    string        `json:"name" yaml:"name"`
	Comment string        `json:"comment" yaml:"comment"`
	Path    string        `json:"path" yaml:"path"`
	Option  StorageOption `json:"option" yaml:"option"`
}

type StorageOption struct {
	IsIso      bool `json:"is_iso" yaml:"is_iso"`
	IsCloudimg bool `json:"is_cloudimg" yaml:"is_cloudimg"`
}

type NIC struct {
	Name      string `json:"name" yaml:"name"`
	Comment   string `json:"comment" yaml:"comment"`
	Interface string `json:"interface" yaml:"interface"`
}

type ImageList struct {
	Name     string   `json:"name" yaml:"name"`
	BasePath string   `json:"base_path" yaml:"base_path"`
	Path     []string `json:"path" yaml:"path"`
}

type ImageTemplate struct {
	Name         string        `json:"name" yaml:"name"`
	Comment      string        `json:"comment" yaml:"comment"`
	Disable      bool          `json:"disable" yaml:"disable"`
	Path         string        `json:"path" yaml:"path"`
	SpecPlans    []SpecPlan    `json:"spec_plans" yaml:"spec_plans"`
	StoragePlans []StoragePlan `json:"storage_plans" yaml:"storage_plans"`
}

type SpecPlan struct {
	Name    string `json:"name" yaml:"name"`
	Disable bool   `json:"disable" yaml:"disable"`
	Arch    string `json:"arch" yaml:"arch"`
	CPU     uint   `json:"cpu" yaml:"cpu"`
	Memory  uint   `json:"memory" yaml:"memory"`
}

type StoragePlan struct {
	Name      string            `json:"name" yaml:"name"`
	Disable   bool              `json:"disable" yaml:"disable"`
	StorageID string            `json:"storage_id" yaml:"storage_id"`
	Size      []uint            `json:"size" yaml:"size"`
	Option    StoragePlanOption `json:"option" yaml:"option"`
}

type StoragePlanOption struct {
	IsNotExtension bool `json:"is_not_extension" yaml:"is_not_extension"`
}
