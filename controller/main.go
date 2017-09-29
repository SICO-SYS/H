/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"log"

	"github.com/SiCo-Ops/cfg"
	"github.com/SiCo-Ops/dao/redis"
)

const (
	configPath string = "config.json"
)

var (
	config     cfg.ConfigItems
	publicPool = redis.NewPool()
)

func init() {
	data, err := cfg.ReadFilePath(configPath)
	if err != nil {
		data = cfg.ReadConfigServer()
		if data == nil {
			log.Fatalln("config.json not exist and configserver was down")
		}
	}
	cfg.Unmarshal(data, &config)

	publicPool = redis.InitPool(config.RedisPublicHost, config.RedisPublicPort, config.RedisPublicAuth)

}
