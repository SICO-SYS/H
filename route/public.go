/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package route

import (
	"github.com/SiCo-Ops/H/controller"
)

func Public() {
	r := HTTPHandler.PathPrefix("/public").Subrouter()
	r.Path("/version").HandlerFunc(controller.GetCfgVersion).Methods("GET")
	r.Path("/token").HandlerFunc(controller.GetPublicToken).Methods("GET")
}
