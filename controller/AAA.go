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
	Token     string `json:"token"`
	ID        string `json:"id"`
	Signature string `json:"signature"`
}

type TokenRegInfo struct {
	Token string `json:"token"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// He PublicService GenerateTokenRPC
func AAAPublicGenerateTokenRPC(in *pb.AAAGenerateTokenCall) *pb.AAAGenerateTokenBack {
	defer func() {
		recover()
	}()
	cc := rpc.RPCConn(RPCAddr["He"])
	defer cc.Close()
	c := pb.NewAAAPublicServiceClient(cc)
	r, err := c.GenerateTokenRPC(context.Background(), in)
	if err != nil {
		raven.CaptureError(err, nil)
		return &pb.AAAGenerateTokenBack{Code: 301}
	}
	return r
}

// He PrivateService AuthenticationRPC
func AAAPrivateAuthenticationRPC(in *pb.AAATokenCall) *pb.AAATokenBack {
	defer func() {
		recover()
	}()
	cc := rpc.RPCConn(RPCAddr["He"])
	defer cc.Close()
	c := pb.NewAAAPrivateServiceClient(cc)
	r, err := c.AuthenticationRPC(context.Background(), in)
	if err != nil {
		raven.CaptureError(err, nil)
		return &pb.AAATokenBack{Code: 301}
	}
	return r
}

// He PrivateService AuthorizationRPC
func AAAPrivateAuthorizationRPC(in *pb.AAAServiceCall) *pb.AAAServiceBack {
	defer func() {
		recover()
	}()
	cc := rpc.RPCConn(RPCAddr["He"])
	defer cc.Close()
	c := pb.NewAAAPrivateServiceClient(cc)
	r, err := c.AuthorizationRPC(context.Background(), in)
	if err != nil {
		raven.CaptureError(err, nil)
		return &pb.AAAServiceBack{Code: 301}
	}
	return r
}

// He PrivateService AuthorizationRPC
func AAAPrivateAccountingRPC(in *pb.AAAEventCall) *pb.AAAEventBack {
	defer func() {
		recover()
	}()
	cc := rpc.RPCConn(RPCAddr["He"])
	defer cc.Close()
	c := pb.NewAAAPrivateServiceClient(cc)
	r, err := c.AccountingRPC(context.Background(), in)
	if err != nil {
		raven.CaptureError(err, nil)
		return &pb.AAAEventBack{Code: 301}
	}
	return r

}

// return the sinature is valid
func AAAValidateToken(id string, signature string) (bool, int64) {
	in := &pb.AAATokenCall{}
	in.Id = id
	in.Signature = signature
	r := AAAPrivateAuthenticationRPC(in)
	if r.Code != 0 {
		return false, r.Code
	}
	return r.IsValid, 0
}

// GET /v1/AAA/token
func AAAGenerateToken(rw http.ResponseWriter, req *http.Request) {
	defer func() {
		recover()
	}()
	data, ok := ValidatePostData(rw, req)
	if !ok {
		return
	}
	v := &TokenRegInfo{}
	json.Unmarshal(data, v)
	isPublicTokenValid, errcode := PublicValidateToken(v.Token)
	if errcode != 0 {
		httpResponse("json", rw, responseErrMsg(errcode))
		return
	}
	if !isPublicTokenValid {
		httpResponse("json", rw, responseErrMsg(8))
		return
	}

	if v.Email == "" && v.Phone == "" {
		httpResponse("json", rw, responseErrMsg(1003))
		return
	}
	in := &pb.AAAGenerateTokenCall{Email: v.Email, Phone: v.Phone}
	r := AAAPublicGenerateTokenRPC(in)
	if r.Code != 0 {
		httpResponse("json", rw, responseErrMsg(r.Code))
		return
	}
	rsp, _ := json.Marshal(&responseData{Code: 0, Data: &PrivateToken{ID: r.Id, Key: r.Key}})
	httpResponse("json", rw, rsp)
	return
}

// POST /v1/AAA/authentication
func AAAAuthentication(rw http.ResponseWriter, req *http.Request) {
	data, ok := ValidatePostData(rw, req)
	if !ok {
		return
	}
	v := &AuthenticationToken{}
	json.Unmarshal(data, v)
	isPublicTokenValid, errcode := PublicValidateToken(v.Token)
	if errcode != 0 {
		httpResponse("json", rw, responseErrMsg(errcode))
		return
	}
	if !isPublicTokenValid {
		httpResponse("json", rw, responseErrMsg(7))
		return
	}

	isPrivateTokenValid, errcode := AAAValidateToken(v.ID, v.Signature)
	if errcode != 0 {
		httpResponse("json", rw, responseErrMsg(errcode))
		return
	}
	if !isPrivateTokenValid {
		httpResponse("json", rw, responseErrMsg(1000))
		return
	}
	httpResponse("json", rw, responseSuccess())
	return
}
