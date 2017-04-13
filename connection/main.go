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
	RedisPool *redis.Pool
	config    = cfg.Config
	err       error
)

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
				WriteLog("error", err.Error())
			}
			return c, err
		},
	}
}
