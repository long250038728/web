## python MCP处理

### uv环境
MCP官方推荐使用uv管理python工程， 注意*mac 需要执行`brew install coreutils`
```bash
pip install uv==0.5.24
python3.10 -m pip install --upgrade pip
```

---
### 初始化项目
1. 通过uv创建项目，删除默认的hello文件
2. uv按照mcp工具
3. pip按照mcp环境
```bash
uv init mcp_server
cd mcp_server
rm -rf hello.py
uv add "mcp[cli]" -i https://pypi.tuna.tsinghua.edu.cn/simple
pip install mcp
touch server.py
touch sse_server.py
```

---
server
```python
from mcp.server.fastmcp import FastMCP

# 创建mcp server
mcp = FastMCP("mcp_server")

# 编写工具tool
@mcp.tool()
def get_score_by_name(name: str) -> str:
    """根据员工的姓名获取该员工的绩效得分"""
    if not name:
        return "你未输入名称所以无法查询"
    return f"{name}的分数是100"
```

server_sse
```python
import  uvicorn
from fastapi import FastAPI
from starlette.routing import Mount, Route
from starlette.applications import Starlette
from mcp.server.fastmcp import FastMCP
from mcp.server.sse import SseServerTransport
# =====================================

# 创建mcp server
mcp = FastMCP("mcp_server")

# 编写工具tool
@mcp.tool()
def get_score_by_name(name: str) -> str:
    """根据员工的姓名获取该员工的绩效得分"""
    if not name:
        return "你未输入名称所以无法查询"
    return f"{name}的分数是100"
# =====================================


def create_sse_server(mcp :FastMCP):
    transport = SseServerTransport("/messages/")

    async def handle_sse(request):
        async with transport.connect_sse(request.scope, request.receive, request._send) as streams:
            await mcp._mcp_server.run(streams[0],streams[1],mcp._mcp_server.create_initialization_options())
    
    routes = [
        Route("/sse",endpoint=handle_sse),
        Mount("/messages/",app=transport.handle_post_message)
    ]
    return Starlette(routes=routes)
# =====================================
app = FastAPI()
app.mount("/",create_sse_server(mcp))

if __name__ == "__main__":
    uvicorn.run(app,host="0.0.0.0",port=8765)
```

---

## MCP运行
#### cline添加mcp依赖
``` json
{
  "mcpServers": {
    "mcp_python_sse_server": {
      "url": "http://127.0.0.1:8765/sse",  //python server_sse.py
    },
    "mcp_go_sse_server": {
      "url": "http://127.0.0.1:8080/sse",  //go run server_sse.go
    },
    "mcp_python_stdio_server": {
      "command": "uv",
      "args": [
        "run",
        "--with",
        "mcp[cli]",
        "--with-editable",
        "/Users/linlong/Desktop/mcp_server",
        "mcp",
        "run",
        "/Users/linlong/Desktop/mcp_server/server.py"
      ],
    },
    "mcp_go_stdio_server": {
      "command": "/Users/linlong/Desktop/sse/cmd/mcp/mcp" // go build -o mcp main.go
    }
  }
}
```

---

### 自定义mcp client

初始化项目
```bash
uv init mcp_client
cd mcp_client
rm -rf hello.py
uv add "mcp[cli]" -i https://pypi.tuna.tsinghua.edu.cn/simple
pip install mcp
touch client.py
touch sse_client.py
```

client.py
```python
from mcp import ClientSession, StdioServerParameters
from mcp.client.stdio import stdio_client


# Create server parameters for stdio connection
server_params = StdioServerParameters(
    command="uv", # Executable
    args=[
        "run",
        "--with",
        "mcp[cli]",
        "--with-editable",
        "/Users/linlong/Desktop/mcp_server",
        "mcp",
        "run",
        "/Users/linlong/Desktop/mcp_server/server.py"
    ],# Optional command line arguments
    env=None # Optional environment variables
)

async def run():
    async with stdio_client(server_params) as (read, write):
        async with ClientSession(read, write) as session:
            # Initialize the connection
            await session.initialize()

            # List available tools
            tools = await session.list_tools()
            print("Tools:", tools)

            # call a tool
            score = await session.call_tool(name="get_score_by_name",arguments={"name": "张三"})
            print("score: ", score)

if __name__ == "__main__":
    import asyncio
    asyncio.run(run())
```

运行
```bash
uv run mcp_client/client.py
```