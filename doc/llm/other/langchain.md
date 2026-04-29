## langchain环境搭建

### 项目创建
```bash
mkdir langchain_v1_project && cd langchain_v1_project
uv init
uv python pin 3.10
uv venv
source .venv/bin/activate
```

### 添加依赖库
```bash
uv add langchain langchain-openai
uv add langgraph
```

### 添加环境变量
```
export OPENAI_API_BASE="https://ark.cn-beijing.volces.com/api/v3/"
export OPENAI_BASE_URL="https://ark.cn-beijing.volces.com/api/v3/"
export OPENAI_API_KEY="xxxxxxxxxxxxxxx"
```


### 代码示例
```python
from dataclasses import dataclass
from pydantic import BaseModel,Field
from typing import Literal,Any

from langchain_openai import ChatOpenAI

from langchain.agents import create_agent
from langchain.agents.structured_output import ToolStrategy
from langchain.tools import tool,ToolRuntime
from langchain.agents.middleware import before_model, after_model, AgentState,AgentMiddleware

from langchain.messages import AIMessage
from langgraph.runtime import Runtime


# 上下文对象
@dataclass
class Context:
    """当前runtime的上下文"""
    city: str

# 响应对象
@dataclass
class ResponseFormat:
    """Response schema for the agent."""
    # A punny response (always required)
    punny_response: str
    # Any interesting information about the weather if available
    weather_conditions: str | None = None

# from pydantic import BaseModel
# class ResponseFormat(BaseModel):
#     punny_response: str
#     weather_conditions: str | None = None


@tool
def get_weater(city: str) -> str:
    """获取城市天气情况"""
    return f"{city}今天天气不错"

@tool("get_ctiy",description="获取城市名称")
def get_ctiy(runtime: ToolRuntime[Context]) -> str:
    """获取城市"""
    return runtime.context.city #通过上下文中的参数获取

# 参数定义类对象
class User(BaseModel):
    name: str = Field(description="会员姓名")
    age: int = Field(description="会员年龄")
    sex: Literal["男", "女"] = Field(description="会员性别，值为男或者女" ,default="男")

@tool(args_schema=User)
def get_user_description(name: str,age: int,sex: str)->str:
    """根据用户信息得到别人给这个人的总结描述"""
    return f"这是一个年龄{age}性别{sex},名字叫做{name}"

# xxxx_model每次调用model执行一次  xxxx_agent代理开始之前/后调用一次
class LoggingMiddleware(AgentMiddleware):
    def before_model(self, state: AgentState, runtime: Runtime) -> dict[str, Any] | None:
        print(f"About to call model with {len(state['messages'])} messages")
        return None

    def after_model(self, state: AgentState, runtime: Runtime) -> dict[str, Any] | None:
        print(f"Model returned: {state['messages'][-1].content}")
        return None


def main():
    model = ChatOpenAI(model="doubao-1-5-pro-32k-250115")

    agent = create_agent(
        model,
        tools=[get_weater,get_ctiy,get_user_description],
        system_prompt="你是一位幽默有趣的天气预报播报员",
        context_schema=Context,
        response_format=ToolStrategy(ResponseFormat),
        middleware=[LoggingMiddleware()]
    )
    
    resp = agent.invoke(
        {
         "messages": [{"role":"user","content":"我的名字叫小红，性别女，年龄24，今天需要穿什么衣服呢，结合天气情况，同时我想知道我的姓名年龄性别，别人给我的描述是什么"}]
        },   
        context=Context(city="北京"),                                #设置上下文信息
        config= {"configurable": {"thread_id": "1"}}                #配置
    )
    print(resp)

if __name__ == "__main__":
    main()
```