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
)

type AuthToken struct {
	Id        string `json:"id"`
	Signature string `json:"signature"`
}

type ThirdKeypair struct {
	Auth    AuthToken `json:"auth"`
	APItype string    `json:"apitype"`
	Name    string    `json:"name"`
	ID      string    `json:"id"`
	Key     string    `json:"key"`
}

func AAA(k string, s string) bool {
	defer func() {
		recover()
		if rcv := recover(); rcv != nil {
			LogProduce("error", "gRPC connect error")
		}
	}()
	cc := dao.RpcConn(RpcAddr["He"])
	defer cc.Close()
	c := pb.NewAAA_SecretClient(cc)
	in := &pb.AAA_APIToken{}
	in.Id = k
	in.Signature = s
	r, err := c.AAA_Auth(context.Background(), in)
	if err != nil {
		LogErrMsg(50, "controller.AAA")
		return false
	}

	if r.Code == 0 {
		return true
	}
	return false
}

func PostThirdKeypair(rw http.ResponseWriter, req *http.Request) {
	defer func() {
		recover()
		if rcv := recover(); rcv != nil {
			LogProduce("error", "gRPC connect error")
		}
	}()
	data, ok := AuthPostData(req)
	v := &ThirdKeypair{}
	if ok {
		json.Unmarshal(data, v)
	} else {
		rsp := &ResponseData{2, "request must follow application/json"}
		rspdata, _ := json.Marshal(rsp)
		rw.Write(rspdata)
		return
	}
	if v.Name == "" || v.APItype == "" || v.ID == "" || v.Key == "" {
		rsp := &ResponseData{2, "Missing params, pls follow the guide"}
		rspdata, _ := json.Marshal(rsp)
		rw.Write(rspdata)
		return
	}
	cc := dao.RpcConn(RpcAddr["He"])
	defer cc.Close()
	c := pb.NewAAA_SecretClient(cc)
	in := &pb.AAA_ThirdpartyKey{}
	in.Apitoken = &pb.AAA_APIToken{}
	in.Apitoken.Id = v.Auth.Id
	in.Apitoken.Signature = v.Auth.Signature
	in.Apitype = v.APItype
	in.Name = v.Name
	in.Id = v.ID
	in.Key = v.Key
	r, err := c.AAA_ThirdKeypair(context.Background(), in)
	if err != nil {
		LogErrMsg(50, "controller.PostThirdKeypair")
	}

	if r.Code == 0 {
		rsp := &ResponseData{0, "Success"}
		rspdata, _ := json.Marshal(rsp)
		rw.Header().Add("Content-Type", "application/json")
		rw.Write(rspdata)
		return
	}
	rsp := &ResponseData{2, r.Msg}
	rspdata, _ := json.Marshal(rsp)
	rw.Header().Add("Content-Type", "application/json")
	rw.Write(rspdata)
}

func AAA_Auth(rw http.ResponseWriter, req *http.Request) {
	data, ok := AuthPostData(req)
	v := &AuthToken{}
	if ok {
		json.Unmarshal(data, v)

	} else {
		rsp := &ResponseData{2, "request must follow application/json"}
		rspdata, _ := json.Marshal(rsp)
		rw.Header().Add("Content-Type", "application/json")
		rw.Write(rspdata)
		return
	}
	if AAA(v.Id, v.Signature) {
		rsp := &ResponseData{0, "Success"}
		rspdata, _ := json.Marshal(rsp)
		rw.Header().Add("Content-Type", "application/json")
		rw.Write(rspdata)
		return
	}
	rsp := &ResponseData{2, "AAA_AuthFailed"}
	rspdata, _ := json.Marshal(rsp)
	rw.Header().Add("Content-Type", "application/json")
	rw.Write(rspdata)

}
