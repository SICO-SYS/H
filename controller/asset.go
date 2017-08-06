/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"net/http"

	"github.com/SiCo-Ops/Pb"
	"github.com/SiCo-Ops/dao"
	"github.com/SiCo-Ops/public"
)

type AssetTemplate struct {
	Auth  AuthToken         `json:"auth"`
	Name  string            `json:"name"`
	Param map[string]string `json:"param"`
}

func Asset_addTemplate(rw http.ResponseWriter, req *http.Request) {
	v := AssetTemplate{}
	if data, ok := AuthPostData(rw, req); ok {
		json.Unmarshal(data, &v)
	} else {
		return
	}
	in := &pb.Asset_Req{}
	if needAAA {
		if AAA(v.Auth.Id, v.Auth.Signature) {
			in.Id = v.Auth.Id
			in.Name = v.Name
			in.Param = v.Param
			cc := dao.RpcConn(RpcAddr["Be"])
			defer cc.Close()
			c := pb.NewAseetClient(cc)
			res, _ := c.AssetTemplate(context.Background(), in)
			if res.Code == 2 {
				LogErrMsg(51, "")
				rsp, _ := json.Marshal(&ResponseData{2, "request error"})
				httprsp(rw, rsp)
				return
			}
			if res.Code == 1 {
				rsp, _ := json.Marshal(&ResponseData{1, "Cannot use same template name"})
				httprsp(rw, rsp)
				return
			}
			rsp, _ := json.Marshal(&ResponseData{0, "Success add template"})
			httprsp(rw, rsp)
			return
		} else {
			rsp, _ := json.Marshal(&ResponseData{2, "AAA failed"})
			httprsp(rw, rsp)
			return
		}
	} else {
		in.Id = v.Auth.Id
		in.Name = v.Name
		in.Param = v.Param
		cc := dao.RpcConn(RpcAddr["Be"])
		defer cc.Close()
		c := pb.NewAseetClient(cc)
		res, _ := c.AssetTemplate(context.Background(), in)
		if res.Code != 0 {
			LogErrMsg(51, "")
			rsp, _ := json.Marshal(&ResponseData{2, "request error"})
			httprsp(rw, rsp)
			return
		}
	}
}

func Asset_synchronize(rw http.ResponseWriter, req *http.Request) {
	cloud := GetRouteName(req, "cloud")
	bsns := GetRouteName(req, "bsns")
	if !AuthBsns(cloud, bsns) {
		rsp, _ := json.Marshal(&ResponseData{Code: 2, Data: "Cloud not support yet ,damn"})
		httprsp(rw, rsp)
		return
	}

	data, ok := AuthPostData(rw, req)
	if !ok {
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

	in := &pb.CloudRequest{Bsns: bsns, Action: "list_ins", Region: v.Region, CloudId: cloud_id, CloudKey: cloud_key}
	in.Params["Limit"] = "1"
	res, ok := Cloud_CommonCall(in, "qcloud")
	if res.Code == 0 {
		v := make(map[string]interface{})
		json.Unmarshal(res.Data, &v)
		var count int
		if bsns == "cvm" {
			totalcount, ok := v["Response"].(map[string]interface{})["TotalCount"].(float64)
			if ok {
				count = int(totalcount)
			}
		} else {
			totalcount, ok := v["totalCount"].(string)
			if ok {
				count = public.Atoi(totalcount)
			}
		}

		var looptime int
		if count%100 == 0 {
			looptime = count / 100
		} else {
			looptime = count/100 + 1
		}
		in.Params["Limit"] = "100"
		for i := 0; i < looptime; i++ {
			in.Params["Offset"] = public.Itoa(i * 100)
			fmt.Println(in)
		}
	}
	rsp, _ := json.Marshal(&Cloud_Res{Code: 2, Msg: res.Msg})
	httprsp(rw, rsp)

	// collectionName := "asset."
}
