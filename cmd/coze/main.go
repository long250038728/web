package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/long250038728/web/cmd/coze/client"
	"github.com/long250038728/web/cmd/coze/handle"
	"net/http"
)

func main() {
	handler := gin.Default()
	cozeCli, err := client.NewCozeClient()
	if err != nil {
		fmt.Println(err.Error())
	}
	cli := handle.NewHandle(cozeCli)

	handler.Use(func(context *gin.Context) {
		// 设置跨域相关的头部信息
		context.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		context.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		context.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		// 如果是预检请求（OPTIONS 方法），则直接返回成功响应
		if context.Request.Method == "OPTIONS" {
			context.AbortWithStatus(http.StatusNoContent) // 204 No Content
			return
		}
	})

	//================================================================================
	handler.POST("/ai/conversation/list", func(context *gin.Context) {
		req := &handle.ConversationsMessageListRequest{}
		//if err := context.ShouldBindJSON(req); err != nil {
		//	context.JSON(http.StatusOK, gin.H{"data": map[string]any{}, "return_message": err.Error(), "return_code": "000001"})
		//	return
		//}
		resp, err := cli.ConversationsMessageList(context.Request.Context(), req)
		if err != nil {
			context.JSON(http.StatusOK, gin.H{"data": map[string]any{}, "return_message": err.Error(), "return_code": "000001"})
			return
		}
		context.JSON(http.StatusOK, gin.H{"data": resp, "return_message": "操作成功", "return_code": "000000"})
	})

	handler.POST("/ai/chat/stream", func(context *gin.Context) {
		req := &handle.StreamChatRequest{}
		//if err := context.ShouldBindJSON(&req); err != nil {
		//	context.JSON(http.StatusOK, gin.H{"data": map[string]any{}, "return_message": err.Error(), "return_code": "000001"})
		//	return
		//}

		ch, err := cli.StreamChat(context.Request.Context(), req)
		if err != nil {
			context.JSON(http.StatusOK, gin.H{"data": map[string]any{}, "return_message": err.Error(), "return_code": "000001"})
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

	//================================================================================

	handler.POST("/ai/conversation/clear", func(context *gin.Context) {
		req := &handle.ConversationsClearRequest{}
		if err := context.ShouldBindJSON(req); err != nil {
			context.JSON(http.StatusOK, gin.H{"data": map[string]any{}, "return_message": err.Error(), "return_code": "000001"})
			return
		}
		resp, err := cli.ConversationsClear(context.Request.Context(), req)
		if err != nil {
			_, _ = context.Writer.Write([]byte(err.Error()))
			return
		}
		context.JSON(http.StatusOK, gin.H{"data": map[string]any{"item": resp}, "return_message": "操作成功", "return_code": "000000"})
	})

	server := &http.Server{Addr: fmt.Sprintf(":%d", 8080), Handler: handler}
	fmt.Println(server.ListenAndServe())
}
