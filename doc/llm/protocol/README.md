### SSE
1. 通过http请求进行通信，后续服务端客户往该链接发送数据(客户端收)
2. 当客户端需要给服务器端发送数据时，需要使用新的http请求，为了与之前的请求关联起来，所以需要带上请求1的session id等信息

### websocket
1. 无法通过http请求进行通信(使用ws请求进行通讯，此时需要使用特定可websocket库进行操作)，，后续服务端客户往该链接发送数据(客户端收)
2. 当客户端需要给服务器端发送数据时，使用之前已经ws链接进行发送，他们共用同一个链接

### SSE 与 websocket 区别及使用场景
1. SSE: 接收跟发送数据是两个独立的链接(SSE是单向的),适用于服务器需要连续推送数据给客户端的场景
2. websocket: 接收跟发送数据是同一个链接(websocket是双向的)，使用服务器与客户端需要互相发送数据的场景


### 流程
1. 在每个连接创建一个自定义session对象， 放到一个全局的列表中，通过session的chan通道获取最新消息
2. 发生数据是找到对应的session对象，往chan通道进行发送
3. 所以sse与websocket他们只要不断开就会占用连接，请求多的话连接会消耗大量内存及TCP连接

---

###  SSE 示例
1. 设置http响应头信息
```
w.Header().Set("Content-Type", "text/event-stream")
w.Header().Set("Cache-Control", "no-cache")
w.Header().Set("Connection", "keep-alive")
w.Header().Set("Access-Control-Allow-Origin", "*")
```

2. 发送数据给客户端
```
flusher, ok := w.(http.Flusher)
if !ok {
    http.Error(w, "SSE not supported", http.StatusInternalServerError)
    return
}
fmt.Fprintf(w, "message")
flusher.Flush()
```

3. sse断开处理
```
select {
    case <-r.Context().Done():
        return
}
```

4. 服务端接收数据
使用新的http请求进行通信，通过请求参数中的session_id等参数找到对应的sse连接


###  websocket 示例
1. 创建工具
```
upgrader := websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源的连接，生产环境中应该更严格
	},
}
```

2. 升级HTTP连接为WebSocket连接
```
conn, err := h.upgrader.Upgrade(w, r, nil)
if err != nil {
	fmt.Printf("Failed to upgrade connection from %s: %v\n", r.RemoteAddr, err)
	http.Error(w, fmt.Sprintf("Could not upgrade connection: %v", err), http.StatusInternalServerError)
	return
}
defer conn.Close()
```

3. 发送数据给客户端
```
if err := conn.WriteJSON(map[string]interface{}{"type": "session", "session_id": sess.ID}); err != nil {
	return
}
if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
	return
}
```

4. 服务端接收数据
```
for {
    messageType, p, err := conn.ReadMessage()
}
```