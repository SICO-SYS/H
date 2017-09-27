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

// B ConfigService PushRPC
func ConfigPushRPC(in *pb.ConfigPushCall) *pb.ConfigPushBack {
	defer func() {
		recover()
	}()
	cc := rpc.RPCConn(RPCAddr["B"])
	c := pb.NewConfigServiceClient(cc)
	r, err := c.PushRPC(context.Background(), in)
	if err != nil {
		raven.CaptureError(err, nil)
		return &pb.ConfigPushBack{Code: 304}
	}
	return r
}

// B ConfigService PullRPC
func ConfigPullRPC(in *pb.ConfigPullCall) *pb.ConfigPullBack {
	defer func() {
		recover()
	}()
	cc := rpc.RPCConn(RPCAddr["B"])
	c := pb.NewConfigServiceClient(cc)
	r, err := c.PullRPC(context.Background(), in)
	if err != nil {
		raven.CaptureError(err, nil)
		return &pb.ConfigPullBack{Code: 304}
	}
	return r
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
	r := ConfigPullRPC(in)
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
	r := ConfigPushRPC(in)
	if r.Code != 0 {
		httpResponse("json", rw, responseErrMsg(r.Code))
		return
	}
	httpResponse("json", rw, responseSuccess())
	return
}
