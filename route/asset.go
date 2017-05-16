/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package route

import (
	"github.com/SiCo-DevOps/H/controller"
)

func Asset() {
	v1 := Handler.PathPrefix("/v1/asset").Subrouter()
	v1.Path("/template").HandlerFunc(controller.Asset_addTemplate)
}
