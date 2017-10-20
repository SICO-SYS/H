/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package route

import (
	"github.com/SiCo-Ops/H/controller"
)

func Hook() {
	v1 := HTTPHandler.PathPrefix("/v1/hook").Subrouter()
	v1.Path("/").HandlerFunc(controller.HookCreate).Methods("POST")
	v1.Path("/").HandlerFunc(controller.HookCreate).Methods("PUT")
	v1.Path("/{hookName}").HandlerFunc(controller.HookRecive).Methods("POST")
}
