package websocket

import (
	"encoding/json"
	"github.com/google/uuid"
	"log"
)

func Login(message *string, client *Client) {
	// {"type":"login","data":"[{\"name\":\"HP-01\", \"type\":1},{\"name\":\"HP-02\", \"type\":1}]","message":""}
	var login []PrinterInfo
	err := json.Unmarshal([]byte(*message), &login)
	if err != nil {
		log.Println(err)
		return
	}

	for _, v := range login {
		tag, _ := uuid.NewUUID()
		// 打印机列表
		PrinterList[tag.String()] = v
		// 客户端所拥有打印机
		hub.tags[client] = append(hub.tags[client], tag.String())
		// 打印机所在客户端
		hub.printer[tag.String()] = client
	}
}

// Quit 客户端下线 清除客户端下的打印机 该处只进行了正常下线处理
func Quit(_ *string, client *Client) {
	// 删除客户端下的打印机
	if hub.tags[client] != nil {
		for _, v := range hub.tags[client] {
			delete(hub.printer, v)
		}
	}

	client.send <- []byte("{\"type\":\"quit\", \"message\":\"ok\"}")
	// 删除客户端
	delete(hub.clients, client)
}

func FileInfo(message *string, client *Client) {
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

func PrintMessage(message *string, client *Client) {
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

func Status(message *string, client *Client) {
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
