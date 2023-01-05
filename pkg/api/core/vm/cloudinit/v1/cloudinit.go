package v1

import (
	"fmt"
	"github.com/pkg/sftp"
	"github.com/vmmgr/controller/pkg/api/core/tool/remote"
	"github.com/vmmgr/controller/pkg/api/core/vm/cloudinit"
	"gopkg.in/yaml.v2"
	"path/filepath"
)

type Handler struct {
	DirPath       string
	Auth          remote.Auth
	MetaData      cloudinit.MetaData   `json:"meta"`
	UserData      cloudinit.UserData   `json:"user"`
	NetworkConfig cloudinit.NetworkCon `json:"network"`
}

func NewHandler(input Handler) *Handler {
	return &Handler{
		DirPath:       input.DirPath,
		Auth:          input.Auth,
		MetaData:      input.MetaData,
		UserData:      input.UserData,
		NetworkConfig: input.NetworkConfig,
	}
}

func (c *Handler) Generate() error {
	// SFTP
	conn, err := c.Auth.SSHClient()
	if err != nil {
		return fmt.Errorf("[c.Auth.SSHClient()] %s", err)
	}
	defer conn.Close()

	client, err := sftp.NewClient(conn)
	if err != nil {
		return fmt.Errorf("[sftp.NewClient(conn)] %s", err)
	}
	defer client.Close()

	// make directory
	client.MkdirAll(c.DirPath)
	// MetaData
	file, err := client.Create(c.DirPath + "/meta-data")
	if err != nil {
		return fmt.Errorf("[Error: create metadata file] %s", err)
	}
	metaDataYAML, err := yaml.Marshal(c.MetaData)
	if err != nil {
		return fmt.Errorf("[Error: marshal metadata yaml file] %s", err)
	}
	_, err = file.Write(metaDataYAML)
	if err != nil {
		return fmt.Errorf("[Error: write metadata] %s", err)
	}

	// UserData
	file, err = client.Create(c.DirPath + "/user-data")
	if err != nil {
		return fmt.Errorf("[Error: create user-data file] %s", err)
	}
	userDataYAML, err := yaml.Marshal(c.UserData)
	if err != nil {
		return fmt.Errorf("[Error: marshal user-data yaml file] %s", err)
	}
	userDataYAML = []byte(fmt.Sprintf("#cloud-config\n%s", userDataYAML))
	_, err = file.Write(userDataYAML)
	if err != nil {
		return fmt.Errorf("[Error: write user-data] %s", err)
	}

	// NetworkConfig
	// NetworkConfig Version指定
	file, err = client.Create(c.DirPath + "/network-config")
	if err != nil {
		return fmt.Errorf("[Error: create network-config file] %s", err)

	}

	c.NetworkConfig.Version = 1
	networkConfigYAML, err := yaml.Marshal(c.NetworkConfig)
	if err != nil {
		return fmt.Errorf("[Error: marshal network-config yaml file] %s", err)
	}
	_, err = file.Write(networkConfigYAML)
	if err != nil {
		return fmt.Errorf("[Error: write neteworkconfig] %s", err)
	}

	//ファイルがすでに存在する場合は削除する
	client.Remove(filepath.Join(c.DirPath, "cloudinit.img"))

	ah := c.Auth
	_, err = ah.SSHClientExecCmd(
		fmt.Sprintf("cloud-localds -N %s %s %s %s",
			filepath.Join(c.DirPath, "network-config"), filepath.Join(c.DirPath, "cloudinit.img"),
			filepath.Join(c.DirPath, "user-data"), filepath.Join(c.DirPath, "meta-data")),
	)

	return err
}
