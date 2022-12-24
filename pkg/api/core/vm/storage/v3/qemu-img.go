package v3

import (
	"github.com/vmmgr/controller/pkg/api/core/tool/remote"
	"github.com/vmmgr/controller/pkg/api/core/vm/storage"
	"log"
	"strconv"
)

func (h *StorageHandler) convertImage(d storage.Convert) (string, error) {
	sh := remote.Auth{
		Config: h.SSHHost,
	}

	//qemu-img convert -f [src_type] -O [dest_type] [src_path] [dest_path]
	command := "qemu-img" + " convert" + " -f " + d.SrcType + " -O " + d.DstType + " " + d.SrcFile + " " + d.DstFile
	result, err := sh.SSHClientExecCmd(command)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return result, nil
}

func (h *StorageHandler) generateImage(fileType, filePath string, fileSize uint) (string, error) {
	sh := remote.Auth{
		Config: h.SSHHost,
	}

	size := strconv.Itoa(int(fileSize)) + "M"
	//qemu-img create -f [file_type] [file_path] 100M
	command := "qemu-img " + "create " + "-f " + fileType + " " + filePath + " " + size
	result, err := sh.SSHClientExecCmd(command)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return result, nil
}

func (h *StorageHandler) infoImage(filePath string) (string, error) {
	sh := remote.Auth{
		Config: h.SSHHost,
	}

	//qemu-img info [file_path]
	command := "qemu-img " + "info " + filePath
	result, err := sh.SSHClientExecCmd(command)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return result, nil
}

func (h *StorageHandler) capacityExpansion(filePath string, size uint) (string, error) {
	sh := remote.Auth{
		Config: h.SSHHost,
	}

	//qemu-img resize [file_path]
	command := "qemu-img " + "info " + filePath
	result, err := sh.SSHClientExecCmd(command)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return result, nil
}
