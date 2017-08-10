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
	// "github.com/SiCo-Ops/dao/mongo"
)

// var (
// 	cloud_id     string
// 	cloud_key    string
// 	cloud_region string
// 	cloud_action string
// )

type ThirdToken struct {
	PrivateToken AuthenticationToken `json:"token"`
	Cloud        string              `json:"cloud"`
	Name         string              `json:"name"`
	ID           string              `json:"id"`
	Key          string              `json:"key"`
}

// type Cloud_Req struct {
// 	Auth   AuthToken         `json:"auth"`
// 	Region string            `json:"region"`
// 	Action string            `json:"action"`
// 	Name   string            `json:"name"`
// 	Param  map[string]string `json:"params"`
// }

// type Cloud_Res struct {
// 	Code int64  `json:"code"`
// 	Msg  string `json:"msg"`
// 	Data string `json:"data"`
// }

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
	if !AAAValidateToken(v.PrivateToken.ID, v.PrivateToken.Signature) {
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
	r, err := c.TokenSet(context.Background(), in)
	if err != nil {
		raven.CaptureError(err, nil)
	}
	if r.Id == "" {
		rsp, _ := json.Marshal(&ResponseData{0, "Success"})
		httprsp(rw, rsp)
		return
	}
	rsp, _ := json.Marshal(ResponseErrmsg(2))
	httprsp(rw, rsp)
}

// func CloudTokenFind(cloud string, id string, signature string, alias string) (string, string) {
// 	in := &pb.AAA_ThirdpartyKey{}
// 	in.Apitoken = &pb.AAA_APIToken{}
// 	in.Apitoken.ID = id
// 	in.Apitoken.Signature = signature
// 	in.Apitype = cloud
// 	in.Name = alias
// 	cc := dao.RpcConn(RPCAddr["He"])
// 	defer cc.Close()
// 	c := pb.NewAAA_SecretClient(cc)
// 	res, _ := c.AAA_GetThirdKey(context.Background(), in)
// 	if res != nil {
// 		return res.ID, res.Key
// 	}
// 	return "", ""
// }

// func AuthBsns(cloud string, bsns string) bool {
// 	d, err := ioutil.ReadFile("cloud.json")
// 	if err != nil {
// 		LogFatalMsg(0, "controller.AuthBsns")
// 	}
// 	var v map[string][]string
// 	json.Unmarshal(d, &v)
// 	if value, ok := v[cloud]; ok {
// 		for _, v := range value {
// 			if v == bsns {
// 				return true
// 			}
// 		}
// 		return false
// 	}
// 	return false
// }

// func Cloud_CommonCall(in *pb.CloudRequest, cloud string) (*pb.CloudResponse, bool) {
// 	cc := dao.RpcConn(RpcAddr["Li"])
// 	defer cc.Close()
// 	c := pb.NewCloud_APIClient(cc)
// 	switch cloud {
// 	case "qcloud":
// 		res, _ := c.Qcloud(context.Background(), in)
// 		return res, true
// 	}
// 	return nil, false
// }

// func Cloud_rawCall(rw http.ResponseWriter, req *http.Request) {
// 	cloud := GetRouteName(req, "cloud")
// 	bsns := GetRouteName(req, "bsns")
// 	if !AuthBsns(cloud, bsns) {
// 		rsp, _ := json.Marshal(&ResponseData{Code: 2, Data: "Cloud not support yet ,damn"})
// 		httprsp(rw, rsp)
// 		return
// 	}

// 	data, ok := AuthPostData(rw, req)
// 	if !ok {
// 		return
// 	}
// 	v := &Cloud_Req{}
// 	json.Unmarshal(data, v)

// 	/*
// 		Control need AAA server to get Cloud id & key
// 	*/
// 	if needAAA {
// 		cloud_id, cloud_key = AAA_GetThirdKey(cloud, v.Auth.ID, v.Auth.Signature, v.Name)
// 		if cloud_id == "" {
// 			rsp, _ := json.Marshal(&ResponseData{2, "AAA failed"})
// 			httprsp(rw, rsp)
// 			return
// 		}
// 	} else {
// 		cloud_id = v.Auth.ID
// 		cloud_key = v.Auth.Signature
// 	}

// 	in := &pb.CloudRequest{Bsns: bsns, Action: v.Action, Region: v.Region, CloudID: cloud_id, CloudKey: cloud_key}
// 	in.Params = v.Param
// 	res, ok := Cloud_CommonCall(in, "qcloud")
// 	if res.Code == 0 {
// 		// rsp, _ := json.Marshal(&Cloud_Res{Code: 0, Data: string(res.Data)})
// 		rsp := res.Data
// 		httprsp(rw, rsp)
// 		return
// 	}
// 	rsp, _ := json.Marshal(&Cloud_Res{Code: 2, Msg: res.Msg})
// 	httprsp(rw, rsp)
// }
