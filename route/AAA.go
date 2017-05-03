/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package route

import (
	"github.com/SiCo-DevOps/H/controller"
)

func AAA() {
	v1 := Handler.PathPrefix("/v1/AAA").Subrouter()
	v1.Path("/keypair").HandlerFunc(controller.AAA_PostThirdKeypair).Methods("POST")
	v1.Path("/").HandlerFunc(controller.AAA_Auth).Methods("POST")
	v1.Path("/authorization").HandlerFunc(controller.GetCfgVersion).Methods("POST")
}
