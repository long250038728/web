package main

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func serverSSE() {
	// 创建 MCP 服务器
	s := server.NewMCPServer("mcp", "1.0.0")

	// 添加工具
	tool := mcp.NewTool("get_score_by_name",
		mcp.WithDescription("根据员工的姓名获取该员工的绩效得分"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("员工的姓名"),
		),
	)

	// 添加工具处理函数
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name := request.GetArguments()["name"].(string)
		result := fmt.Sprintf("The score of %s is 95", name)
		return mcp.NewToolResultText(result), nil
	})

	fmt.Println(server.NewSSEServer(s).Start("192.168.1.4:8080"))
}
