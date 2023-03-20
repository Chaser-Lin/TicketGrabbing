package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func Load() error {
	content, err := getConfigContent()
	if err != nil {
		return err
	}
	return yaml.Unmarshal(content, &Conf)
}

func getConfigContent() ([]byte, error) {
	return ioutil.ReadFile("conf.yaml")
}
