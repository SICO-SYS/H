/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package route

import (
	"github.com/SiCo-Ops/H/controller"
)

func Cloud() {
	v1 := HTTPHandler.PathPrefix("/v1/cloud").Subrouter()
	v1.HandleFunc("/{cloud}/{service}", controller.CloudAPICall).Methods("POST")
	v1.Path("/token").HandlerFunc(controller.CloudTokenRegistry).Methods("POST")
}
