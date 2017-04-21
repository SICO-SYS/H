/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package route

import (
	"github.com/gorilla/mux"
)

var (
	Handler *mux.Router
)

func init() {
	Handler = mux.NewRouter()
	Handler.PathPrefix("/happy")
	Handler.StrictSlash(true)
	OpenAPI()
	AAA()
}
