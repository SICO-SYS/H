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
	r.Path("/version").HandlerFunc(controller.PublicCfgVersion).Methods("GET")
	r.Path("/token").HandlerFunc(controller.PublicGenerateToken).Methods("GET")
}
