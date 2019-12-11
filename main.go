package main

import (
	"godb/config"
	"godb/db"
	"godb/server"
	"fmt"
	"flag"
	"log"
	"godb/utils"
)

var (
	confPath = flag.String("config", "config.json", "配置文件")
)

func main() {
	flag.Parse()
	conf, err := config.InitConfig(*confPath)
	if err != nil {
		log.Println("load config err:", err)
	}
	fmt.Printf("json string: %s\n",conf.Version)
	fmt.Println(utils.Encrypt("suneil"))
	// db.InitRedisDB()
	db.InitMySqlDB()
	server.StartHTTP()
}

