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
	Driver                        string
	DbName                        string
	UserTableName                 string
	WorkspaceTableName            string
	WorkspaceAndUserTableName     string
	RoleTableName                 string
	ChannelsTableName             string
	ChannelsAndUserTableName      string

	//jwt-token
	TokenHourLifeSpan string
	SecretKey         string
}

var Config ConfigList

func loadConfigFile() (*ini.File, error) {
	file, err := ini.Load("config.ini")
	if err == nil {
		return file, nil
	}
	file, err = ini.Load("../config.ini")

	return file, err
}

func init() {
	cfg, err := loadConfigFile()
	if err != nil {
		fmt.Printf("Failed to read file: %v", err)
		os.Exit(1)

	}

	Config = ConfigList{
		BackendBaseUrl:  cfg.Section("web").Key("backendBaseUrl").String(),
		FrontendBaseUrl: cfg.Section("web").Key("frontendBaseUrl").String(),

		Driver:                        cfg.Section("db").Key("driver").String(),
		DbName:                        cfg.Section("db").Key("dbName").String(),
		UserTableName:                 cfg.Section("db").Key("usersTableName").String(),
		WorkspaceTableName:            cfg.Section("db").Key("workspacesTableName").String(),
		WorkspaceAndUserTableName:     cfg.Section("db").Key("workspacesAndUsersTableName").String(),
		RoleTableName:                 cfg.Section("db").Key("rolesTableName").String(),
		ChannelsTableName:             cfg.Section("db").Key("channelsTableName").String(),
		ChannelsAndUserTableName:      cfg.Section("db").Key("channelsAndUsersTableName").String(),

		TokenHourLifeSpan: cfg.Section("jwt-token").Key("tokenHourLifespan").String(),
		SecretKey:         cfg.Section("jwt-token").Key("secretKey").String(),
	}
}
