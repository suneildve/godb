package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	Server 			map[string]string 	`json:"server"`
	Redis			*RedisConfig 		`json:"redis"`
	DB 				map[string]string 	`json:"db"`
	Design			map[string]string 	`json:"design"`
	Version 		string 				`json:"version"`
}

type RedisConfig struct {
	Addr 		string 	`json:"addr"`
	Password 	string 	`json:"password"`
	PoolSize 	int 	`json:"poolsize"`
	DBs 		[]int 	`json:"dbs"`
}

var config *Config

func InitConfig(confPath string) (*Config, error) {
	if data, err := ioutil.ReadFile(confPath); err != nil {
		return nil, err
	} else {
		var conf = &Config{}
		err := json.Unmarshal(data, conf)
		if err != nil {
			fmt.Println("json decode error")
			return nil, err
		}
		config = conf
		return conf, nil
	}
}

//GetConfig 获取配置
func GetConfig() *Config {
	return config
}