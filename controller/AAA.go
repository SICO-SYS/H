/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"encoding/json"
	"golang.org/x/net/context"
	"io/ioutil"
	"net/http"

	"github.com/SiCo-DevOps/Pb"
	"github.com/SiCo-DevOps/dao"
	. "github.com/SiCo-DevOps/log"
)

type AuthToken struct {
	Key   string `json:"key"`
	Token string `json:"token"`
}

type ThirdKeypair struct {
	Auth    AuthToken `json:"auth"`
	APItype string    `json:"apitype"`
	ID      string    `json:"id"`
	Token   string    `json:"token"`
}

func AAA(k string, s string) bool {
	defer func() {
		recover()
		if rcv := recover(); rcv != nil {
			LogProduce("error", "gRPC connect error")
		}
	}()
	cc := dao.RpcConn("He")
	defer cc.Close()
	c := pb.NewAAA_SecretClient(cc)
	in := &pb.AAA_APIToken{}
	in.Key = Sha256Encrypt(k)
	in.Token = s
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
	header := req.Header.Get("Content-Type")
	body, _ := ioutil.ReadAll(req.Body)
	req.Body.Close()
	data := &ThirdKeypair{}
	json.Unmarshal(body, data)
	if header != "application/json" {
		rsp := &ResponseData{1, "request must follow application/json"}
		rspdata, _ := json.Marshal(rsp)
		rw.Write(rspdata)
		return
	}
	defer func() {
		recover()
		if rcv := recover(); rcv != nil {
			LogProduce("error", "gRPC connect error")
		}
	}()
	cc := dao.RpcConn("He")
	defer cc.Close()
	c := pb.NewAAA_SecretClient(cc)
	in := &pb.AAA_ThirdpartyKey{}
	in.Apitoken = &pb.AAA_APIToken{}
	in.Apitoken.Key = Sha256Encrypt(data.Auth.Key)
	in.Apitoken.Token = data.Auth.Token
	in.Apitype = data.APItype
	in.Id = data.ID
	in.Key = data.Token
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
	header := req.Header.Get("Content-Type")
	body, _ := ioutil.ReadAll(req.Body)
	req.Body.Close()
	keypair := &AuthToken{}
	json.Unmarshal(body, keypair)
	if header != "application/json" {
		rsp := &ResponseData{1, "request must follow application/json"}
		rspdata, _ := json.Marshal(rsp)
		rw.Header().Add("Content-Type", "application/json")
		rw.Write(rspdata)
		return
	}
	if AAA(keypair.Key, keypair.Token) {
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
