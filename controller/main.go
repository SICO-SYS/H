/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"github.com/SiCo-Ops/cfg/v2"
	"github.com/SiCo-Ops/dao/redis"
)

var (
	config     cfg.ConfigItems
	configPool = redis.Pool("", "", "")
	publicPool = redis.Pool("", "", "")
	RPCAddr    map[string]string
)

func init() {
	data := cfg.ReadLocalFile()

	if data != nil {
		cfg.Unmarshal(data, &config)
	}

	configPool = redis.Pool(config.RedisConfigHost, config.RedisConfigPort, config.RedisConfigAuth)
	configs, _ := redis.Hgetall(configPool, "system.config")
	cfg.Map2struct(configs, &config)
	publicPool = redis.Pool(config.RedisPublicHost, config.RedisPublicPort, config.RedisPublicAuth)

	RPCAddr = map[string]string{
		"He": config.RpcHeHost + ":" + config.RpcHePort,
		"Li": config.RpcLiHost + ":" + config.RpcLiPort,
		"Be": config.RpcBeHost + ":" + config.RpcBePort,
		"B":  config.RpcBHost + ":" + config.RpcBPort,
		"C":  config.RpcCHost + ":" + config.RpcCPort,
		"N":  config.RpcNHost + ":" + config.RpcNPort,
		"O":  config.RpcOHost + ":" + config.RpcOPort,
		"F":  config.RpcFHost + ":" + config.RpcFPort,
		"Ne": config.RpcNeHost + ":" + config.RpcNePort,
	}
}
