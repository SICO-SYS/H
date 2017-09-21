/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package route

import (
	"github.com/SiCo-Ops/H/controller"
)

func Config() {
	v1 := HTTPHandler.PathPrefix("/v1/config").Subrouter()
	v1.Path("/{environment}").HandlerFunc(controller.ConfigPull).Methods("GET")
	v1.Path("/{environment}").HandlerFunc(controller.ConfigPush).Methods("PUT")
}
