/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package cfg

type Log struct {
	Logpath string `json:"logpath"`
	WARNING string `json:"warning"`
	ERROR   string `json:"error"`
	FATAL   string `json:"fatal"`
}
