package websocket

import (
	"github.com/gorilla/websocket"
	"sonui.cn/cloudprint/utils"
)

var hub = &Hub{
	broadcast:  make(chan []byte),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	clients:    make(map[*Client]bool),
	printer:    make(map[string]*Client),
	blackQueue: utils.Queue[*string]{},
	colorQueue: utils.Queue[*string]{},
	tags:       make(map[*Client][]string),
}

type Client struct {
	hub *Hub

	conn *websocket.Conn

	send chan []byte
}

type PrinterInfo struct {
	Name string `json:"name"`
	Type int    `json:"type"`
	Tag  string `json:"tag"`
}

// PrinterList 打印机ID -> 打印机参数
var PrinterList = make(map[string]PrinterInfo)

type Hub struct {
	// 客户端
	clients map[*Client]bool

	// 客户端名下的打印机ID
	tags map[*Client][]string

	// 打印机所在客户端
	printer map[string]*Client

	// 黑白打印机队列
	blackQueue utils.Queue[*string]

	// 彩色打印机队列
	colorQueue utils.Queue[*string]

	// 广播器
	broadcast chan []byte

	// 新客户端注册
	register chan *Client

	// 客户端下线
	unregister chan *Client
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				// 删除客户端名下的打印机ID 检测是否为打印机连接
				if h.tags[client] != nil {
					for _, v := range h.tags[client] {
						delete(h.printer, v)
					}
					delete(h.tags, client)
				}
				// 删除客户端
				delete(h.clients, client)
				// 关闭管道
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
