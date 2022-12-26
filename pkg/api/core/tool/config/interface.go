package config

type SSHHost struct {
	User     string `json:"user"`
	HostName string `json:"host_name"`
	KeyPath  string `json:"key_path"`
}
