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

type PrivateToken struct {
	Id  string `json:"id"`
	Key string `json:"key"`
}

type TokenRegInfo struct {
	Token string `json:"token"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

func GetPublicToken(rw http.ResponseWriter, req *http.Request) {
	key := public.GenerateHexString()
	err := redis.RedisSetWithExpire(redis.OpenAccessPool, key, config.OpenAccess.TokenValid, config.OpenAccess.TokenExpired)
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

// func RegAPIToken(rw http.ResponseWriter, req *http.Request) {
// 	defer func() {
// 		recover()
// 	}()
// 	data, ok := ValidatePostData(rw, req)
// 	v := &TokenRegInfo{}
// 	if ok {
// 		json.Unmarshal(data, v)
// 	} else {
// 		return
// 	}
// 	if ValidateOpenToken(v.Token) {
// 		if v.Email == "" {
// 			rsp, _ := json.Marshal(&ResponseData{Code: 2, Data: "Must need random & email"})
// 			httprsp(rw, rsp)
// 			return
// 		}
// 		cc := rpc.RPCConn(RPCAddr["He"])
// 		defer cc.Close()
// 		c := pb.NewAAAPublicServiceClient(cc)
// 		r, _ := c.(context.Background(), &pb.AAA_RegRequest{Email: v.Email})
// 		if r != nil {
// 			rsp, _ := json.Marshal(&SecretToken{Id: r.Id, Key: r.Key})
// 			httprsp(rw, rsp)
// 			return
// 		}
// 		rsp, _ := json.Marshal(&SecretToken{"", ""})
// 		httprsp(rw, rsp)
// 	} else {
// 		rw.WriteHeader(http.StatusUnauthorized)
// 		rw.Write([]byte("Permission Denied"))
// 	}
// }

func ValidateOpenToken(k string) bool {
	data, err1, err2 := redis.RedisGetWithKey(redis.OpenAccessPool, k)
	if err1 != nil {
		return false
	}
	if err2 != nil {
		return false
	}
	ok, err := redis.RedisValueIsBool(data)
	if err != nil {
		return false
	}
	return ok
}
