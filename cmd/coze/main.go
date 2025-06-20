package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/long250038728/web/cmd/coze/client"
	"net/http"
)

func main() {
	handler := gin.Default()

	handler.GET("/coze/coze/sse", func(context *gin.Context) {
		req := &client.ChatRequest{
			ConversationID: "7517117256539881526",
			BotID:          "7479292866154561548",
			UserID:         "1",
			Content:        "珠宝系统的作用",
		}

		ch, err := (&client.Client{}).StreamChat(context.Request.Context(), req)
		if err != nil {
			_, _ = context.Writer.Write([]byte(err.Error()))
			return
		}

		context.Header("Content-Type", "text/event-stream; charset=utf-8")
		context.Header("Cache-Control", "no-store")
		context.Header("Connection", "keep-alive")

		ginWriter := context.Writer
		f, _ := ginWriter.(http.Flusher)

		for message := range ch {
			_, _ = fmt.Fprintf(ginWriter, message.Content)
			f.Flush()
		}
	})

	// curl -X POST "http://192.168.1.5:8080/coze/coze/sse2"
	// -H "Content-Type: application/json"
	// -d '{"bot_id":"7479292866154561548","conversation_id":"7517117256539881526","user_id":"1","content":"aaa"}'

	handler.POST("/coze/coze/sse2", func(context *gin.Context) {
		req := &client.ChatRequest{}
		if err := context.ShouldBindJSON(&req); err != nil {
			context.JSON(400, gin.H{"error": err.Error()})
			return
		}

		ch, err := (&client.Client{}).StreamChat(context.Request.Context(), req)
		if err != nil {
			_, _ = context.Writer.Write([]byte(err.Error()))
			return
		}

		context.Header("Content-Type", "text/event-stream; charset=utf-8")
		context.Header("Cache-Control", "no-store")
		context.Header("Connection", "keep-alive")
		ginWriter := context.Writer
		f, _ := ginWriter.(http.Flusher)

		for message := range ch {
			_, _ = fmt.Fprintf(ginWriter, message.Content)
			f.Flush()
		}
	})

	server := &http.Server{Addr: fmt.Sprintf("%s:%d", "192.168.1.5", 8080), Handler: handler}
	fmt.Println(server.ListenAndServe())
}
