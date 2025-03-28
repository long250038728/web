from mcp import ClientSession, StdioServerParameters
from mcp.client.stdio import stdio_client

async def connect_to_sse_server(server_url: str):
    """Connect to an MCP server running with SSE transport"""
    # Store the context managers so they stay alive
    async with sse_client(url=server_url) as (read, write):
        async with ClientSession(read, write) as session:
            await session.initialize()
            # List available tools to verify connection
            print("Initialized SSE client...")
            print("Listing tools...")
            response = await session.list_tools()
            tools = response.tools
            print("\nConnected to server with tools:", [tool.name for tool in tools])

            # call a tool
            score = await session.call_tool(name="get_score_by_name",arguments={"name": "张三"})

            print("score: ", score)