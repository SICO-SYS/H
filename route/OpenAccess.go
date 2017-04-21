/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package route

import (
	"github.com/SiCo-DevOps/H/controller"
)

func OpenAPI() {
	r := Handler.PathPrefix("/open").Subrouter()
	r.NewRoute().Path("/config").HandlerFunc(controller.GetCfgVersion).Methods("GET")
	r.NewRoute().Path("/Token").HandlerFunc(controller.GetOpenToken).Methods("GET")
	r.NewRoute().Path("/APIToken").HandlerFunc(controller.GetAPIToken).Methods("GET")
}
