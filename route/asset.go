/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package route

import (
	"github.com/SiCo-Ops/H/controller"
)

func Asset() {
	v1 := HTTPHandler.PathPrefix("/v1/asset").Subrouter()
	v1.Path("/{cloud}/{bsns}").HandlerFunc(controller.GetCfgVersion).Methods("POST")
}

func Template() {
	v1 := HTTPHandler.PathPrefix("/v1/template").Subrouter()
	v1.Path("/").HandlerFunc(controller.GetCfgVersion).Methods("POST")
}
