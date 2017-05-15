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

func AAA_PostThirdKeypair(rw http.ResponseWriter, req *http.Request) {
	defer func() {
		recover()
		if rcv := recover(); rcv != nil {
			LogProduce("error", "gRPC connect error")
		}
	}()
	data, ok := AuthPostData(rw, req)
	v := &ThirdKeypair{}
	if ok {
		json.Unmarshal(data, v)
	} else {
		return
	}
	if v.Name == "" || v.APItype == "" || v.ID == "" {
		rsp, _ := json.Marshal(&ResponseData{2, "Missing params, pls follow the guide"})
		httprsp(rw, rsp)
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
		rsp, _ := json.Marshal(&ResponseData{0, "Success"})
		httprsp(rw, rsp)
		return
	}
	rsp, _ := json.Marshal(&ResponseData{2, r.Msg})
	httprsp(rw, rsp)
}

func AAA_Auth(rw http.ResponseWriter, req *http.Request) {
	data, ok := AuthPostData(rw, req)
	v := &AuthToken{}
	if ok {
		json.Unmarshal(data, v)
	} else {
		return
	}
	if AAA(v.Id, v.Signature) {
		rsp, _ := json.Marshal(&ResponseData{0, "Success"})
		httprsp(rw, rsp)
		return
	}
	rsp, _ := json.Marshal(&ResponseData{2, "AAA_AuthFailed"})
	httprsp(rw, rsp)
}

func AAA_GetThirdKey(cloud string, id string, signature string, alias string) (string, string) {
	in := &pb.AAA_ThirdpartyKey{}
	in.Apitoken = &pb.AAA_APIToken{}
	in.Apitoken.Id = id
	in.Apitoken.Signature = signature
	in.Apitype = cloud
	in.Name = alias
	cc := dao.RpcConn(RpcAddr["He"])
	defer cc.Close()
	c := pb.NewAAA_SecretClient(cc)
	res, _ := c.AAA_GetThirdKey(context.Background(), in)
	if res != nil {
		return res.Id, res.Key
	}
	return "", ""
}
