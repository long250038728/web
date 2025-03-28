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
