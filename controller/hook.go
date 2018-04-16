/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"encoding/json"
	"github.com/getsentry/raven-go"
	"net/http"

	"github.com/SiCo-Ops/Pb"
	"github.com/SiCo-Ops/dao/grpc"
)

type hookCreate struct {
	PrivateToken AuthenticationToken `json:"token"`
	HookType     string              `json:"hooktype"`
}

type hookUpdate struct {
	PrivateToken AuthenticationToken
	HookName     string
}

func HookCreate(rw http.ResponseWriter, req *http.Request) {
	data, ok := ValidatePostData(rw, req)
	if !ok {
		return
	}
	v := &hookCreate{}
	json.Unmarshal(data, v)
	isPrivateTokenValid, errcode := AAAValidateToken(v.PrivateToken.ID, v.PrivateToken.Signature)
	if errcode != 0 {
		httpResponse("json", rw, responseErrMsg(errcode))
		return
	}
	if config.AAAstatus == "active" && !isPrivateTokenValid {
		httpResponse("json", rw, responseErrMsg(1000))
		return
	}

	in := &pb.HookCreateCall{Id: v.PrivateToken.ID, Hooktype: v.HookType}
	cc, err := rpc.Conn(config.RpcCHost, config.RpcCPort)
	if err != nil {
		raven.CaptureError(err, nil)
		httpResponse("json", rw, responseErrMsg(305))
		return
	}
	r := rpc.HookCreateRPC(cc, in)
	if r.Code != 0 {
		httpResponse("json", rw, responseErrMsg(r.Code))
		return
	}
	rsp, _ := json.Marshal(&responseData{Code: 0, Data: r.Hookname})
	httpResponse("json", rw, rsp)
	return
}

func HookRecive(rw http.ResponseWriter, req *http.Request) {
	hookName := getRouteName(req, "hookName")
	cc, err := rpc.Conn(config.RpcCHost, config.RpcCPort)
	if err != nil {
		raven.CaptureError(err, nil)
		httpResponse("json", rw, responseErrMsg(305))
		return
	}
	in := &pb.HookQueryCall{Hookname: hookName}
	r := rpc.HookQueryRPC(cc, in)
	if r.Code != 0 {
		httpResponse("json", rw, responseErrMsg(r.Code))
		return
	}
	hookType := r.Hooktype
	var data []byte
	switch hookType {
	case "travis":
		data = []byte(req.PostForm.Get("payload"))
	default:
		var ok bool
		data, ok = ValidatePostData(rw, req)
		if !ok {
			return
		}
	}
	cc, err = rpc.Conn(config.RpcCHost, config.RpcCPort)
	if err != nil {
		raven.CaptureError(err, nil)
		httpResponse("json", rw, responseErrMsg(305))
		return
	}
	hookin := &pb.HookReceiveCall{Hookname: hookName, Hooktype: hookType, Payload: data}
	hookres := rpc.HookReceiveRPC(cc, hookin)
	if hookres.Code != 0 {
		httpResponse("json", rw, responseErrMsg(hookres.Code))
		return
	}
	if hookres.Params == nil {
		httpResponse("json", rw, responseSuccess())
		return
	}

	Nclient, Nerr := rpc.Conn(config.RpcNHost, config.RpcNPort)
	if Nerr != nil {
		raven.CaptureError(Nerr, nil)
		httpResponse("json", rw, responseErrMsg(306))
		return
	}
	Nin := &pb.OrchestrationCheckCall{Hookid: r.Hookid, Id: r.Belong, Type: hookType, Params: hookres.Params}
	Nres := rpc.OrchestrationCheckRPC(Nclient, Nin)
	if Nres.Code != 0 {
		httpResponse("json", rw, responseErrMsg(Nres.Code))
		return
	}
	if Nres.Task != nil {
		rsp, _ := json.Marshal(&responseData{Code: 0, Data: Nres.Task})
		httpResponse("json", rw, rsp)
		return
	}
	httpResponse("json", rw, responseSuccess())
	return
}
