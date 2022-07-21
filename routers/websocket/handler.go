package websocket

import (
	"encoding/json"
	"log"
	"time"
)

func NewTask(message *string, client *Client) {
	var newTask struct {
		Task    int      `json:"task"` // 1: 获取文件信息（暂不启用） 2: 打印黑白文件 3: 打印彩色文件
		Fid     []string `json:"fid"`
		OrderId string   `json:"orderId"`
	}
	err := json.Unmarshal([]byte(*message), &newTask)
	if err != nil {
		log.Println(err)
		return
	}
	var c *string
	// 优化方案:
	// 添加空闲队列 进行维护 从空闲队列中取出打印机进行打印

	i := 0
	for *c == "" && i < 30 {
		if newTask.Task == 3 {
			c = hub.colorQueue.Pop()
		} else {
			c = hub.blackQueue.Pop()
			if *c == "" {
				c = hub.colorQueue.Pop()
				if *c != "" {
					hub.colorQueue.Push(c) // 放到队尾
				}
			} else {
				hub.blackQueue.Push(c) // 放到队尾
			}
		}
		if *c != "" {
			break
		}
		time.Sleep(time.Second * 1)
		i++
	}

	if *c == "" {
		// TODO: 设置状态为打印失败
	} else {
		// 发送给打印机
		hub.printer[*c].send <- []byte(*message)
	}
}

func Login(message *string, client *Client) {
	// {"type":"login","data":"[{\"name\":\"HP-01\", \"type\":1},{\"name\":\"HP-02\", \"type\":1}]","message":""}
	var login []PrinterInfo
	err := json.Unmarshal([]byte(*message), &login)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("%+v\n", login)
	for k, v := range login {
		tag := login[k].Tag
		// 打印机列表
		PrinterList[tag] = v
		// 客户端所拥有打印机
		hub.tags[client] = append(hub.tags[client], tag)
		// 打印机所在客户端
		hub.printer[tag] = client

		if v.Type == 1 {
			// 黑白打印机
			hub.blackQueue.Push(&tag)
		} else if v.Type == 2 {
			// 彩色打印机
			hub.colorQueue.Push(&tag)
		}
	}
}

// Quit 客户端下线 清除客户端下的打印机 该处只进行了正常下线处理
func Quit(_ *string, client *Client) {
	// 从队列中删除该客户端下的打印机

	// 删除客户端下的打印机
	if hub.tags[client] != nil {
		for _, v := range hub.tags[client] {
			delete(hub.printer, v)
		}
		delete(hub.tags, client)
	}

	client.send <- []byte("{\"type\":\"quit\", \"message\":\"ok\"}")
	// 删除客户端
	delete(hub.clients, client)
}

// FileInfo 文件信息上报 暂不启用
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

// PrintMessage 打印状态上报d
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

// Status 打印机状态上报
func Status(message *string, client *Client) {
	log.Println("Status:", *message)
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
