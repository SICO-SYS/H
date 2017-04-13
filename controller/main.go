/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"github.com/SiCo-DevOps/H/cfg"
	"github.com/SiCo-DevOps/H/connection"
)

var (
	RedisPool = connection.RedisPool
	config    = cfg.Config
	err       error
)

type ResponseData struct {
	Code int8        `json:"code"`
	Data interface{} `json:"data"`
}

func ResponseMessage(c int8) string {
	msg := ""
	switch c {
	case 0:
		msg = "[Success] Processed"
	case 1:
		msg = "[Failed] Authentication Failed"
	case 2:
		msg = "[Failed] Authorization Failed"
	case 3:
		msg = "[Failed] Missing Params"
	case 4:
		msg = "[Failed] Params Format Incorrect"
	case 127:
		msg = "[Error] Seems Hack"
	default:
		msg = "[Error] Unknown"
	}
	return msg
}

func RedisSetShort(k string, v interface{}, t int16) error {
	conn := RedisPool.Get()
	err = conn.Close()
	defer conn.Close()
	conn.Do("SET", k, v)
	conn.Do("EXPIRE", k, t)
	return nil
}

func RedisSetLong(k string, v interface{}) error {
	conn := RedisPool.Get()
	defer conn.Close()
	conn.Do("SET", k, v)
	return nil
}
