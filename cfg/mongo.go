/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package cfg

type MongoInfo struct {
	Host string `json:"host"`
	Port string `json:"port"`
	User string `json:"user"`
	Pass string `json:"pass"`
	Name string `json:"name"`
}

type Mongo struct {
	Default MongoInfo `json:"default"`
	Geo     MongoInfo `json:"geo"`
}
