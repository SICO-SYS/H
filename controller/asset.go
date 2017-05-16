/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"encoding/json"
	"golang.org/x/net/context"
	"net/http"

	"github.com/SiCo-DevOps/Pb"
	"github.com/SiCo-DevOps/dao"
	. "github.com/SiCo-DevOps/log"
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
			if res.Code != 0 {
				LogErrMsg(51, "")
				rsp, _ := json.Marshal(&ResponseData{2, "request error"})
				httprsp(rw, rsp)
				return
			}
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
