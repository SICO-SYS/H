/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package connection

import (
	"github.com/garyburd/redigo/redis"

	"github.com/SiCo-DevOps/H/cfg"
	. "github.com/SiCo-DevOps/H/log"
)

var (
	config = cfg.Config
	err    error
)

func init() {
	RedisPool = &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", config.Redis.Default.Host+":"+config.Redis.Default.Port)
			if err != nil {
				WriteLog("error", err.Error())
			}
			return c, err
		},
	}

	err = RedisPool.Get().Close()
	if err != nil {
		WriteLog("error", "Cannot Open redis connection")
	} else {
		WriteLog("info", "Success connect redis")
	}
}
