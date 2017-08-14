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
	err := redis.SetWithExpire(redis.PublicPool, key, config.OpenAccess.TokenValid, config.OpenAccess.TokenExpired)
	rspdata := &ResponseData{}
	if err != nil {
		raven.CaptureError(err, nil)
		rspdata = ResponseErrmsg(126)
	} else {
		rspdata = &ResponseData{0, &PublicToken{Token: key}}
	}
	rsp, _ := json.Marshal(rspdata)
	httprsp(rw, rsp)
}

func ValidateOpenToken(k string) bool {
	data, err1, err2 := redis.GetWithKey(redis.PublicPool, k)
	if err1 != nil {
		return false
	}
	if err2 != nil {
		return false
	}
	ok, err := redis.ValueIsBool(data)
	if err != nil {
		return false
	}
	return ok
}
