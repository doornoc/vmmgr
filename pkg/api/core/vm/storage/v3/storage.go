package v3

import (
	"fmt"
	"github.com/vmmgr/controller/pkg/api/core/tool/config"
	"github.com/vmmgr/controller/pkg/api/core/tool/remote"
	"github.com/vmmgr/controller/pkg/api/core/tool/template"
	vmInterface "github.com/vmmgr/controller/pkg/api/core/vm"
	"github.com/vmmgr/controller/pkg/api/core/vm/storage"
	"libvirt.org/go/libvirtxml"
	"log"
	"time"
)

//type StorageHandler struct {
//	UUID      string
//	Conn      *libvirt.Connect
//	Input     storage.Storage
//	Store        vm.VirtualMachine
//	Address   *vm.Address
//	Auth      *storage.Auth
//	SrcImaCon core.ImaCon
//	DstAuth   storage.Auth
//	SrcPath   string
//	DstPath   string
//	CtrlType  uint // 1:Admin 2:User
//}

type StorageHandler struct {
	SSHHost     config.SSHHost
	Template    template.Template
	Address     *vmInterface.Address
	Xml         *libvirtxml.Domain
	Chan        *chan vmInterface.WebSocketResult
	IsCloudinit bool
}

func NewStorageHandler(handler StorageHandler) *StorageHandler {
	return &handler
}

//func (h *StorageHandler) AddFromImage() error {
//	// ImaConからイメージ取得(時間がかかるので、go funcにて処理)
//	log.Println("From: " + h.SrcPath)
//	log.Println("To: " + h.DstPath)
//
//	sh := remote.Auth{
//		Config: h.SSHHost,
//	}
//
//	//qemu-img create -f qcow2 file.qcow2 100M
//	command := h.SrcImaCon.AppPath + " copy --uuid " + h.UUID + " --url " + url +
//		" --src " + h.SrcPath + " --dst " + h.DstPath + " --addr " + h.DstAuth.IP + ":" +
//		strconv.Itoa(int(h.DstAuth.Port)) + " --user " + h.DstAuth.User + " --config " + h.SrcImaCon.ConfigPath
//	log.Println(command)
//	result, err := sh.SSHClientExecCmd(command)
//	if err != nil {
//		log.Println(err)
//		return err
//	}
//
//	log.Println(result)
//
//	log.Println("Done: Image Create")
//
//	return nil
//}

func (h *StorageHandler) Add(name string, input []vmInterface.VMDisk) error {
	for idx, disk := range input {
		path := disk.Path
		extension, err := convertFileTypeToString(disk.FileType)
		if err != nil {
			return err
		}

		if disk.Type == 2 || disk.Type >= 13 {
			return fmt.Errorf("disk type error")
		}

		//　テンプレートからbaseパスを抽出
		for _, tpl := range h.Template.Storage {
			if tpl.Name == path {
				path = GenTemplatePath(tpl.Path, name, idx, extension)
				break
			}
		}

		// Pathが見つからない場合
		if path == "" {
			return fmt.Errorf("Error: Not found... ")
		}

		isExist, err := h.fileExist(path)
		// cloudinit=fase時は、fileが存在しないかCheckする
		if isExist && !h.IsCloudinit {
			return fmt.Errorf("Error: file already exists... ")
		} else if !isExist && h.IsCloudinit {
			// cloudinit=true時は、fileが存在するかCheckする
			return fmt.Errorf("Error: file is not exists... ")
		} else if err != nil {
			return fmt.Errorf("[h.fileExist] %s", err.Error())
		}

		// cloudinit=false時は、ここでストレージを作成
		if (0 == disk.Type || disk.Type >= 10) && !h.IsCloudinit {
			// イメージの作成
			out, err := h.generateImage(extension, path, disk.Size)
			if err != nil {
				log.Println(out)
				return err
			}
			log.Println("Done: [" + path + "] Image Create")
		}
		input[idx].Path = path
		//controller.SendServer(h.Input.Info, 1, 100, "Done: Image Create", nil)
	}
	err := h.xmlGenerate(input)
	if err != nil {
		return err
	}

	return nil
}

func (h *StorageHandler) ImageCopy(uuid, srcPath, destPath string) error {
	sh := remote.Auth{
		Config: h.SSHHost,
	}
	command := fmt.Sprintf(".vmmgr/node copy --controller %s --uuid %s --src %s --dest %s --debug true &", config.Conf.LocalUrl, uuid, srcPath, destPath)
	log.Println(command)
	_, err := sh.SSHClientExecCmd(command)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("[SSHClientExecCmd(%s)] %s", command, err.Error())
	}

	return nil
}

func (h *StorageHandler) ConvertRawAndSize(srcPath, destPath, srcType string, size uint, notice *vmInterface.WebSocketResult) error {
	// convert raw
	if srcType != "" {
		_, err := h.ConvertImage(storage.Convert{
			SrcFile: srcPath,
			SrcType: "qcow2",
			DstFile: destPath,
			DstType: srcType,
		})
		if err != nil {
			return fmt.Errorf("[ConvertImage(qcow2=>raw)] %s", err.Error())
		}
		vmInterface.ClientBroadcast <- vmInterface.WebSocketResult{
			Type:      vmInterface.MessageTypeCreateVM,
			CreatedAt: time.Now(),
			VMDetail:  notice.VMDetail,
			Data:      map[string]string{"create_progress": "70", "copy_progress": "100", "message": "expansion size... (require time!!!)"},
		}
	}
	// expansion size
	_, err := h.CapacityExpansion(destPath, size)
	if err != nil {
		return fmt.Errorf("[CapacityExpansion] %s", err.Error())
	}
	vmInterface.ClientBroadcast <- vmInterface.WebSocketResult{
		Type:      vmInterface.MessageTypeCreateVM,
		CreatedAt: time.Now(),
		VMDetail:  notice.VMDetail,
		Data:      map[string]string{"create_progress": "85", "copy_progress": "100", "message": "delete old file..."},
	}

	// delete old file
	if srcType != "" {
		err = h.deleteFile(srcPath)
		if err != nil {
			return fmt.Errorf("[deleteFile(srcPath)] %s", err.Error())
		}
		vmInterface.ClientBroadcast <- vmInterface.WebSocketResult{
			Type:      vmInterface.MessageTypeCreateVM,
			CreatedAt: time.Now(),
			VMDetail:  notice.VMDetail,
			Data:      map[string]string{"create_progress": "90", "copy_progress": "100", "message": "finish convert and size..."},
		}
	}

	return nil
}
