/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"encoding/json"
	"golang.org/x/net/context"
	"net/http"

	"github.com/SiCo-DevOps/Pb"
	"github.com/SiCo-DevOps/dao"
	. "github.com/SiCo-DevOps/log"
	"github.com/SiCo-DevOps/public"
)

type OpenToken struct {
	Token string `json:"token"`
}

type SecretToken struct {
	Id     string `json:"id"`
	Key    string `json:"key"`
	Random string `json:"random"`
}

type RegUser struct {
	Token  string `json:"token"`
	Random string `json:"random"`
	Email  string `json:"email"`
}

func GetOpenToken(rw http.ResponseWriter, req *http.Request) {
	key := public.GenHexString()
	err = dao.RedisSetShort(key, config.OpenAccess.TokenValid, config.OpenAccess.TokenExpired)
	rspdata := &ResponseData{}
	if err != nil {
		rspdata = ResponseErrmsg(126)
		LogErrMsg(10, "controller.GetOpenToken")
	} else {
		rspdata = &ResponseData{0, &OpenToken{Token: key}}
	}
	rsp, _ := json.Marshal(rspdata)
	rw.Header().Add("Content-Type", "application/json")
	rw.Write(rsp)
}

func RegAPIToken(rw http.ResponseWriter, req *http.Request) {
	defer func() {
		recover()
		LogErrMsg(5, "controller.RegAPIToken")
	}()
	data, ok := AuthPostData(req)
	v := &RegUser{}
	if ok {
		json.Unmarshal(data, v)
	} else {
		rsp, _ := json.Marshal(&ResponseData{2, "request must follow application/json"})
		rw.Header().Add("Content-Type", "application/json")
		rw.Write(rsp)
		return
	}
	if AuthOpenToken(v.Token) {
		if v.Random == "" || v.Email == "" {
			rsp, _ := json.Marshal(&ResponseData{Code: 2, Data: "Must need random & email"})
			rw.Header().Add("Content-Type", "application/json")
			rw.Write(rsp)
			return
		}
		cc := dao.RpcConn(RpcAddr["He"])
		defer cc.Close()
		c := pb.NewAAA_OpenClient(cc)
		r, err := c.AAA_RegUser(context.Background(), &pb.AAA_RegRequest{Random: v.Random, Email: v.Email})
		if err != nil {
			LogErrMsg(50, "controller.RegAPIToken")
		}
		if r != nil {
			rsp, _ := json.Marshal(&SecretToken{Id: r.Id, Key: r.Key, Random: v.Random})
			rw.Header().Add("Content-Type", "application/json")
			rw.Write(rsp)
			return
		}
		rsp, _ := json.Marshal(&SecretToken{"", "", v.Random})
		rw.Header().Add("Content-Type", "application/json")
		rw.Write(rsp)
	} else {
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("Permission Denied"))
	}
}

func AuthOpenToken(k string) bool {
	data, err1, err2 := dao.RedisGetValue(k)
	if err1 != nil {
		LogErrMsg(1, "controller.AuthOpenToken")
		return false
	}
	if err2 != nil {
		LogErrMsg(19, "controller.AuthOpenToken")
		return false
	}
	ok, err := dao.RedisBool(data)
	if err != nil {
		LogErrMsg(11, "controller.AuthOpenToken")
		return false
	}
	return ok
}
