/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"encoding/json"
	"github.com/getsentry/raven-go"
	"golang.org/x/net/context"
	"net/http"

	"github.com/SiCo-Ops/Pb"
	"github.com/SiCo-Ops/dao/grpc"
)

type PrivateToken struct {
	ID  string `json:"id"`
	Key string `json:"key"`
}

type AuthenticationToken struct {
	ID        string `json:"id"`
	Signature string `json:"signature"`
}

type TokenRegInfo struct {
	Token string `json:"token"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

func AAAValidateToken(id string, signature string) bool {
	defer func() {
		recover()
		if rcv := recover(); rcv != nil {
			raven.CaptureMessage("Maybe gRPC connect error", nil)
		}
	}()
	cc := rpc.RPCConn(RPCAddr["He"])
	defer cc.Close()
	c := pb.NewAAAPrivateServiceClient(cc)
	in := &pb.AAATokenCall{}
	in.Id = id
	in.Signature = signature
	r, err := c.AuthenticationRPC(context.Background(), in)
	if err != nil {
		return false
	}
	return r.Valid
}

func AAAAuthentication(rw http.ResponseWriter, req *http.Request) {
	data, ok := ValidatePostData(rw, req)
	v := &AuthenticationToken{}
	if ok {
		json.Unmarshal(data, v)
	} else {
		return
	}
	if AAAValidateToken(v.ID, v.Signature) {
		rsp, _ := json.Marshal(&ResponseData{0, "Success"})
		httprsp(rw, rsp)
		return
	}
	rsp, _ := json.Marshal(&ResponseData{2, "AAA_AuthFailed"})
	httprsp(rw, rsp)
}

func AAARegToken(rw http.ResponseWriter, req *http.Request) {
	defer func() {
		recover()
	}()
	data, ok := ValidatePostData(rw, req)
	v := &TokenRegInfo{}
	if ok {
		json.Unmarshal(data, v)
	} else {
		return
	}
	if ValidateOpenToken(v.Token) {
		if v.Email == "" {
			rsp, _ := json.Marshal(&ResponseData{Code: 2, Data: "Must need email"})
			httprsp(rw, rsp)
			return
		}
		cc := rpc.RPCConn(RPCAddr["He"])
		defer cc.Close()
		c := pb.NewAAAPublicServiceClient(cc)
		r, _ := c.GenerateTokenRPC(context.Background(), &pb.AAAGenerateTokenCall{Email: v.Email, Phone: v.Phone})
		if r != nil {
			rsp, _ := json.Marshal(&PrivateToken{ID: r.Id, Key: r.Key})
			httprsp(rw, rsp)
			return
		}
		rsp, _ := json.Marshal(&PrivateToken{"", ""})
		httprsp(rw, rsp)
	} else {
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("Permission Denied"))
	}
}
