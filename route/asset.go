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
	v1.Path("/template").HandlerFunc(controller.Asset_addTemplate).Methods("POST")
	v1.Path("/sync/{cloud}/{bsns}").HandlerFunc(controller.Asset_synchronize).Methods("POST")
	v1.Path("/instance/{cloud}/{bsns}").HandlerFunc(controller.Asset_addTemplate)
}
