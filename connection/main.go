/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package connection

import (
	"H/cfg"
	"H/log"
	"github.com/garyburd/redigo/redis"
)

var (
	RedisPool *redis.Pool
	config    = cfg.Config
	err       error
)

func GetRedisValue(v interface{}) (bool, error) {
	return redis.Bool(v, err)
}

func init() {
	// defer func() {
	// 	log.Println(recover())
	// }()
	RedisPool = &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", config.Redis.Default.Host+":"+config.Redis.Default.Port)
			if err != nil {
				log.WriteLog("error", err.Error())
			}
			return c, err
		},
	}
}
