/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

func GetCfgVersion(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("[Success] config version  === " + config.Version))
}

func GetRouteName(req *http.Request, name string) string {
	return mux.Vars(req)[name]
}

func ValidatePostData(rw http.ResponseWriter, req *http.Request) ([]byte, bool) {
	header := req.Header.Get("Content-Type")
	if header != "application/json" {
		rsp, _ := json.Marshal(&ResponseData{2, "request must follow application/json"})
		httprsp(rw, rsp)
		return nil, false
	}
	body, _ := ioutil.ReadAll(req.Body)
	req.Body.Close()
	return body, true
}

func httprsp(rw http.ResponseWriter, rsp []byte) {
	rw.Header().Add("Content-Type", "application/json")
	rw.Write(rsp)
}

type ResponseData struct {
	Code int8        `json:"code"`
	Data interface{} `json:"data"`
}

func ErrorMessage(c int8) string {
	msg := ""
	switch c {
	case 0:
		msg = "[!!!] Do not hack the system"
	// 1 - 10 Receive an incorrect request
	case 1:
		msg = "[Failed] AAA Failed"
	case 2:
		msg = "[Failed] Request Params Incorrect"
	case 3:
		msg = "[Failed] Request Timeout"
	case 4:
		msg = "[Failed] Request Forbidden"
	// 100 - 120 System Error
	// 120 - 127 Middleware Error
	case 125:
		msg = "[Error] MQ crash"
	case 126:
		msg = "[Error] DB crash"
	case 127:
		msg = "[Error] RPC crash"
	default:
		msg = "[Error] Unknown problem"
	}
	return msg
}

func ResponseErrmsg(c int8) *ResponseData {
	return &ResponseData{c, ErrorMessage(c)}
}
