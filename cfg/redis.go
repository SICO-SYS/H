/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package cfg

type RedisInfo struct {
	Host string `json:"host"`
	Port string `json:"port"`
	Pass string `json:"pass"`
}

type Redis struct {
	Default RedisInfo `json:"default"`
}
