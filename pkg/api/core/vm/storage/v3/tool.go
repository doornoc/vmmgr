package v3

import (
	"fmt"
	"github.com/vmmgr/controller/pkg/api/core/vm"
	"github.com/vmmgr/controller/pkg/api/core/vm/storage"
)

func GenTemplatePath(templatePath, name string, idx int, extension string) string {
	return fmt.Sprintf("%s/%s_%d.%s", templatePath, name, idx, extension)
}

func convertStorageVM(disk vm.VMDisk) storage.VMStorage {
	return storage.VMStorage{
		Type:     disk.Type,
		FileType: disk.FileType,
		Path:     disk.Path,
		ReadOnly: disk.ReadOnly,
	}
}

func convertFileTypeToString(fileType uint) (string, error) {
	switch fileType {
	case 0:
		// qcow2
		return "qcow2", nil
	case 1:
		// raw
		return "raw", nil
	}
	return "", fmt.Errorf("error: file format type")
}
