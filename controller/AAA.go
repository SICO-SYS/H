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

type AuthenticationToken struct {
	Id        string `json:"id"`
	Signature string `json:"signature"`
}

type ThirdToken struct {
	Auth    AuthenticationToken `json:"auth"`
	APIType string              `json:"apitype"`
	Name    string              `json:"name"`
	ID      string              `json:"id"`
	Key     string              `json:"key"`
}

func AAAValidateUser(id string, signature string) bool {
	defer func() {
		recover()
		if rcv := recover(); rcv != nil {
			raven.CaptureMessage("Maybe gRPC connect error", nil)
		}
	}()
	cc := rpc.RPCConn(RPCAddr["He"])
	defer cc.Close()
	c := pb.NewAAAPrivateServiceClient(cc)
	in := &pb.AAATokenCall{}
	in.Id = id
	in.Signature = signature
	r, err := c.AuthenticationRPC(context.Background(), in)
	if err != nil {
		return false
	}
	return r.Valid
}

// func CloudTokenRegistry(rw http.ResponseWriter, req *http.Request) {
// 	defer func() {
// 		recover()
// 		if rcv := recover(); rcv != nil {
// 			raven.CaptureMessage("controller.AAARegisThirdpartyKeypair", nil)
// 		}
// 	}()
// 	data, ok := ValidatePostData(rw, req)
// 	v := &ThirdKeypair{}
// 	if ok {
// 		json.Unmarshal(data, v)
// 	} else {
// 		return
// 	}
// 	if v.Name == "" || v.APItype == "" || v.ID == "" {
// 		rsp, _ := json.Marshal(&ResponseData{2, "Missing params, pls follow the guide"})
// 		httprsp(rw, rsp)
// 		return
// 	}
// 	cc := rpc.RPCConn(RPCAddr["He"])
// 	defer cc.Close()
// 	c := pb.NewAAAPrivateServiceClient(cc)
// 	in := &pb.
// 	in.Apitoken = &pb.AAA_APIToken{}
// 	in.Apitoken.Id = v.Auth.Id
// 	in.Apitoken.Signature = v.Auth.Signature
// 	in.Apitype = v.APItype
// 	in.Name = v.Name
// 	in.Id = v.ID
// 	in.Key = v.Key
// 	r, err := c.AAA_ThirdKeypair(context.Background(), in)
// 	if err != nil {
// 		LogErrMsg(50, "controller.PostThirdKeypair")
// 	}

// 	if r.Code == 0 {
// 		rsp, _ := json.Marshal(&ResponseData{0, "Success"})
// 		httprsp(rw, rsp)
// 		return
// 	}
// 	rsp, _ := json.Marshal(&ResponseData{2, r.Msg})
// 	httprsp(rw, rsp)
// }

func AAAAuthentication(rw http.ResponseWriter, req *http.Request) {
	data, ok := ValidatePostData(rw, req)
	v := &AuthenticationToken{}
	if ok {
		json.Unmarshal(data, v)
	} else {
		return
	}
	if AAAValidateUser(v.Id, v.Signature) {
		rsp, _ := json.Marshal(&ResponseData{0, "Success"})
		httprsp(rw, rsp)
		return
	}
	rsp, _ := json.Marshal(&ResponseData{2, "AAA_AuthFailed"})
	httprsp(rw, rsp)
}

// func AAA_GetThirdKey(cloud string, id string, signature string, alias string) (string, string) {
// 	in := &pb.AAA_ThirdpartyKey{}
// 	in.Apitoken = &pb.AAA_APIToken{}
// 	in.Apitoken.Id = id
// 	in.Apitoken.Signature = signature
// 	in.Apitype = cloud
// 	in.Name = alias
// 	cc := dao.RpcConn(RPCAddr["He"])
// 	defer cc.Close()
// 	c := pb.NewAAA_SecretClient(cc)
// 	res, _ := c.AAA_GetThirdKey(context.Background(), in)
// 	if res != nil {
// 		return res.Id, res.Key
// 	}
// 	return "", ""
// }
