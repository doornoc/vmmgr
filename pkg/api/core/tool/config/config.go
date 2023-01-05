package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

type Config struct {
	Port        int      `yaml:"port"`
	LocalUrl    string   `yaml:"url"`
	TimeZone    string   `yaml:"timezone"`
	AcceptHosts []string `yaml:"accept_hosts"`
}

var Conf Config

func GetConfig(inputConfPath string) error {
	configPath := "./config.yaml"
	if inputConfPath != "" {
		configPath = inputConfPath
	}
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}
	var data Config
	err = yaml.Unmarshal(file, &data)
	if err != nil {
		log.Fatal(err)
	}
	Conf = data
	return nil
}
