package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type ListernerConfig struct {
	Address string `yaml:"address"`
	Cluster string `yaml:"cluster"`
	Tlsname string `yaml:"tls"`
}

type ClusterConfig struct {
	Name     string   `yaml:"name"`
	Endpoint []string `yaml:"endpoints"`
	TlsName  string   `yaml:"tls"`
}

type TlsConfig struct {
	Name string `yaml:"name"`
	Cert string `yaml:"cert"`
	Key  string `yaml:"key"`
	CA   string `yaml:"ca"`
}

type GlobalConfig struct {
	Listeners []ListernerConfig `yaml:"listeners"`
	TlsCfg    []TlsConfig       `yaml:"tls"`
	Clusters  []ClusterConfig   `yaml:"clusters"`
}

var globalconfig *GlobalConfig

func LoadConfig(filename string) error {

	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	config := new(GlobalConfig)
	config.Listeners = make([]ListernerConfig, 0)
	config.Clusters = make([]ClusterConfig, 0)
	config.TlsCfg = make([]TlsConfig, 0)

	err = yaml.Unmarshal(body, config)
	if err != nil {
		return err
	}
	
	globalconfig = config
	return nil
}

func listenerGetAll() []ListernerConfig {
	return globalconfig.Listeners
}

func ClusterGet(name string) *ClusterConfig {
	for _, v := range globalconfig.Clusters {
		if v.Name == name {
			return &v
		}
	}
	return nil
}

func TlsGet(name string) *TlsConfig {
	for _, v := range globalconfig.TlsCfg {
		if v.Name == name {
			return &v
		}
	}
	return nil
}
