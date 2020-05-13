package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

var (
	Config tomlConfig
)

type tomlConfig struct {
	Server  serverConfig
	BaseLog logConfig
	DataLog logConfig
	Redis   redisConfig
	Mysql   mysqlConfig
	Trade   tradeConfig
	Service serviceConfig
}

type serverConfig struct {
	Address     string
	Location    string
	BaseLogPath string
	ServerName  string
	ServiceName string
	ServerTag   string
}

type logConfig struct {
	FileName       string
	FileMaxSize    int
	FileMaxBackups int
	FileMaxAge     int
}

type redisConfig struct {
	Address   string
	Password  string
	MaxIdle   int
	MaxActive int
	IdleMills int
}

type mysqlConfig struct {
	User        string
	Password    string
	Address     string
	DataBase    string
	MaxOpenConn int
	MaxIdleConn int
}

type tradeConfig struct {
	Endpoint        			string
	WSEndpoint 					string
	ApiKey                   	string
	SecretKey                   string
	Passphrase                  string
	TimeoutSecond               int
	IsPrint                    	bool
	I18n                 	 	string
}

type serviceConfig struct {
	NotifyUrl string
}

func LoadConfig(pathToToml string) {
	if _, err := toml.DecodeFile(pathToToml, &Config); err != nil {
		fmt.Println(err)
		panic(err)
	}
}
