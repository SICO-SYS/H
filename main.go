/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package main

import (
	"H/route"
	"net/http"
)

func Run() {
	http.ListenAndServe("0.0.0.0:80", route.Handler)
}

func main() {
	Run()
}
