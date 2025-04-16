package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	manager     *Manager
	connections sync.Map
	upgrader    websocket.Upgrader
}

func NewWebSocketHandler() *WebSocketHandler {
	return &WebSocketHandler{
		manager: NewManager(),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // 允许所有来源的连接，生产环境中应该更严格
			},
		},
	}
}

func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Incoming WebSocket connection from %s\n", r.RemoteAddr)

	// 升级HTTP连接为WebSocket连接
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Failed to upgrade connection from %s: %v\n", r.RemoteAddr, err)
		http.Error(w, fmt.Sprintf("Could not upgrade connection: %v", err), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	// 创建新的会话
	sess := h.manager.CreateSession()
	defer h.manager.RemoveSession(sess.ID)

	// 存储WebSocket连接
	h.connections.Store(sess.ID, conn)
	defer h.connections.Delete(sess.ID)

	// 发送会话ID给客户端
	if err := conn.WriteJSON(map[string]interface{}{"type": "session", "session_id": sess.ID}); err != nil {
		return
	}

	// 创建一个用于发送ping消息的ticker
	pingTicker := time.NewTicker(5 * time.Second)
	defer pingTicker.Stop()

	// 创建一个用于通知goroutine退出的channel
	done := make(chan struct{})
	defer close(done)

	// 主动发给客户端
	go func() {
		for {
			select {
			case msg, ok := <-sess.MessageCh:
				if !ok {
					fmt.Printf("Message channel closed for session %s\n", sess.ID)
					return
				}
				if err := conn.WriteJSON(map[string]interface{}{"type": "message", "message": msg.Data}); err != nil {
					fmt.Printf("Error sending message to session %s: %v\n", sess.ID, err)
					return
				}
			case <-pingTicker.C:
				if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
					fmt.Printf("Failed to send ping message to session %s: %v\n", sess.ID, err)
					return
				}
			case <-done:
				return
			}
		}
	}()

	// 接收客户端消息
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("Error reading message from session %s: %v\n", sess.ID, err)
			break
		}
		switch messageType {
		case websocket.TextMessage:
			var jsonMsg struct {
				Message string `json:"message"`
			}
			if err := json.Unmarshal(p, &jsonMsg); err != nil {
				fmt.Printf("Error parsing JSON message from session %s: %v\n", sess.ID, err)
				continue
			}
			if err := conn.WriteJSON(map[string]interface{}{"type": "response", "session_id": sess.ID}); err != nil {
				fmt.Printf("Failed to send ping JSON to session %s: %v\n", sess.ID, err)
				return
			}
		case websocket.PongMessage:
			fmt.Printf("Received pong from session %s\n", sess.ID)
		}
	}
}
