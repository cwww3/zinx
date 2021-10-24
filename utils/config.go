package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type GlobalConfig struct {
	ServerConfig `json:"server"`
}

var Config *GlobalConfig

type ServerConfig struct {
	Name            string
	Host            string
	Port            int
	MaxPackageSize  uint32
	WorkerSize      int
	RequestPoolSize uint32
}

func init() {
	Config = &GlobalConfig{
		ServerConfig{
			Name:            "Zinx",
			Host:            "0.0.0.0",
			Port:            8999,
			MaxPackageSize:  512,
			RequestPoolSize: 100,
			WorkerSize:      10,
		},
	}
	Config.Reload()
}

func (c *GlobalConfig) Reload() {
	data, err := ioutil.ReadFile("config/config.json")
	if err != nil {
		fmt.Println("read config err", err)
		panic(err)
	}
	err = json.Unmarshal(data, &Config)
	if err != nil {
		fmt.Println("unmarshal config err", err)
		panic(err)
	}
}
