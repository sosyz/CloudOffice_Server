package websocket

import (
	"encoding/json"
	"log"
)

func Login(message *string, tag string) {
	var login struct {
		Name string `json:"tag"`
		Type int    `json:"type"`
	}
	err := json.Unmarshal([]byte(*message), &login)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("%v\n%v", *message, login)
}

func Quiet(message *string, tag string) {

}

func FileInfo(message *string, tag string) {
	var fileInfo struct {
		Name  string `json:"name"`
		Fid   string `json:"fid"`
		Pages int    `json:"pages"`
		Words int    `json:"words"`
	}
	err := json.Unmarshal([]byte(*message), &fileInfo)
	if err != nil {
		log.Println(err)
		return
	}
}

func PrintMessage(message *string, tag string) {
	var printMessage struct {
		Fid     string `json:"fid"`
		OrderId string `json:"orderId"`
		Status  int    `json:"status"`
	}
	err := json.Unmarshal([]byte(*message), &printMessage)
	if err != nil {
		log.Println(err)
		return
	}
}

func Status(message *string, tag string) {
	var status struct {
		// 状态
		Status int `json:"status"`

		// 等待打印文件ID
		Wait []string `json:"wait"`

		// 打印中文件ID
		Printing string `json:"printing"`
	}
	err := json.Unmarshal([]byte(*message), &status)
	if err != nil {
		log.Println(err)
		return
	}
}
