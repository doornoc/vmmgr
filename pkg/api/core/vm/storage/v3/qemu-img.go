package v3

import (
	"fmt"
	"github.com/pkg/sftp"
	"github.com/vmmgr/controller/pkg/api/core/tool/remote"
	"github.com/vmmgr/controller/pkg/api/core/vm/storage"
	"log"
	"strings"
)

func (h *StorageHandler) ConvertImage(d storage.Convert) (string, error) {
	sh := remote.Auth{
		Config: h.SSHHost,
	}

	//qemu-img convert -t none -T none -O [dest_type] [src_path] [dest_path] -o preallocation=falloc
	command := fmt.Sprintf("qemu-img convert -p -t none -T none -O %s %s %s -o preallocation=falloc", d.DstType, d.SrcFile, d.DstFile)
	log.Println("[ConvertImage]", command)
	result, err := sh.SSHClientExecCmd(command)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (h *StorageHandler) generateImage(fileType, filePath string, fileSize uint) (string, error) {
	sh := remote.Auth{
		Config: h.SSHHost,
	}

	//qemu-img create -f [file_type] [file_path] 100M
	command := fmt.Sprintf("qemu-img create -f %s %s %dM", fileType, filePath, fileSize)
	log.Println("[generateImage]", command)
	result, err := sh.SSHClientExecCmd(command)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (h *StorageHandler) infoImage(filePath string) (string, error) {
	sh := remote.Auth{
		Config: h.SSHHost,
	}

	//qemu-img info [file_path]
	command := fmt.Sprintf("qemu-img info %s", filePath)
	log.Println("[infoImage]", command)
	result, err := sh.SSHClientExecCmd(command)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (h *StorageHandler) CapacityExpansion(filePath string, size uint) (string, error) {
	sh := remote.Auth{
		Config: h.SSHHost,
	}

	command := fmt.Sprintf("qemu-img resize %s %dM", filePath, size)
	log.Println("[CapacityExpansion]", command)
	//qemu-img resize [file_path]
	result, err := sh.SSHClientExecCmd(command)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (h *StorageHandler) deleteFile(filepath string) error {
	sh := remote.Auth{
		Config: h.SSHHost,
	}

	conn, err := sh.SSHClient()
	if err != nil {
		return fmt.Errorf("[sh.SSHClient()] %s", err)
	}
	defer conn.Close()

	client, err := sftp.NewClient(conn)
	if err != nil {
		return fmt.Errorf("[sftp.NewClient(conn)] %s", err)
	}
	defer client.Close()

	return client.Remove(filepath)

}

func (h *StorageHandler) fileExist(filepath string) (bool, error) {

	sh := remote.Auth{
		Config: h.SSHHost,
	}

	file := "/etc/resolv.conf"
	command := "FILE=" + file
	command += `
if [ -f "$FILE" ]; then
    echo true
else 
    echo false
fi
`
	result, err := sh.SSHClientExecCmd(command)
	if err != nil {
		return false, fmt.Errorf("[sh.SSHClientExecCmd(command)] %s", err)
	}

	if strings.Contains(result, "true") {
		return true, nil
	}

	return false, nil
}
