/*

LICENSE:  MIT
Author:   sine
Email:    sinerwr@gmail.com

*/

package cfg

type MysqlInfo struct {
	Host string `json:"host"`
	Port string `json:"port"`
	User string `json:"user"`
	Pass string `json:"pass"`
	Name string `json:"name"`
}

type Mysql struct {
	Default MysqlInfo `json:"default"`
	User    MysqlInfo `json:"user"`
}
