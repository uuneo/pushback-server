package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有跨域
	},
}

type SignalMessage struct {
	To   string      `json:"to"`
	Data interface{} `json:"data"`
}

type Client struct {
	ID   string
	Conn *websocket.Conn
}

var (
	clients   = make(map[string]*Client)
	clientsMu sync.RWMutex
)

// HomeController 处理首页请求
// 支持两种功能:
// 1. 通过id参数移除未推送数据
// 2. 生成二维码图片
func HomeController(c *gin.Context) {
	id := c.Query("id")
	if id != "" {
		RemoveNotPushedData(id)
		c.Status(http.StatusOK)
		return
	}

	call := c.Query("call")
	if call != "" {
		WebsocketHandler(c, call)
		return
	}

	url := "https://" + c.Request.Host

	code := c.Query("code")

	if code != "" {
		url = code
	}
	png, err := qrcode.Encode(url, qrcode.High, 1024)

	if err != nil {
		c.JSON(http.StatusOK, failed(http.StatusInternalServerError, "failed to generate QR code: %v", err))
		return
	}

	c.Data(http.StatusOK, "image/png", png)
}

func WebsocketHandler(c *gin.Context, user string) {

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer func() { _ = conn.Close() }()

	client := &Client{ID: user, Conn: conn}
	clientsMu.Lock()
	clients[user] = client
	clientsMu.Unlock()
	defer func() {
		clientsMu.Lock()
		delete(clients, user)
		clientsMu.Unlock()
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		var signal SignalMessage
		if err := json.Unmarshal(msg, &signal); err != nil {
			log.Println("Invalid message format:", err)
			continue
		}
		if signal.To == "" {
			log.Println("Missing 'to' field in message")
			continue
		}
		clientsMu.RLock()
		target, ok := clients[signal.To]
		clientsMu.RUnlock()
		if ok {
			if err := target.Conn.WriteJSON(SignalMessage{To: user, Data: signal.Data}); err != nil {
				log.Println("Write error:", err)
			}
		} else {
			log.Println("Target client not found:", signal.To)
		}
	}
}
