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

type AssetTemplate struct {
	PrivateToken AuthenticationToken `json:"token"`
	Name         string              `json:"name"`
	Param        map[string]string   `json:"param"`
}

type AssetSynchronizeRequest struct {
	PrivateToken   AuthenticationToken `json:"token"`
	CloudTokenName string              `json:"name"`
	Region         string              `json:"region"`
}

func AssetCreateTemplate(rw http.ResponseWriter, req *http.Request) {
	data, ok := ValidatePostData(rw, req)
	if !ok {
		return
	}
	v := AssetTemplate{}
	json.Unmarshal(data, &v)

	isPrivateTokenValid, errcode := AAAValidateToken(v.PrivateToken.ID, v.PrivateToken.Signature)
	if errcode != 0 {
		httpResponse("json", rw, responseErrMsg(errcode))
		return
	}
	if config.AAAstatus == "active" && !isPrivateTokenValid {
		httpResponse("json", rw, responseErrMsg(1000))
		return
	}

	in := &pb.AssetTemplateCall{}
	in.Id = v.PrivateToken.ID
	in.Name = v.Name
	in.Params, _ = json.Marshal(v.Param)
	cc, err := rpc.Conn(config.RpcBeHost, config.RpcBePort)
	if err != nil {
		raven.CaptureError(err, nil)
		httpResponse("json", rw, responseErrMsg(303))
		return
	}
	r := rpc.TemplateCreateRPC(cc, in)
	if r.Code != 0 {
		httpResponse("json", rw, responseErrMsg(r.Code))
		return
	}
	httpResponse("json", rw, responseSuccess())
	return
}

func AssetSynchronize(rw http.ResponseWriter, req *http.Request) {
	data, ok := ValidatePostData(rw, req)
	if !ok {
		return
	}

	v := &AssetSynchronizeRequest{}
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

	cloud := getRouteName(req, "cloud")
	service := getRouteName(req, "service")
	region := v.Region
	action, errcode := getActionMap(cloud, service, "DescribeInstances")
	if errcode != 0 {
		httpResponse("json", rw, responseErrMsg(errcode))
		return
	}

	cloudTokenID, cloudTokenKey, errcode := cloudTokenGet(v.PrivateToken.ID, cloud, v.CloudTokenName)
	if errcode != 0 {
		httpResponse("json", rw, responseErrMsg(errcode))
		return
	}

	in := &pb.CloudAPICall{Cloud: cloud, Service: service, Region: region, Action: action, CloudId: cloudTokenID, CloudKey: cloudTokenKey}
	var nextToken string = ""
	var totalCount int64 = 0
	var page int64
	for page = 0; true; page++ {
		in, isLoop := CloudAPICallForLoop(in, nextToken, page, totalCount)
		if !isLoop {
			break
		}
		cloudcc, clouderr := rpc.Conn(config.RpcLiHost, config.RpcLiPort)
		if clouderr != nil {
			raven.CaptureError(clouderr, nil)
			httpResponse("json", rw, responseErrMsg(302))
			return
		}
		cloudResponse := rpc.CloudAPIRequestRPC(cloudcc, in)
		if cloudResponse.Code != 0 {
			httpResponse("json", rw, responseErrMsg(cloudResponse.Code))
			return
		}
		assetcc, asseterr := rpc.Conn(config.RpcBeHost, config.RpcBePort)
		if asseterr != nil {
			raven.CaptureError(asseterr, nil)
			httpResponse("json", rw, responseErrMsg(303))
			return
		}
		assetResponse := rpc.AssetSynchronizeRPC(assetcc, &pb.AssetSynchronizeCall{Id: v.PrivateToken.ID, Cloud: cloud, Service: service, Data: cloudResponse.Data})
		if assetResponse.Code != 0 {
			httpResponse("json", rw, responseErrMsg(assetResponse.Code))
			return
		}
		nextToken = assetResponse.NextToken
		totalCount = assetResponse.TotalCount
	}
	httpResponse("json", rw, responseSuccess())
	return
}
