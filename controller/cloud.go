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
	// "io/ioutil"
	"net/http"

	"github.com/SiCo-Ops/Pb"
	"github.com/SiCo-Ops/dao/grpc"
	"github.com/SiCo-Ops/public"
)

type ThirdToken struct {
	PrivateToken AuthenticationToken `json:"token"`
	Cloud        string              `json:"cloud"`
	Name         string              `json:"name"`
	ID           string              `json:"id"`
	Key          string              `json:"key"`
}

type CloudAPIRequest struct {
	PrivateToken   AuthenticationToken `json:"token"`
	CloudTokenName string              `json:"name"`
	Region         string              `json:"region"`
	Action         string              `json:"action"`
	Param          map[string]string   `json:"params"`
}

type CloudAPIRawRequest struct {
	Token         string            `json:"token"`
	CloudTokenID  string            `json:"cloudid"`
	CloudTokenKey string            `json:"cloudkey"`
	Region        string            `json:"region"`
	Action        string            `json:"action"`
	Param         map[string]string `json:"params"`
}

type CloudAPIResponse struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

// Li APIservice RequestRPC
func CloudAPIRequestRPC(in *pb.CloudAPICall) *pb.CloudAPIBack {
	defer func() {
		recover()
	}()
	cc := rpc.RPCConn(RPCAddr["Li"])
	defer cc.Close()
	c := pb.NewCloudAPIServiceClient(cc)
	r, err := c.RequestRPC(context.Background(), in)
	if err != nil {
		raven.CaptureError(err, nil)
		return &pb.CloudAPIBack{Code: 302}
	}
	return r
}

// Li TokenService SetRPC
func CloudTokenSetRPC(in *pb.CloudTokenCall) *pb.CloudTokenBack {
	defer func() {
		recover()
	}()
	cc := rpc.RPCConn(RPCAddr["Li"])
	defer cc.Close()
	c := pb.NewCloudTokenServiceClient(cc)
	r, err := c.SetRPC(context.Background(), in)
	if err != nil {
		raven.CaptureError(err, nil)
		return &pb.CloudTokenBack{Code: 302}
	}
	return r
}

// Li TokenService GetRPC
func CloudTokenGetRPC(in *pb.CloudTokenCall) *pb.CloudTokenBack {
	defer func() {
		recover()
	}()
	cc := rpc.RPCConn(RPCAddr["Li"])
	defer cc.Close()
	c := pb.NewCloudTokenServiceClient(cc)
	r, err := c.GetRPC(context.Background(), in)
	if err != nil {
		raven.CaptureError(err, nil)
		return &pb.CloudTokenBack{Code: 302}
	}
	return r
}

func CloudAPICallForLoop(in *pb.CloudAPICall, nextToken string, page, totalCount int64) (out *pb.CloudAPICall, isLoop bool) {
	in.Params = make(map[string]string)
	var requestCount int64 = 100
	switch in.Cloud {
	case "qcloud":
		in.Params["Limit"] = public.Int64ToString(requestCount)
		in.Params["Offset"] = public.Int64ToString(page * requestCount)
		if page != 0 && float64(page+1) >= float64(totalCount)/float64(requestCount) {
			return nil, false
		}
	case "aliyun":
		in.Params["PageNumber"] = public.Int64ToString(page + 1)
		in.Params["PageSize"] = public.Int64ToString(requestCount)
		if page != 0 && float64(page+1) >= float64(totalCount)/float64(requestCount) {
			return nil, false
		}
	case "aws":
		in.Params["MaxResults"] = public.Int64ToString(requestCount)
		in.Params["NextToken"] = nextToken
		if page != 0 && nextToken == "" {
			return nil, false
		}
	default:
		return nil, false
	}
	return in, true
}

func cloudTokenGet(id string, cloud string, name string) (string, string, int64) {
	in := &pb.CloudTokenCall{}
	in.AAATokenID = id
	in.Cloud = cloud
	in.Name = name
	r := CloudTokenGetRPC(in)
	if r.Code != 0 {
		return "", "", r.Code
	}
	if r.Id == "" {
		return "", "", 0
	}
	return r.Id, r.Key, 0
}

// func CloudServiceIsSupport(cloud string, service string) (bool, int64) {
// 	d, err := ioutil.ReadFile("cloud.json")
// 	if err != nil {
// 		return false, 9
// 	}
// 	var v map[string][]string
// 	json.Unmarshal(d, &v)
// 	if value, ok := v[cloud]; ok {
// 		for _, v := range value {
// 			if v == service {
// 				return true, 0
// 			}
// 		}
// 		return false, 0
// 	}
// 	return false, 0
// }

// POST /v1/cloud/token
func CloudTokenRegistry(rw http.ResponseWriter, req *http.Request) {
	data, ok := ValidatePostData(rw, req)
	if !ok {
		return
	}
	v := &ThirdToken{}
	json.Unmarshal(data, v)
	if v.Name == "" || v.Cloud == "" || v.ID == "" {
		httpResponse("json", rw, responseErrMsg(2000))
		return
	}
	isPrivateTokenValid, errcode := AAAValidateToken(v.PrivateToken.ID, v.PrivateToken.Signature)
	if errcode != 0 {
		httpResponse("json", rw, responseErrMsg(errcode))
		return
	}
	if config.AAAstatus == "active" && !isPrivateTokenValid {
		httpResponse("json", rw, responseErrMsg(1000))
		return
	}
	in := &pb.CloudTokenCall{}
	in.Cloud = v.Cloud
	in.Name = v.Name
	in.Id = v.ID
	in.Key = v.Key
	in.AAATokenID = v.PrivateToken.ID
	r := CloudTokenSetRPC(in)
	if r.Code != 0 {
		httpResponse("json", rw, responseErrMsg(r.Code))
		return
	}
	httpResponse("json", rw, responseSuccess())
	return
}

func CloudAPICall(rw http.ResponseWriter, req *http.Request) {
	data, ok := ValidatePostData(rw, req)
	if !ok {
		return
	}
	v := &CloudAPIRequest{}
	json.Unmarshal(data, v)

	isPrivateTokenValid, errcode := AAAValidateToken(v.PrivateToken.ID, v.PrivateToken.Signature)
	if errcode != 0 {
		httpResponse("json", rw, responseErrMsg(errcode))
		return
	}
	if !isPrivateTokenValid {
		httpResponse("json", rw, responseErrMsg(1000))
		return
	}

	cloud := getRouteName(req, "cloud")
	service := getRouteName(req, "service")
	action := v.Action

	cloudTokenID, cloudTokenKey, errcode := cloudTokenGet(v.PrivateToken.ID, cloud, v.CloudTokenName)
	if errcode != 0 {
		httpResponse("json", rw, responseErrMsg(errcode))
		return
	}

	in := &pb.CloudAPICall{Cloud: cloud, Service: service, Action: action, Region: v.Region, CloudId: cloudTokenID, CloudKey: cloudTokenKey}
	in.Params = v.Param
	r := CloudAPIRequestRPC(in)
	if r.Code != 0 {
		httpResponse("json", rw, responseErrMsg(r.Code))
		return
	}
	if cloud == "aws" {
		httpResponse("xml", rw, r.Data)
	} else {
		httpResponse("json", rw, r.Data)
	}
}

func CloudAPICallRaw(rw http.ResponseWriter, req *http.Request) {
	data, ok := ValidatePostData(rw, req)
	if !ok {
		return
	}
	v := &CloudAPIRawRequest{}
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

	cloud := getRouteName(req, "cloud")
	service := getRouteName(req, "service")

	in := &pb.CloudAPICall{Cloud: cloud, Service: service, Action: v.Action, Region: v.Region, CloudId: v.CloudTokenID, CloudKey: v.CloudTokenKey}
	in.Params = v.Param
	r := CloudAPIRequestRPC(in)
	if r.Code != 0 {
		httpResponse("json", rw, responseErrMsg(r.Code))
		return
	}
	if cloud == "aws" {
		httpResponse("xml", rw, r.Data)
	} else {
		httpResponse("json", rw, r.Data)
	}
}
