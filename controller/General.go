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

var (
	errmsg string
)

func PublicCfgVersion(rw http.ResponseWriter, req *http.Request) {
	rsp, _ := json.Marshal(&responseData{Code: 0, Data: "[Success] " + config.Version})
	httpResponse("json", rw, rsp)
}

func getRouteName(req *http.Request, name string) string {
	return mux.Vars(req)[name]
}

func ValidatePostData(rw http.ResponseWriter, req *http.Request) ([]byte, bool) {
	header := req.Header.Get("Content-Type")
	if header != "application/json" {
		httpResponse("json", rw, responseErrMsg(2))
		return nil, false
	}
	body, _ := ioutil.ReadAll(req.Body)
	req.Body.Close()
	return body, true
}

func getActionMap(cloud string, service string, action string) (string, int64) {
	d, err := ioutil.ReadFile("ActionMap.json")
	if err != nil {
		return "", 3
	}

	actionMap := make(map[string]interface{})
	json.Unmarshal(d, &actionMap)

	cloudMap, ok := actionMap[action].(map[string]interface{})
	if !ok {
		return "", 6
	}
	serviceMap, ok := cloudMap[cloud].(map[string]interface{})
	if !ok {
		return "", 4
	}
	value, ok := serviceMap[service].(string)
	if !ok {
		return "", 5
	}
	return value, 0
}

func httpResponse(contentType string, rw http.ResponseWriter, rsp []byte) {
	switch contentType {
	case "xml":
		rw.Header().Add("Content-Type", "application/xml")
	default:
		rw.Header().Add("Content-Type", "application/json")
	}
	rw.Write(rsp)
}

type responseData struct {
	Code int64       `json:"code"`
	Data interface{} `json:"data"`
}
