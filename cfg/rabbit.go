/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package cfg

type RabbitInfo struct {
	Host string `json:"host"`
	Port string `json:"port"`
	User string `json:"user"`
	Pass string `json:"pass"`
}

type Rabbit struct {
	Default RabbitInfo `json:"default"`
}
