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
	"io/ioutil"
	"net/http"

	"github.com/SiCo-Ops/Pb"
	"github.com/SiCo-Ops/dao/grpc"
	// "github.com/SiCo-Ops/dao/mongo"
)

var (
	cloudTokenID  string
	cloudTokenKey string
	cloudRegion   string
	cloudService  string
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

func CloudTokenRegistry(rw http.ResponseWriter, req *http.Request) {
	defer func() {
		recover()
		if rcv := recover(); rcv != nil {
			raven.CaptureMessage("controller.CloudTokenRegistry", nil)
		}
	}()
	data, ok := ValidatePostData(rw, req)
	v := &ThirdToken{}
	if ok {
		json.Unmarshal(data, v)
	} else {
		return
	}
	if v.Name == "" || v.Cloud == "" || v.ID == "" {
		rsp, _ := json.Marshal(ResponseErrmsg(2))
		httprsp(rw, rsp)
		return
	}
	if config.AAAEnable && !AAAValidateToken(v.PrivateToken.ID, v.PrivateToken.Signature) {
		rsp, _ := json.Marshal(ResponseErrmsg(1))
		httprsp(rw, rsp)
		return
	}
	cc := rpc.RPCConn(RPCAddr["Li"])
	defer cc.Close()
	c := pb.NewCloudTokenServiceClient(cc)
	in := &pb.CloudTokenCall{}
	in.Cloud = v.Cloud
	in.Name = v.Name
	in.Id = v.ID
	in.Key = v.Key
	in.AAATokenID = v.PrivateToken.ID
	r, err := c.TokenSet(context.Background(), in)
	if err != nil {
		raven.CaptureError(err, nil)
	}
	if r.Id != "" {
		rsp, _ := json.Marshal(&ResponseData{0, "Success"})
		httprsp(rw, rsp)
		return
	}
	rsp, _ := json.Marshal(ResponseErrmsg(2))
	httprsp(rw, rsp)
}

func CloudTokenGet(id string, cloud string, name string) (string, string) {
	in := &pb.CloudTokenCall{}
	in.AAATokenID = id
	in.Cloud = cloud
	in.Name = name
	cc := rpc.RPCConn(RPCAddr["Li"])
	defer cc.Close()
	c := pb.NewCloudTokenServiceClient(cc)
	res, err := c.TokenGet(context.Background(), in)
	if err != nil {
		raven.CaptureError(err, nil)
	}
	if res.Id != "" {
		return res.Id, res.Key
	}
	return "", ""
}

func CloudServiceIsSupport(cloud string, service string) bool {
	d, err := ioutil.ReadFile("cloud.json")
	if err != nil {
		raven.CaptureError(err, nil)
		return false
	}
	var v map[string][]string
	json.Unmarshal(d, &v)
	if value, ok := v[cloud]; ok {
		for _, v := range value {
			if v == service {
				return true
			}
		}
		return false
	}
	return false
}

func CloudAPICall(rw http.ResponseWriter, req *http.Request) {
	data, ok := ValidatePostData(rw, req)
	if !ok {
		return
	}
	v := &CloudAPIRequest{}
	json.Unmarshal(data, v)

	if !AAAValidateToken(v.PrivateToken.ID, v.PrivateToken.Signature) {
		rsp, _ := json.Marshal(ResponseErrmsg(2))
		httprsp(rw, rsp)
		return
	}

	cloud := getRouteName(req, "cloud")
	service := getRouteName(req, "service")
	action, ok := actionMap(cloud, service, v.Action)
	if !ok {
		rsp, _ := json.Marshal(ResponseErrmsg(4))
		httprsp(rw, rsp)
		return
	}

	cloudTokenID, cloudTokenKey = CloudTokenGet(v.PrivateToken.ID, cloud, v.CloudTokenName)

	in := &pb.CloudAPICall{Cloud: cloud, Service: service, Action: action, Region: v.Region, CloudId: cloudTokenID, CloudKey: cloudTokenKey}
	in.Params = v.Param
	cc := rpc.RPCConn(RPCAddr["Li"])
	defer cc.Close()
	c := pb.NewCloudAPIServiceClient(cc)
	var res *pb.CloudAPIBack
	res, _ = c.RequestRPC(context.Background(), in)
	if res.Code == 0 {
		rsp := res.Data
		httprsp(rw, rsp)
		return
	}
	rsp, _ := json.Marshal(res)
	httprsp(rw, rsp)
}

func CloudAPICallRaw(rw http.ResponseWriter, req *http.Request) {
	data, ok := ValidatePostData(rw, req)
	if !ok {
		return
	}
	v := &CloudAPIRawRequest{}
	json.Unmarshal(data, v)
	if !ValidateOpenToken(v.Token) {
		rsp, _ := json.Marshal(ResponseErrmsg(5))
		httprsp(rw, rsp)
		return
	}

	cloud := getRouteName(req, "cloud")
	service := getRouteName(req, "service")

	in := &pb.CloudAPICall{Cloud: cloud, Service: service, Action: v.Action, Region: v.Region, CloudId: v.CloudTokenID, CloudKey: v.CloudTokenKey}
	in.Params = v.Param
	cc := rpc.RPCConn(RPCAddr["Li"])
	defer cc.Close()
	c := pb.NewCloudAPIServiceClient(cc)
	var res *pb.CloudAPIBack
	res, _ = c.RequestRPC(context.Background(), in)
	if res.Code == 0 {
		rsp := res.Data
		httprsp(rw, rsp)
		return
	}
	rsp, _ := json.Marshal(res)
	httprsp(rw, rsp)

}
