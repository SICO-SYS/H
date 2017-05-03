/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"encoding/json"
	"golang.org/x/net/context"
	"io/ioutil"
	"net/http"

	"github.com/SiCo-DevOps/Pb"
	"github.com/SiCo-DevOps/dao"
	. "github.com/SiCo-DevOps/log"
)

var (
	cloud_id     string
	cloud_key    string
	cloud_region string
	cloud_action string
)

type Cloud_Req struct {
	Auth   AuthToken         `json:"auth"`
	Region string            `json:"region"`
	Action string            `json:"action"`
	Name   string            `json:"name"`
	Param  map[string]string `json:"params"`
}

type Cloud_Res struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

func AuthBsns(cloud string, bsns string) bool {
	d, err := ioutil.ReadFile("cloud.json")
	if err != nil {
		LogFatalMsg(0, "controller.AuthBsns")
	}
	var v map[string][]string
	json.Unmarshal(d, &v)
	if value, ok := v[cloud]; ok {
		for _, v := range value {
			if v == bsns {
				return true
			}
		}
		return false
	}
	return false
}

func Cloud_CommonCall(in *pb.CloudRequest, cloud string) (*pb.CloudResponse, bool) {
	cc := dao.RpcConn(RpcAddr["Li"])
	defer cc.Close()
	c := pb.NewCloud_APIClient(cc)
	switch cloud {
	case "qcloud":
		res, _ := c.Qcloud(context.Background(), in)
		return res, true
	}
	return nil, false
}

func Cloud_SyncResourse(rw http.ResponseWriter, req *http.Request) {
	cloud := GetRouteName(req, "cloud")
	bsns := GetRouteName(req, "bsns")
	cloud_region := req.URL.Query().Get("region")
	id := req.URL.Query().Get("id")
	signature := req.URL.Query().Get("signature")
	alias := req.URL.Query().Get("name")
	if !AuthBsns(cloud, bsns) {
		rsp, _ := json.Marshal(&ResponseData{Code: 2, Data: "Cloud not support yet ,damn"})
		httprsp(rw, rsp)
		return
	}

	/*
		Control need AAA server to get Cloud id & key
	*/
	if needAAA {
		cloud_id, cloud_key = AAA_GetThirdKey(cloud, id, signature, alias)
	} else {
		cloud_id = req.URL.Query().Get("id")
		cloud_key = req.URL.Query().Get("key")
	}

	cc := dao.RpcConn(RpcAddr["Li"])
	defer cc.Close()
	c := pb.NewCloud_APIClient(cc)
	in := &pb.CloudRequest{Bsns: bsns, Action: "listall", Region: cloud_region, CloudId: cloud_id, CloudKey: cloud_key}
	res, _ := c.Qcloud(context.Background(), in)
	if res.Code == 0 {
		rsp, _ := json.Marshal(&ResponseData{Code: 0, Data: res.Msg})
		httprsp(rw, rsp)
		return
	}
	rsp, _ := json.Marshal(&ResponseData{Code: 2, Data: res.Msg})
	httprsp(rw, rsp)
}

func Cloud_Call(rw http.ResponseWriter, req *http.Request) {
	cloud := GetRouteName(req, "cloud")
	bsns := GetRouteName(req, "bsns")
	if !AuthBsns(cloud, bsns) {
		rsp, _ := json.Marshal(&ResponseData{Code: 2, Data: "Cloud not support yet ,damn"})
		httprsp(rw, rsp)
		return
	}

	data, ok := AuthPostData(req)
	if !ok {
		rsp, _ := json.Marshal(&ResponseData{2, "request must follow application/json"})
		httprsp(rw, rsp)
		return
	}
	v := &Cloud_Req{}
	json.Unmarshal(data, v)

	/*
		Control need AAA server to get Cloud id & key
	*/
	if needAAA {
		cloud_id, cloud_key = AAA_GetThirdKey(cloud, v.Auth.Id, v.Auth.Signature, v.Name)
		if cloud_id == "" {
			rsp, _ := json.Marshal(&ResponseData{2, "AAA failed"})
			httprsp(rw, rsp)
			return
		}
	} else {
		cloud_id = v.Auth.Id
		cloud_key = v.Auth.Signature
	}

	in := &pb.CloudRequest{Bsns: bsns, Action: v.Action, Region: v.Region, CloudId: cloud_id, CloudKey: cloud_key}
	in.Params = []*pb.CloudParams{}
	for param_key, param_value := range v.Param {
		in.Params = append(in.Params, &pb.CloudParams{Key: param_key, Value: param_value})
	}
	res, ok := Cloud_CommonCall(in, "qcloud")
	if res.Code == 0 {
		// rsp, _ := json.Marshal(&Cloud_Res{Code: 0, Data: string(res.Data)})
		rsp := res.Data
		httprsp(rw, rsp)
		return
	}
	rsp, _ := json.Marshal(&Cloud_Res{Code: 2, Msg: res.Msg})
	httprsp(rw, rsp)
}
