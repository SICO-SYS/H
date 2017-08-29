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
	v1.Path("/cloud/{cloud}/{service}").HandlerFunc(controller.PublicCfgVersion).Methods("POST")
	v1.Path("/synchronize/{cloud}/{service}").HandlerFunc(controller.AssetSynchronize).Methods("POST")
	v1.Path("/template").HandlerFunc(controller.AssetCreateTemplate).Methods("POST")
}
