package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	clients   = make(map[string]*Client)
	clientsMu sync.RWMutex
	upgrader  = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有跨域
	}}
	broadcastChan = make(chan BroadcastMessage, 1024)
)

type Client struct {
	ID   string
	Conn *websocket.Conn
}

type BroadcastMessage struct {
	SenderID string
	Message  []byte
}

func init() {
	go broadcastWorker()
}

func broadcastWorker() {
	for msg := range broadcastChan {
		clientsMu.RLock()
		for id, cli := range clients {
			if id == msg.SenderID {
				continue
			}
			err := cli.Conn.WriteMessage(websocket.TextMessage, msg.Message)
			if err != nil {
				log.Printf("Broadcast write error to %s: %v", id, err)
			}
		}
		clientsMu.RUnlock()
	}
}

func WebsocketHandler(c *gin.Context, user string) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, map[string][]string{})

	defer func() {
		clientsMu.Lock()
		_ = conn.Close()
		delete(clients, user)
		clientsMu.Unlock()
	}()

	conn.SetReadLimit(8192)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	client := &Client{ID: user, Conn: conn}
	clientsMu.Lock()
	clients[user] = client
	clientsMu.Unlock()

	for {
		var msg []byte

		_, msg, err = conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		// 只尝试解析出 `to` 字段
		var temp struct {
			To string `json:"to"`
		}
		err = json.Unmarshal(msg, &temp)
		if err != nil {
			log.Println("Invalid JSON format:", err)
			continue
		}

		if temp.To != "" {
			// 定向转发
			clientsMu.RLock()
			target, ok := clients[temp.To]
			clientsMu.RUnlock()
			if ok {
				err = target.Conn.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					log.Println("Write error:", err)
				}
			} else {
				log.Println("Target client not found:", temp.To)
			}
		} else {
			broadcastChan <- BroadcastMessage{
				SenderID: user,
				Message:  msg,
			}
		}
	}
}
