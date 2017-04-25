/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package controller

import (
	"io/ioutil"
	"net/http"

	"github.com/SiCo-DevOps/cfg"
)

var (
	config  = cfg.Config
	errcode int8
	err     error
	RpcAddr = map[string]string{
		"He": "He.SiCo",
		"Li": "Li.SiCo",
		"Be": "Be.SiCo",
		"B":  "B.SiCo",
		"C":  "C.SiCo",
		"N":  "N.SiCo",
	}
)

func AuthPostData(req *http.Request) ([]byte, bool) {
	header := req.Header.Get("Content-Type")
	if header != "application/json" {
		return nil, false
	}
	body, _ := ioutil.ReadAll(req.Body)
	req.Body.Close()
	return body, true
}
