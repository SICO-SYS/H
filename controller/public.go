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
	err := redis.SetWithExpire(publicPool, key, config.PublicTokenStatus, public.String2Int(config.PublicTokenExpire))
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
	data, err := redis.ExpiredAfterGetWithKey(publicPool, k)
	if err != nil {
		return false
	}
	ok, err := redis.ValueIsString(data)
	if err != nil {
		return false
	}
	if ok != "active" {
		return false
	}
	return true
}
