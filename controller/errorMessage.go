/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"encoding/json"
)

func errorMessage(errcode int64) string {
	switch errcode {
	// 1 - 99 Gateway
	case 1:
		errmsg = "[Gateway] Thanks for testing hack"
	case 2:
		errmsg = "[Gateway.ValidatePostData] Request Content-Type make sure application/json"
	case 3:
		errmsg = "[Gateway.actionMap] ActionMap.json file not found"
	case 4:
		errmsg = "[Gateway.actionMap] Not support specified cloud yet"
	case 5:
		errmsg = "[Gateway.actionMap] Not support specified service yet"
	case 6:
		errmsg = "[Gateway.actionMap] Not support specified action yet"
	case 7:
		errmsg = "[Gateway] Public Token Service is shutdown by manually"
	case 8:
		errmsg = "[Gateway] Invalid public token"
	case 9:
		errmsg = "[Gateway.CloudServiceIsSupport] cloud.json file not found"
	// 100 - 199 Redis error
	case 100:
		errmsg = "[Redis] DB problem"
	case 101:
		errmsg = "[Redis] PublicDB has some problem, pls follow sentry to resolve"
	case 102:
		errmsg = "[Redis] ConfigDB has some problem, pls follow sentry to resolve"
	//200-299 mongo error
	case 200:
		errmsg = "[Mongo] DB problem"
	case 201:
		errmsg = "[Mongo] UserDB has some problem, pls follow sentry to resolve"
	case 202:
		errmsg = "[Mongo] CloudDB has some problem, pls follow sentry to resolve"
	case 203:
		errmsg = "[Mongo] AssetDB has some problem, pls follow sentry to resolve"
	case 204:
		errmsg = "[Mongo] ConfigDB has some problem, pls follow sentry to resolve"
	// 300 - 399 gRPC error
	case 300:
		errmsg = "[gRPC] RPC call failed"
	case 301:
		errmsg = "[gRPC] He has some problem, pls follow sentry to resolve"
	case 302:
		errmsg = "[gRPC] Li has some problem, pls follow sentry to resolve"
	case 303:
		errmsg = "[gRPC] Be has some problem, pls follow sentry to resolve"
	case 304:
		errmsg = "[gRPC] B has some problem, pls follow sentry to resolve"
	case 305:
		errmsg = "[gRPC] C has some problem, pls follow sentry to resolve"
	case 306:
		errmsg = "[gRPC] N has some problem, pls follow sentry to resolve"
	case 307:
		errmsg = "[gRPC] O has some problem, pls follow sentry to resolve"
	case 308:
		errmsg = "[gRPC] F has some problem, pls follow sentry to resolve"
	case 309:
		errmsg = "[gRPC] Ne has some problem, pls follow sentry to resolve"
	// 1000 - 1999 AAA error
	case 1000:
		errmsg = "[AAA] Authentication failed"
	case 1001:
		errmsg = "[AAA] Authorization failed"
	case 1002:
		errmsg = "[AAA] Accounting failed"
	case 1003:
		errmsg = "[AAA] Missing email or phone"
	case 1004:
		errmsg = "[AAA] Generate token retry more than 5 times, already report event automatically"
	case 1999:
		errmsg = "[AAA] unknown error"
	//2000 - 2999 Cloud problem
	case 2000:
		errmsg = "[cloud] Missing Name or Cloud or ID"
	case 2001:
		errmsg = "[cloud] Request Third-party API error,pls follow sentry to resolve"
	case 2002:
		errmsg = "[cloud] The cloud you specified not support yet"
	case 2003:
		errmsg = "[cloud] Cloud token not exist"

	//3000-3999 Asset problem
	case 3000:
		errmsg = "[asset] Synchronize not support this service yet"
	case 3001:
		errmsg = "[asset] Template"
	//General error message
	default:
		errmsg = "[Common] platform error"
	}
	return errmsg
}

func responseErrMsg(errcode int64) []byte {
	rsp, _ := json.Marshal(&responseData{errcode, errorMessage(errcode)})
	return rsp
}

func responseSuccess() []byte {
	rsp, _ := json.Marshal(&responseData{0, "Success"})
	return rsp
}
