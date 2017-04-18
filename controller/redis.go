/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"github.com/garyburd/redigo/redis"

	"github.com/SiCo-DevOps/dao"
)

var (
	RedisPool = dao.RedisPool
)

func RedisSetShort(k string, v interface{}, t int16) error {
	conn := RedisPool.Get()
	err = conn.Err()
	defer conn.Close()
	conn.Do("SET", k, v)
	conn.Do("EXPIRE", k, t)
	return err
}

func RedisSetLong(k string, v interface{}) error {
	conn := RedisPool.Get()
	err = conn.Err()
	defer conn.Close()
	conn.Do("SET", k, v)
	return err
}

func RedisGetValue(k string) (interface{}, error, error) {
	conn := RedisPool.Get()
	err = conn.Err()
	defer conn.Close()
	data, err2 := conn.Do("GET", k)
	return data, err, err2
}

func RedisBool(v interface{}) (bool, error) {
	return redis.Bool(v, err)
}

func RedisString(v interface{}) (string, error) {
	return redis.String(v, err)
}
