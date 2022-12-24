package template

import (
	"encoding/json"
	"github.com/vmmgr/controller/pkg/api/core/tool/config"
	"github.com/vmmgr/controller/pkg/api/core/tool/remote"
	"log"
	"strings"
)

func Get(host config.SSHHost) (Template, error) {
	var tpl Template
	sh := remote.Auth{
		Config: host,
	}

	//cat data.json
	command := "cat data.json"
	result, err := sh.SSHClientExecCmd(command)
	if err != nil {
		log.Println(err)
		return tpl, err
	}
	log.Println(result)

	err = json.Unmarshal([]byte(result), &tpl)
	if err != nil {
		log.Fatal(err)
	}

	return tpl, nil
}

func GetList(host config.SSHHost, basePath string, name string) (ImageList, error) {
	var imageList = ImageList{BasePath: basePath, Name: name}
	sh := remote.Auth{
		Config: host,
	}

	//cat data.json
	command := "ls -m " + basePath
	result, err := sh.SSHClientExecCmd(command)
	if err != nil {
		log.Println(err)
		return imageList, err
	}

	if len(result) == 0 {
		return imageList, nil
	}

	var path []string
	path = strings.Split(result, ",")
	for idx := range path {
		path[idx] = basePath + "/" + strings.TrimSpace(path[idx])
	}

	imageList.Path = path

	return imageList, nil
}
