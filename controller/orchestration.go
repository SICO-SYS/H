/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"encoding/json"
	// "golang.org/x/crypto/ssh"
	"net/http"

	"github.com/SiCo-Ops/Pb"
	"github.com/SiCo-Ops/dao/grpc"
)

type OrchestrationCreateRequest struct {
	PrivateToken AuthenticationToken `json:"token"`
	HookName     string              `json:"hookname"`
	Project      string              `json:"project"`
	Key          string              `json:"key"`
	Value        string              `json:"value"`
	Task         []string            `json:"task"`
}

func OrchestrationCreate(rw http.ResponseWriter, req *http.Request) {
	data, ok := ValidatePostData(rw, req)
	if !ok {
		return
	}
	v := &OrchestrationCreateRequest{}
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

	carbonClient, carbonErr := rpc.Conn(config.RpcCHost, config.RpcCPort)
	if carbonErr != nil {
		httpResponse("json", rw, responseErrMsg(305))
		return
	}
	in := &pb.HookAuthCall{Id: v.PrivateToken.ID, Hookname: v.HookName}
	carbonRes := rpc.HookAuthRPC(carbonClient, in)
	if carbonRes.Code != 0 {
		httpResponse("json", rw, responseErrMsg(carbonRes.Code))
		return
	}
	if carbonRes.Hookid == "" {
		httpResponse("json", rw, responseErrMsg(5001))
		return
	}
	nitrogenClient, nitrogenErr := rpc.Conn(config.RpcNHost, config.RpcNPort)
	if nitrogenErr != nil {
		httpResponse("json", rw, responseErrMsg(306))
		return
	}
	orchestrationIn := &pb.OrchestrationCreateCall{Hookid: carbonRes.Hookid, Project: v.Project, Key: v.Key, Value: v.Value, Belong: v.PrivateToken.ID, Task: v.Task}
	orchestrationRes := rpc.OrchestrationCreateRPC(nitrogenClient, orchestrationIn)
	if orchestrationRes.Code != 0 {
		httpResponse("json", rw, responseErrMsg(orchestrationRes.Code))
		return
	}
	httpResponse("json", rw, responseSuccess())
	return
}
