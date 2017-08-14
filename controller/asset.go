/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"encoding/json"
	// "fmt"
	// "github.com/getsentry/raven-go"
	"golang.org/x/net/context"
	"net/http"

	"github.com/SiCo-Ops/Pb"
	"github.com/SiCo-Ops/dao/grpc"
	// "github.com/SiCo-Ops/dao/mongo"
	// "github.com/SiCo-Ops/public"
)

type AssetTemplate struct {
	PrivateToken AuthenticationToken `json:"token"`
	Name         string              `json:"name"`
	Param        map[string]string   `json:"param"`
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
		rsp, _ := json.Marshal(ResponseErrmsg(2))
		httprsp(rw, rsp)
		return
	}
	in.Id = v.PrivateToken.ID
	in.Name = v.Name
	in.Param = v.Param
	cc := rpc.RPCConn(RPCAddr["Be"])
	defer cc.Close()
	c := pb.NewAssetClient(cc)
	res, _ := c.CreateTemplateRPC(context.Background(), in)
	if res.Code == 0 {
		rsp, _ := json.Marshal(&ResponseData{0, "Success add template"})
		httprsp(rw, rsp)
		return
	}
	rsp, _ := json.Marshal(res)
	httprsp(rw, rsp)
	return
}

// func Asset_synchronize(rw http.ResponseWriter, req *http.Request) {
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
// 		cloud_id, cloud_key = AAA_GetThirdKey(cloud, v.Auth.Id, v.Auth.Signature, v.Name)
// 		if cloud_id == "" {
// 			rsp, _ := json.Marshal(&ResponseData{2, "AAA failed"})
// 			httprsp(rw, rsp)
// 			return
// 		}
// 	} else {
// 		cloud_id = v.Auth.Id
// 		cloud_key = v.Auth.Signature
// 	}

// 	in := &pb.CloudRequest{Bsns: bsns, Action: "list_ins", Region: v.Region, CloudId: cloud_id, CloudKey: cloud_key}
// 	in.Params["Limit"] = "1"
// 	res, ok := Cloud_CommonCall(in, "qcloud")
// 	if res.Code == 0 {
// 		v := make(map[string]interface{})
// 		json.Unmarshal(res.Data, &v)
// 		var count int
// 		if bsns == "cvm" {
// 			totalcount, ok := v["Response"].(map[string]interface{})["TotalCount"].(float64)
// 			if ok {
// 				count = int(totalcount)
// 			}
// 		} else {
// 			totalcount, ok := v["totalCount"].(string)
// 			if ok {
// 				count = public.Atoi(totalcount)
// 			}
// 		}

// 		var looptime int
// 		if count%100 == 0 {
// 			looptime = count / 100
// 		} else {
// 			looptime = count/100 + 1
// 		}
// 		in.Params["Limit"] = "100"
// 		for i := 0; i < looptime; i++ {
// 			in.Params["Offset"] = public.Itoa(i * 100)
// 			fmt.Println(in)
// 		}
// 	}
// 	rsp, _ := json.Marshal(&Cloud_Res{Code: 2, Msg: res.Msg})
// 	httprsp(rw, rsp)

// 	// collectionName := "asset."
// }
