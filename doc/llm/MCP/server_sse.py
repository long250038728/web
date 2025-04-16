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