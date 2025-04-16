package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type SSEHandler struct {
	manager *Manager
}

func NewSSEHandler() *SSEHandler {
	return &SSEHandler{
		manager: NewManager(),
	}
}

func (h *SSEHandler) HandleSSE(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	sess := h.manager.CreateSession()
	defer h.manager.RemoveSession(sess.ID)

	// Send session ID to client
	fmt.Fprintf(w, "data: {\"session_id\": \"%s\"}\n\n", sess.ID)
	flusher.Flush()

	t := time.NewTicker(5 * time.Second)
	defer t.Stop()

	for {
		select {
		case msg, ok := <-sess.MessageCh:
			if !ok {
				fmt.Fprintf(w, "close: %s\n\n", msg.Data)
				return
			}
			fmt.Fprintf(w, "data: %s\n\n", msg.Data)
			flusher.Flush()
		case <-r.Context().Done():
			return
		case <-t.C:
			fmt.Println(fmt.Fprintf(w, "ping\n\n")) //间隔5秒发送ping消息，查看氮气的是否存活
			flusher.Flush()
		}
	}
}

func (h *SSEHandler) HandleMessages(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		http.Error(w, "Missing session_id", http.StatusBadRequest)
		return
	}

	var msg Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Invalid message format", http.StatusBadRequest)
		return
	}

	if !h.manager.SendMessage(sessionID, msg) {
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
