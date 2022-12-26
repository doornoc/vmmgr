package v3

import (
	"fmt"
	"github.com/vmmgr/controller/pkg/api/core/tool/template"
	"github.com/vmmgr/controller/pkg/api/core/vm"
	"github.com/vmmgr/controller/pkg/api/core/vm/storage"
	"libvirt.org/go/libvirtxml"
)

func (h *StorageHandler) xmlGenerate(input []vm.VMDisk) error {

	var disks []libvirtxml.DomainDisk

	//countの定義＆初期化
	var virtIOCount uint = 0
	var otherCount uint = 0

	for _, storageTmp := range input {
		if storageTmp.Path == "" {
			return fmt.Errorf("black: storage path")
		}

		var Number uint
		var AddressNumber uint

		// VirtIOの場合はVirtIO Countに数字を代入＋加算する
		if storageTmp.Type == 0 {
			Number = virtIOCount
			AddressNumber = h.Address.PCICount
			h.Address.PCICount++
			virtIOCount++
		} else {
			Number = otherCount
			// VirtIO以外の場合はOther Countに数字を代入＋加算する
			AddressNumber = h.Address.DiskCount
			otherCount++
			h.Address.DiskCount++
		}

		disks = append(disks, *generateTemplate(storage.GenerateStorageXml{
			Storage:       convertStorageVM(storageTmp),
			Number:        Number,
			PCISlot:       virtIOCount,
			AddressNumber: AddressNumber,
		}))
	}

	h.Xml.Devices.Disks = disks

	return nil
}

func generateTemplate(xmlStruct storage.GenerateStorageXml) *libvirtxml.DomainDisk {
	//デフォルトはブートディスク(VirtIO)

	domDisk := libvirtxml.DomainDisk{}
	var dev string
	var bus string
	// CDROM
	if xmlStruct.Storage.Type == 1 {
		dev = "sda"
		bus = "sata"
		domDisk.Device = "cdrom"
		domDisk.Address = &libvirtxml.DomainAddress{
			Drive: &libvirtxml.DomainAddressDrive{
				Controller: &[]uint{0}[0],
				Bus:        &[]uint{0}[0],
				Target:     &[]uint{0}[0],
				Unit:       &[]uint{xmlStruct.AddressNumber}[0],
			},
		}
		// Boot Disk(SATA)
	} else if xmlStruct.Storage.Type == 11 {
		dev = "sda"
		bus = "sata"
		domDisk.Address = &libvirtxml.DomainAddress{
			Drive: &libvirtxml.DomainAddressDrive{
				Controller: &[]uint{0}[0],
				Bus:        &[]uint{0}[0],
				Target:     &[]uint{0}[0],
				Unit:       &[]uint{xmlStruct.AddressNumber}[0],
			},
		} // Boot Disk(IDE)
	} else if xmlStruct.Storage.Type == 12 {
		dev = "sda"
		bus = "ide"

		domDisk.Address = &libvirtxml.DomainAddress{
			Drive: &libvirtxml.DomainAddressDrive{
				Controller: &[]uint{0}[0],
				Bus:        &[]uint{0}[0],
				Target:     &[]uint{0}[0],
				Unit:       &[]uint{xmlStruct.AddressNumber}[0],
			},
		}
	} else {
		dev = "vda"
		bus = "virtio"
		domDisk.Address = &libvirtxml.DomainAddress{
			PCI: &libvirtxml.DomainAddressPCI{
				Domain:   &[]uint{0}[0],
				Bus:      &[]uint{1}[0],
				Slot:     &[]uint{xmlStruct.PCISlot}[0],
				Function: &[]uint{0}[0],
			},
		}
	}

	if xmlStruct.Storage.ReadOnly {
		domDisk.ReadOnly = &libvirtxml.DomainDiskReadOnly{}
	}

	domDisk.Target = &libvirtxml.DomainDiskTarget{Bus: bus, Dev: dev[0:2] + string(dev[2]+uint8(xmlStruct.Number))}
	// Driver
	if xmlStruct.Storage.Type == 1 || xmlStruct.Storage.Type == 2 {
		domDisk.Driver = &libvirtxml.DomainDiskDriver{
			Name: "qemu",
			Type: "raw",
		}
	} else {
		domDisk.Driver = &libvirtxml.DomainDiskDriver{
			Name: "qemu",
			Type: template.GetExtensionStr(xmlStruct.Storage.FileType),
		}
	}
	// File Path
	domDisk.Source = &libvirtxml.DomainDiskSource{
		File: &libvirtxml.DomainDiskSourceFile{File: xmlStruct.Storage.Path},
	}

	return &domDisk
}
