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

type configResponse map[string]interface{}

type configPushCall struct {
	PrivateToken AuthenticationToken `json:"token"`
	Params       configResponse      `json:"params"`
}

type configPushBack struct {
	Code int64 `json:"code"`
}

type configPullResponse struct {
	Code   int64          `json:"code"`
	Params configResponse `json:"params"`
}

func ConfigPull(rw http.ResponseWriter, req *http.Request) {
	in := &pb.ConfigPullCall{}
	queryString := req.URL.Query()
	id := queryString.Get("id")
	signature := queryString.Get("signature")

	isPrivateTokenValid, errcode := AAAValidateToken(id, signature)
	if errcode != 0 {
		httpResponse("json", rw, responseErrMsg(errcode))
		return
	}
	if !isPrivateTokenValid {
		httpResponse("json", rw, responseErrMsg(1000))
		return
	}

	in.Id = id
	in.Environment = getRouteName(req, "environment")
	cc, err := rpc.Conn(config.RpcBHost, config.RpcBPort)
	if err != nil {
		raven.CaptureError(err, nil)
		httpResponse("json", rw, responseErrMsg(304))
		return
	}
	r := rpc.ConfigPullRPC(cc, in)
	if r.Code != 0 {
		httpResponse("json", rw, responseErrMsg(r.Code))
		return
	}
	v := configResponse{}
	json.Unmarshal(r.Params, &v)
	rsp, _ := json.Marshal(&responseData{Code: r.Code, Data: v})
	httpResponse("json", rw, rsp)
	return
}

func ConfigPush(rw http.ResponseWriter, req *http.Request) {
	in := &pb.ConfigPushCall{}
	data, ok := ValidatePostData(rw, req)
	v := &configPushCall{}
	if ok {
		json.Unmarshal(data, v)
	} else {
		return
	}
	id := v.PrivateToken.ID
	signature := v.PrivateToken.Signature

	isPrivateTokenValid, errcode := AAAValidateToken(id, signature)
	if errcode != 0 {
		httpResponse("json", rw, responseErrMsg(errcode))
		return
	}
	if !isPrivateTokenValid {
		httpResponse("json", rw, responseErrMsg(1000))
		return
	}

	in.Id = id
	in.Environment = getRouteName(req, "environment")
	params, _ := json.Marshal(v.Params)
	in.Params = params
	cc, err := rpc.Conn(config.RpcBHost, config.RpcBPort)
	if err != nil {
		raven.CaptureError(err, nil)
		httpResponse("json", rw, responseErrMsg(304))
		return
	}
	r := rpc.ConfigPushRPC(cc, in)
	if r.Code != 0 {
		httpResponse("json", rw, responseErrMsg(r.Code))
		return
	}
	httpResponse("json", rw, responseSuccess())
	return
}
