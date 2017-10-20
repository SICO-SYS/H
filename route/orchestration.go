/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package route

import (
	"github.com/SiCo-Ops/H/controller"
)

func Orchestration() {
	v1 := HTTPHandler.PathPrefix("/v1/orchestration").Subrouter()
	v1.Path("/").HandlerFunc(controller.OrchestrationCreate).Methods("POST")
}
