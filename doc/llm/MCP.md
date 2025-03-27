## uv环境
MCP官方推荐使用uv管理python工程
```bash
pip install uv=0.5.24
```

## 初始化项目
1. 删除默认hello文件
2. 添加mcp及创建.venv环境安装依赖
```bash
uv init myproject
cd myproject

rm -rf hello.py

uv add "mcp[cli]"
pip install mcp
```

## tool编写
自定义的请求处理流程
```python
from mcp.server.fastmcp import FastMCP

def handle_request(request):
    """
    处理客户端请求的函数
    :param request: 客户端发送的请求数据
    :return: 响应字符串
    """
    print("收到请求：", request)
    # 添加你的业务逻辑处理，这里返回简单字符串作为示例
    return "Hello from FastMCP!"

if __name__ == "__main__":
    # 创建 FastMCP 实例，设置监听地址、端口和请求处理函数
    server = FastMCP(host="0.0.0.0", port=8080, handler=handle_request)
    print("FastMCP 服务器启动，监听端口 8080 ...")
    # 启动服务器，开始接受请求
    server.start()
```

注册工具供MCP调用
```python
from mcp.server.fastmcp import FastMCP

# 创建mcp server
mcp = FastMCP("myproject")

# 编写工具tool
@mcp.tool()
def get_score_by_name(name: str) -> str:
    """根据员工的姓名获取该员工的绩效得分"""
    return f"{name}的分数是100"
```

## mcp加载json
```json
{
  "mcpServers": {
    "myproject": {
      "command": "uv",
      "args": [
        "run",
        "--with","mcp[cli]",
        "--with-editable","/Users/linlong/Desktop/myproject",
        "mcp","run","/Users/linlong/Desktop/myproject/server.py"
       ] 
  }
}
```

---

## `uv` 命令解析

### 命令
```sh
uv run --with mcp[cli] --with-editable /Users/linlong/Desktop/myproject mcp run /Users/linlong/Desktop/myproject/server.py
```

### 解析
#### 1. `uv`
`uv` 是一个 Python 的包管理和运行工具，类似于 `pip` 和 `venv`，但更现代、高效，支持即时运行和依赖管理。

#### 2. `run`
- `run` 是 `uv` 的一个子命令，用于执行 Python 代码或脚本，并在一个受管理的环境中运行它。

#### 3. `--with mcp[cli]`
- `--with` 用于安装或使用指定的 Python 依赖。
- `mcp[cli]` 指定了 `mcp` 这个 Python 包，并额外包含 `cli` 这个可选的依赖（如果 `mcp` 采用了 extras 机制，即 `setup.py` 或 `pyproject.toml` 中定义的 `[cli]` 额外功能）。

#### 4. `--with-editable /Users/linlong/Desktop/myproject`
- `--with-editable` 允许以“可编辑模式” (`pip install -e .`) 方式安装本地项目。
- `/Users/linlong/Desktop/myproject` 是要安装的项目路径，意味着 `uv` 会在运行环境中使用这个项目的代码，而不需要每次修改代码后重新安装。

#### 5. `mcp run /Users/linlong/Desktop/myproject/server.py`
- `mcp` 可能是一个 Python 可执行程序（比如 `mcp` 是 `mcp` 库提供的 CLI 工具）。
- `run` 可能是 `mcp` 的一个子命令，意味着 `mcp` 需要运行某个 Python 脚本。
- `/Users/linlong/Desktop/myproject/server.py` 是要运行的 Python 代码文件。

### 总结
该命令的作用是：
1. 使用 `uv` 运行 Python 环境。
2. 在环境中动态安装 `mcp[cli]` 依赖包。
3. 以“可编辑模式”安装 `/Users/linlong/Desktop/myproject`，确保代码修改可即时生效。
4. 执行 `mcp run /Users/linlong/Desktop/myproject/server.py`，可能是通过 `mcp` 这个 CLI 工具来启动 `server.py`。

这个命令通常用于开发环境，确保运行的是最新代码，同时减少手动安装和配置的工作量。