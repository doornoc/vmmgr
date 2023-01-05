package v1

import (
	"fmt"
	"github.com/vmmgr/controller/pkg/api/core/tool/config"
	"github.com/vmmgr/controller/pkg/api/core/tool/remote"
	"github.com/vmmgr/controller/pkg/api/core/tool/template"
	"github.com/vmmgr/controller/pkg/api/core/vm"
	cloudinitInterface "github.com/vmmgr/controller/pkg/api/core/vm/cloudinit"
	cloudinit "github.com/vmmgr/controller/pkg/api/core/vm/cloudinit/v1"
	nic "github.com/vmmgr/controller/pkg/api/core/vm/nic/v3"
	storage "github.com/vmmgr/controller/pkg/api/core/vm/storage/v3"
	"libvirt.org/go/libvirt"
	"log"
	"time"
)

func (b *Base) createVM(hostname string, data vm.VMInput) {
	log.Println("WebSocket Store Create " + "(" + hostname + ")")
	vm.ClientBroadcast <- vm.WebSocketResult{
		Type:      vm.MessageTypeCreateVM,
		CreatedAt: time.Now(),
		VMDetail: []vm.VMDetail{{
			Node: hostname,
		}},
		Data: map[string]string{"create_progress": "0", "copy_progress": "0", "message": "create start..."},
	}

	sshHost, err := config.CollectConfig(&hostname)
	if err != nil {
		b.Error("failed to config collect config: " + err.Error())
		return
	}

	conn, err := ConnectLibvirt(sshHost[0].HostName, sshHost[0].User)
	if err != nil {
		b.Error("failed to connect to qemu: " + err.Error())
		return
	}
	defer conn.Close()
	doms, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE | libvirt.CONNECT_LIST_DOMAINS_INACTIVE)
	if err != nil {
		b.Error(fmt.Sprintf("ListAllDomains error: %s", err.Error()))
		return
	}
	// 同じVM名がないか確認
	sameName := false
	for _, dom := range doms {
		vmName, _ := dom.GetName()
		if vmName == data.Name {
			sameName = true
			break
		}
	}
	if sameName {
		b.Error("failed to same vm name")
		return
	}

	// 処理中のVMを確認
	for _, v := range Store {
		if v.VMInfo.Name == data.Name {
			b.Error("failed to same vm name. (" + data.Name + ")")
			return
		}
	}

	// Get Template
	tpl, err := template.Get(sshHost[0])
	if err != nil {
		log.Printf("Get Template error: %s", err)
		return
	}

	if data.IsCloudInit {
		// select image_template
		imageTemplate, err := CheckImageTemplate(tpl.ImageTemplate, data.CloudInit.ID)
		if err != nil {
			b.Error(err.Error())
			return
		}

		// check spec
		if err = CheckImageTemplateSpec(imageTemplate.SpecPlans, data); err != nil {
			b.Error(err.Error())
			return
		}

		// check storage size
		if len(data.Disk) == 0 {
			b.Error("disk none")
			return
		}

		if _, err = CheckImageTemplateStorage(imageTemplate.StoragePlans, data.Disk[0].Path, data.Disk[0].Size); err != nil {
			b.Error(err.Error())
			return
		}

		//　テンプレートからbaseパスを抽出
		var destBasePath string
		var destImagePath string
		for _, tmpTpl := range tpl.Storage {
			if tmpTpl.Name == data.Disk[0].Path {
				destBasePath = tmpTpl.Path
				extension := "raw"
				if data.CloudInit.TemplateType != "" {
					extension = "qcow2"
				}
				destImagePath = fmt.Sprintf("%s/%s_0.%s", tmpTpl.Path, data.Name, extension)
				break
			}
		}

		// register store
		Store[b.UUID] = StoreStruct{
			StartTime:     time.Now(),
			UpdateTime:    time.Now(),
			HostName:      hostname,
			VMInfo:        data,
			Template:      tpl,
			ImageTemplate: *imageTemplate,
			DestBasePath:  destBasePath,
		}

		// storage newHandler
		sto := storage.NewStorageHandler(storage.StorageHandler{
			SSHHost:  sshHost[0],
			Template: tpl,
			Chan:     &vm.ClientBroadcast,
		})

		// notice
		vm.ClientBroadcast <- vm.WebSocketResult{
			Type:      vm.MessageTypeCreateVM,
			CreatedAt: time.Now(),
			VMDetail: []vm.VMDetail{{
				Node: hostname,
			}},
			Data: map[string]string{"create_progress": "10", "copy_progress": "0", "message": "copy template image..."},
		}

		if err = sto.ImageCopy(b.UUID, imageTemplate.Path, destImagePath); err != nil {
			b.Error("[ImageCopy] " + err.Error())
			delete(Store, b.UUID)
			return
		}

		return
	}

	// addr
	addr := vm.Address{
		PCICount:  0,
		DiskCount: 0,
	}

	// init generate xml
	domCfg, err := initGenerateXml(data, doms)
	if err != nil {
		b.Error(err.Error())
	}

	// storage
	sto := storage.NewStorageHandler(storage.StorageHandler{
		SSHHost:  sshHost[0],
		Template: tpl,
		Address:  &addr,
		Xml:      domCfg,
		Chan:     &vm.ClientBroadcast,
	})
	err = sto.Add(data.Name, data.Disk)
	if err != nil {
		b.Error("failed to storage apply: " + err.Error())
		return
	}

	// nic
	ni := nic.NewNICHandler(nic.NICHandler{
		Conn:     conn,
		Doms:     doms,
		SSHHost:  sshHost[0],
		Template: tpl,
		Address:  &addr,
		Xml:      domCfg,
		Chan:     &vm.ClientBroadcast,
	})
	err = ni.Add(data.NIC)
	if err != nil {
		b.Error("failed to nic apply: " + err.Error())
		return
	}

	xml, err := domCfg.Marshal()
	if err != nil {
		b.Error("domCfg marshal error: " + err.Error())
		return
	}
	log.Println("xml", xml)
	dom, err := conn.DomainDefineXML(xml)
	if err != nil {
		b.Error("vm domainDefineXML error: " + err.Error())
		return
	}

	err = dom.Create()
	if err != nil {
		b.Error("vm create error: " + err.Error())
		return
	}
	vm.ClientBroadcast <- vm.WebSocketResult{
		Type:      vm.MessageTypeCreateVM,
		CreatedAt: time.Now(),
		VMDetail: []vm.VMDetail{{
			Node: hostname,
		}},
		Data: map[string]string{"create_progress": "100", "copy_progress": "100", "message": "create finished...", "xml": xml},
	}
}

func CreateForCloudInit(uuidStr string) {
	storeVM := Store[uuidStr]
	vm.ClientBroadcast <- vm.WebSocketResult{
		Type:      vm.MessageTypeCreateVM,
		CreatedAt: time.Now(),
		VMDetail: []vm.VMDetail{{
			Node: storeVM.HostName,
		}},
		Data: map[string]string{"create_progress": "40", "copy_progress": "100", "message": "copy template image..."},
	}

	sshHost, err := config.CollectConfig(&storeVM.HostName)
	if err != nil {
		storeVM.Base.Error("failed to config collect config: " + err.Error())
		return
	}

	conn, err := ConnectLibvirt(sshHost[0].HostName, sshHost[0].User)
	if err != nil {
		storeVM.Base.Error("failed to connect to qemu: " + err.Error())
		return
	}
	defer conn.Close()
	doms, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE | libvirt.CONNECT_LIST_DOMAINS_INACTIVE)
	if err != nil {
		storeVM.Base.Error(fmt.Sprintf("ListAllDomains error: %s", err.Error()))
		return
	}

	// addr
	addr := vm.Address{
		PCICount:  0,
		DiskCount: 0,
	}

	// init generate xml
	domCfg, err := initGenerateXml(storeVM.VMInfo, doms)
	if err != nil {
		storeVM.Base.Error(err.Error())
		return
	}

	// storage
	sto := storage.NewStorageHandler(storage.StorageHandler{
		SSHHost:     sshHost[0],
		Template:    storeVM.Template,
		Address:     &addr,
		Xml:         domCfg,
		Chan:        &vm.ClientBroadcast,
		IsCloudinit: true,
	})

	// nic
	ni := nic.NewNICHandler(nic.NICHandler{
		Conn:     conn,
		Doms:     doms,
		SSHHost:  sshHost[0],
		Template: storeVM.Template,
		Address:  &addr,
		Xml:      domCfg,
		Chan:     &vm.ClientBroadcast,
	})

	storeVM.VMInfo.NIC, err = ni.GenerateOnlyMac(storeVM.VMInfo.NIC)
	if err != nil {
		storeVM.Base.Error("failed to GenerateOnlyMac: " + err.Error())
		return
	}

	networkConfig, err := cloudinit.GenNetworkConfig(storeVM.VMInfo.NIC)
	if err != nil {
		storeVM.Base.Error("failed to GenNetworkConfig: " + err.Error())
		return
	}

	// convert raw and size(require time)
	// TODO: ストレージの複数対応
	// notify
	vm.ClientBroadcast <- vm.WebSocketResult{
		Type:      vm.MessageTypeCreateVM,
		CreatedAt: time.Now(),
		VMDetail: []vm.VMDetail{{
			Node: storeVM.HostName,
		}},
		Data: map[string]string{"create_progress": "50", "copy_progress": "100", "message": "convert raw and size... (require time!!!)"},
	}
	extension := "raw"
	if storeVM.VMInfo.CloudInit.TemplateType != "" {
		extension = "qcow2"
	}

	err = sto.ConvertRawAndSize(
		storage.GenTemplatePath(storeVM.DestBasePath, storeVM.VMInfo.Name, 0, extension),
		storage.GenTemplatePath(storeVM.DestBasePath, storeVM.VMInfo.Name, 0, "raw"),
		storeVM.VMInfo.CloudInit.TemplateType,
		storeVM.VMInfo.Disk[0].Size,
		&vm.WebSocketResult{
			Type: vm.MessageTypeCreateVM,
			VMDetail: []vm.VMDetail{{
				Node: storeVM.HostName,
			}},
		},
	)
	if err != nil {
		storeVM.Base.Error("failed to convert raw and size: " + err.Error())
		return
	}

	// notify
	vm.ClientBroadcast <- vm.WebSocketResult{
		Type:      vm.MessageTypeCreateVM,
		CreatedAt: time.Now(),
		VMDetail: []vm.VMDetail{{
			Node: storeVM.HostName,
		}},
		Data: map[string]string{"create_progress": "90", "copy_progress": "100", "message": "generate cloudinit file..."},
	}

	// cloudinit
	cinh := cloudinit.NewHandler(cloudinit.Handler{
		DirPath: fmt.Sprintf("%s/%s", storeVM.DestBasePath, storeVM.VMInfo.Name),
		Auth:    remote.Auth{Config: sshHost[0]},
		MetaData: cloudinitInterface.MetaData{
			InstanceID:    uuidStr,
			LocalHostName: storeVM.VMInfo.Name,
		},
		UserData:      storeVM.VMInfo.CloudInit.UserData,
		NetworkConfig: networkConfig,
	})
	err = cinh.Generate()
	if err != nil {
		storeVM.Base.Error("failed to cloudinit apply : " + err.Error())
		return
	}

	// change to raw fileType
	storeVM.VMInfo.Disk[0].FileType = 1

	// add cloudinit image
	storeVM.VMInfo.Disk = append(storeVM.VMInfo.Disk, vm.VMDisk{
		Type:     1,
		Path:     fmt.Sprintf("%s/%s/cloudinit.img", storeVM.DestBasePath, storeVM.VMInfo.Name),
		ReadOnly: true,
	})

	err = sto.Add(storeVM.VMInfo.Name, storeVM.VMInfo.Disk)
	if err != nil {
		storeVM.Base.Error("failed to storage apply: " + err.Error())
		return
	}

	err = ni.Add(storeVM.VMInfo.NIC)
	if err != nil {
		storeVM.Base.Error("failed to nic apply: " + err.Error())
		return
	}

	xml, err := domCfg.Marshal()
	if err != nil {
		storeVM.Base.Error("domCfg marshal error: " + err.Error())
		return
	}
	log.Println("xml", xml)
	dom, err := conn.DomainDefineXML(xml)
	if err != nil {
		storeVM.Base.Error("vm domainDefineXML error: " + err.Error())
		return
	}

	err = dom.Create()
	if err != nil {
		storeVM.Base.Error("vm create error: " + err.Error())
		return
	}

	delete(Store, uuidStr)

	vm.ClientBroadcast <- vm.WebSocketResult{
		Type:      vm.MessageTypeCreateVM,
		CreatedAt: time.Now(),
		VMDetail: []vm.VMDetail{{
			Node: storeVM.HostName,
		}},
		Data: map[string]string{"create_progress": "100", "copy_progress": "100", "message": "create finished...", "xml": xml},
	}
}
