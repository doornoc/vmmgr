package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	Controller Controller `json:"controller"`
}

type Controller struct {
	Port     int    `json:"port"`
	TimeZone string `json:"timezone"`
}

var Conf Config

func GetConfig(inputConfPath string) error {
	configPath := "./data.json"
	if inputConfPath != "" {
		configPath = inputConfPath
	}
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}
	var data Config
	err = json.Unmarshal(file, &data)
	if err != nil {
		log.Fatal(err)
	}
	Conf = data
	return nil
}
