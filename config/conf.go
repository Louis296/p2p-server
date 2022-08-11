package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Conf struct {
	HTTPS  HTTPS  `yaml:"https"`
	Server Server `yaml:"server"`
}
type HTTPS struct {
	Cert string `yaml:"cert"`
	Key  string `yaml:"key"`
}
type Server struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	HTMLRoot string `yaml:"html_root"`
}

func GetConfig() (Conf, error) {
	file, err := ioutil.ReadFile("config/conf.yaml")
	if err != nil {
		return Conf{}, err
	}
	var conf Conf
	err = yaml.Unmarshal(file, &conf)
	if err != nil {
		return Conf{}, err
	}
	return conf, nil
}
