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

// return the sinature is valid
func AAAValidateToken(id string, signature string) (bool, int64) {
	in := &pb.AAATokenCall{}
	in.Id = id
	in.Signature = signature
	cc, err := rpc.Conn(config.RpcHeHost, config.RpcHePort)
	if err != nil {
		raven.CaptureError(err, nil)
		return false, 301
	}
	r := rpc.AAATokenAuthenticationRPC(cc, in)
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
	cc, err := rpc.Conn(config.RpcHeHost, config.RpcHePort)
	if err != nil {
		raven.CaptureError(err, nil)
		httpResponse("json", rw, responseErrMsg(301))
		return
	}
	r := rpc.AAATokenGenerateRPC(cc, in)
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
