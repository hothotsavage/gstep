package config

import (
	"encoding/json"
	"github.com/hothotsavage/gstep/util/LogUtil"
	"log"
	"os"
)

type Configuration struct {
	Port uint64
	Auth struct {
		Secret string
	}
	Db struct {
		Database string
		Host     string
		Port     string
		User     string
		Password string
	}
	Department struct {
		RootParentDepartmentId string
	}
	Nacos struct {
		Host        string
		Port        uint64
		Namespace   string
		ServiceIP   string
		ServiceName string
	}
	IsDebugLocal bool
}

// 全局配置
var Config = &Configuration{}

func Setup() {
	//将配置文件:config.json中的配置读取到Config
	//---debug---------------------------------------
	file, err := os.Open("config_debug_local.json")
	//file, err := os.Open("./config_debug_nacos.json")
	//file, err := os.Open("./config_release.json")
	if err != nil {
		log.Printf("cannot open file config.LogUtil: %v", err)
		panic(err)
	}

	decoder := json.NewDecoder(file)
	Config = &Configuration{}
	err = decoder.Decode(Config)
	if err != nil {
		log.Printf("decode config.LogUtil failed: %v", err)
		panic(err)
	}

	log.Printf("global config:")
	LogUtil.PrintPretty(Config)
}
