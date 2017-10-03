/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"encoding/json"
	"github.com/getsentry/raven-go"
	"net/http"

	"github.com/SiCo-Ops/dao/redis"
	"github.com/SiCo-Ops/public"
)

type PublicToken struct {
	Token string `json:"token"`
}

func PublicGenerateToken(rw http.ResponseWriter, req *http.Request) {
	key := public.GenerateHexString()
	if config.PublicTokenStatus != "active" {
		httpResponse("json", rw, responseErrMsg(7))
		return
	}
	err := redis.Set(publicPool, key, config.PublicTokenStatus, int64(public.StringToInt(config.PublicTokenExpire)))
	if err != nil {
		raven.CaptureError(err, nil)
		httpResponse("json", rw, responseErrMsg(101))
		return
	}
	rsp, _ := json.Marshal(&responseData{Code: 0, Data: &PublicToken{Token: key}})
	httpResponse("json", rw, rsp)
	return
}

func PublicValidateToken(k string) (bool, int64) {
	if config.PublicTokenStatus != "active" {
		return false, 7
	}
	data, err := redis.ExpiredAfterGet(publicPool, k)
	if err != nil {
		return false, 101
	}
	ok, err := redis.ValueIsString(data)
	if err != nil {
		return false, 0
	}
	if ok != "active" {
		return false, 7
	}
	return true, 0
}
