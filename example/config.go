package main

import (
	"github.com/swxctx/goai/baidu"
	td "github.com/swxctx/malatd"
	"github.com/swxctx/xlog"
	"github.com/usthooz/gconf"
)

type Config struct {
	SrvConfig *td.SrvConfig `json:"srv_config"`
}

var cfg = &Config{
	SrvConfig: td.NewSrvConfig(),
}

func reload() {
	conf := gconf.NewConf(&gconf.Gconf{
		ConfPath: "./config/config.yaml",
	})

	// get config
	err := conf.GetConf(&cfg)
	if err != nil {
		xlog.Errorf("GetConf Err: %v", err.Error())
	}

	// 初始化厂商SDK
	if err := baidu.NewClient("T4YTCcTsdfsdfdfsd", "qbqc3B2vd", true); err != nil {
		xlog.Errorf("Config: baidu.NewClient err-> %v", err)
	}
}

func init() {
	reload()
}
