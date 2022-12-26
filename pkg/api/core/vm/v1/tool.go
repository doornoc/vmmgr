package v1

import (
	"encoding/json"
	"github.com/vmmgr/controller/pkg/api/core/tool/config"
	"github.com/vmmgr/controller/pkg/api/core/tool/template"
	"github.com/vmmgr/controller/pkg/api/core/vm"
	"log"
	"time"
)

func (b *Base) getTemplate(hostname string) {
	log.Printf("WebSocket Get Template: [%s]", hostname)
	sshHost, err := config.CollectConfig(&hostname)
	if err != nil {
		b.Error(err.Error())
		return
	}

	tpl, err := template.Get(sshHost[0])
	if err != nil {
		b.Error(err.Error())
		return
	}
	log.Println(tpl)
	tplStr, _ := json.Marshal(&tpl)

	// get ISO & cloudimg data
	var imageLists []template.ImageList
	for _, storage := range tpl.Storage {
		if storage.Option.IsIso || storage.Option.IsCloudimg {
			list, err := template.GetList(sshHost[0], storage.Path, storage.Name)
			if err != nil {
				b.Error(err.Error())
			}
			imageLists = append(imageLists, list)
		}
	}
	imageListStr, _ := json.Marshal(&imageLists)

	vm.ClientBroadcast <- vm.WebSocketResult{
		Type:      b.Type,
		CreatedAt: time.Now(),
		Data:      map[string]string{"template": string(tplStr), "image_list": string(imageListStr)},
	}
}

func (b *Base) getHost() {
	log.Printf("WebSocket Get Host ")
	sshHost, err := config.CollectConfig(nil)
	if err != nil {
		b.Error(err.Error())
		return
	}

	sshHostStr, _ := json.Marshal(&sshHost)

	vm.ClientBroadcast <- vm.WebSocketResult{
		Type:      b.Type,
		CreatedAt: time.Now(),
		Data:      map[string]string{"hosts": string(sshHostStr)},
	}
}
