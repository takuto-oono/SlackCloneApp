package config

import (
	"fmt"
	"os"

	"gopkg.in/ini.v1"
)

type ConfigList struct {
	// web関連
	BackendBaseUrl  string
	FrontendBaseUrl string

	//db関連
	Driver        string
	DbName        string
	UserTableName string
}

var Config ConfigList

func init() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	Config = ConfigList{
		BackendBaseUrl:  cfg.Section("web").Key("backendBaseUrl").String(),
		FrontendBaseUrl: cfg.Section("web").Key("frontendBaseUrl").String(),
		Driver:          cfg.Section("db").Key("driver").String(),
		DbName:          cfg.Section("db").Key("dbName").String(),
		UserTableName:   cfg.Section("db").Key("userTableName").String(),
	}
}
