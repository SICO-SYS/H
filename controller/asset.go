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
	"github.com/SiCo-Ops/public"
)

type AssetTemplate struct {
	PrivateToken AuthenticationToken `json:"token"`
	Name         string              `json:"name"`
	Param        map[string]string   `json:"param"`
}

type AssetSynchronizeRequest struct {
	PrivateToken   AuthenticationToken `json:"token"`
	CloudTokenName string              `json:"name"`
	Region         string              `json:"region"`
}

func AssetCreateTemplate(rw http.ResponseWriter, req *http.Request) {
	v := AssetTemplate{}
	if data, ok := ValidatePostData(rw, req); ok {
		json.Unmarshal(data, &v)
	} else {
		return
	}
	in := &pb.AssetTemplateCall{}
	if config.AAAEnable && !AAAValidateToken(v.PrivateToken.ID, v.PrivateToken.Signature) {
		rsp, _ := json.Marshal(ResponseErrmsg(1))
		httprsp(rw, rsp)
		return
	}
	in.Id = v.PrivateToken.ID
	in.Name = v.Name
	in.Param = v.Param
	cc := rpc.RPCConn(RPCAddr["Be"])
	defer cc.Close()
	c := pb.NewTemplateServiceClient(cc)
	res, err := c.CreateRPC(context.Background(), in)
	if err != nil {
		raven.CaptureError(err, nil)
	}
	if res.Code == 0 {
		rsp, _ := json.Marshal(&ResponseData{0, "Success add template"})
		httprsp(rw, rsp)
		return
	}
	rsp, _ := json.Marshal(res)
	httprsp(rw, rsp)
	return
}

func AssetSynchronizeRPC(in *pb.AssetSynchronizeCall) *pb.AssetMsgBack {
	defer func() {
		recover()
	}()
	cc := rpc.RPCConn(RPCAddr["Be"])
	defer cc.Close()
	c := pb.NewAssetServiceClient(cc)
	res, err := c.SynchronizeRPC(context.Background(), in)
	if err != nil {
		raven.CaptureError(err, nil)
		return &pb.AssetMsgBack{Code: -1, Msg: ""}
	}
	return res
}

func AssetSynchronize(rw http.ResponseWriter, req *http.Request) {
	defer func() {
		recover()
	}()
	data, ok := ValidatePostData(rw, req)
	if !ok {
		return
	}

	v := &AssetSynchronizeRequest{}
	json.Unmarshal(data, v)

	if config.AAAEnable && !AAAValidateToken(v.PrivateToken.ID, v.PrivateToken.Signature) {
		rsp, _ := json.Marshal(ResponseErrmsg(1))
		httprsp(rw, rsp)
		return
	}

	cloud := getRouteName(req, "cloud")
	service := getRouteName(req, "service")
	action, ok := actionMap(cloud, service, "DescribeInstances")
	if !ok {
		rsp, _ := json.Marshal(ResponseErrmsg(4))
		httprsp(rw, rsp)
		return
	}

	cloudTokenID, cloudTokenKey = CloudTokenGet(v.PrivateToken.ID, cloud, v.CloudTokenName)

	var moreSource bool = true
	for i := 0; moreSource; i++ {
		in := &pb.CloudAPICall{Cloud: cloud, Service: service, Action: action, Region: v.Region, CloudId: cloudTokenID, CloudKey: cloudTokenKey}
		in.Params = make(map[string]string)
		in.Params["Limit"] = "100"
		in.Params["Offset"] = public.Int2String(i * 100)
		res := CloudAPIRPC(in)
		assetRes := AssetSynchronizeRPC(&pb.AssetSynchronizeCall{Id: cloudTokenID, Cloud: cloud, Service: service, Data: res.Data})
		if assetRes.Code == -1 {
			rsp, _ := json.Marshal(ResponseErrmsg(2))
			httpResponse("json", rw, rsp)
			moreSource = false
			return
		}

		if assetRes.Code == 1 {
			moreSource = false
		}
	}

	if !moreSource {
		rsp, _ := json.Marshal(&ResponseData{Code: 0, Data: "success"})
		httpResponse("json", rw, rsp)
	}
}
