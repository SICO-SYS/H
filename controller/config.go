/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"encoding/json"
	"fmt"
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

func configPushRPC(in *pb.ConfigPushCall) *pb.ConfigPushBack {
	defer func() {
		recover()
	}()
	cc := rpc.RPCConn(RPCAddr["B"])
	c := pb.NewConfigServiceClient(cc)
	res, err := c.PushRPC(context.Background(), in)
	if err != nil {
		raven.CaptureError(err, nil)
		return &pb.ConfigPushBack{Code: 1}
	}
	return res
}

func configPullRPC(in *pb.ConfigPullCall) *pb.ConfigPullBack {
	defer func() {
		recover()
	}()
	cc := rpc.RPCConn(RPCAddr["B"])
	c := pb.NewConfigServiceClient(cc)
	res, err := c.PullRPC(context.Background(), in)
	if err != nil {
		raven.CaptureError(err, nil)
		return &pb.ConfigPullBack{Code: 1}
	}
	return res
}

func ConfigPull(rw http.ResponseWriter, req *http.Request) {
	in := &pb.ConfigPullCall{}
	queryString := req.URL.Query()
	id := queryString.Get("id")
	signature := queryString.Get("signature")
	if config.AAAstatus == "active" {
		if !AAAValidateToken(id, signature) {
			rsp, _ := json.Marshal(ResponseErrmsg(1))
			httpResponse("json", rw, rsp)
			return
		}
	}
	in.Id = id
	in.Environment = getRouteName(req, "environment")
	res := configPullRPC(in)
	if res.Code == 1 {
		rsp, _ := json.Marshal(ResponseErrmsg(2))
		httpResponse("json", rw, rsp)
		return
	}
	v := configResponse{}
	json.Unmarshal(res.Params, &v)
	rsp, _ := json.Marshal(&configPullResponse{Code: res.Code, Params: v})
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
	if config.AAAstatus == "active" {
		if !AAAValidateToken(id, signature) {
			rsp, _ := json.Marshal(ResponseErrmsg(1))
			httpResponse("json", rw, rsp)
			return
		}
	}
	in.Id = id
	in.Environment = getRouteName(req, "environment")
	params, _ := json.Marshal(v.Params)
	in.Params = params
	fmt.Println(v.Params)
	res := configPushRPC(in)
	if res.Code == 1 {
		rsp, _ := json.Marshal(ResponseErrmsg(2))
		httpResponse("json", rw, rsp)
		return
	}
	rsp, _ := json.Marshal(&configPushBack{Code: 0})
	httpResponse("json", rw, rsp)
	return
}
