package db

import (
	"github.com/go-redis/redis"
	"godb/config"
	"fmt"
	"time"
	"log"
)

type RedisDB int
const (
	// RedisDBUseGame       RedisDB = iota 	//游戏角色相关数据 0
	// RedisDBUseBattleLoad                    //战场负载数据 1
	// RedisDBUseBattleInfo                    //战斗相关数据 2
	// RedisDBFriend                           //好友相关数据 3
	// RedisDBGuild                            //公会相关数据 4
	// RedisDBConfig        = 10               //一些及时配置
	// RedisDBUseMax
)

type RedisMgr struct {
	clients map[RedisDB]*redis.Client
}
var redisMgr *RedisMgr

func InitRedisDB() bool {
	r := new(RedisMgr)
	r.clients = make(map[RedisDB]*redis.Client)
	redisMgr = r
	var conf = config.GetConfig()
	if conf.Redis != nil {
		for _, v := range conf.Redis.DBs {
			if !redisMgr.NewRedisClient(v,conf.Redis) {
				return false
			}
		}
	}
	return true
}

func (mgr *RedisMgr) NewRedisClient(dbIndex int, redisConf *config.RedisConfig) bool {
	redclient := redis.NewClient(&redis.Options{
		Addr:         redisConf.Addr, //":6379",
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     redisConf.PoolSize, //10,
		Password:     redisConf.Password,
		PoolTimeout:  30 * time.Second,
		DB:           dbIndex,
	})
	fmt.Printf("dbs %d\n",dbIndex)
	// _, err := redclient.Set("test____", 1, time.Second*10).Result()
	_, err := redclient.Ping().Result()
	if err != nil {
		log.Printf("%v", err)
		return false
	}
	mgr.clients[RedisDB(dbIndex)] = redclient
	return true
}
func GetRedisDB(dbIndex RedisDB) *redis.Client {
	if client, b := redisMgr.clients[dbIndex]; b {
		return client
	}
	log.Printf("GetRedisDB no exist:%v", dbIndex)
	return nil
}