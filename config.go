package main

import (
	"io/ioutil"

	"github.com/mikespook/golib/log"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Log struct {
		File, Level string
	}
	Addr string
}

func InitConfig(filename string) (cfg *Config, err error) {
	var data []byte
	if data, err = ioutil.ReadFile(filename); err != nil {
		return
	}
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return
	}
	err = log.Init(cfg.Log.File, log.StrToLevel(cfg.Log.Level), 0)
	return
}
