## Agent to Agent
1. 由于局限于目前当个agent无法做到一个大而全的功能，需要把功能拆分到各个agent中然后串联起来，在coze及dify中都是实现了自己工作流，但是不能调用其他平台的agent，此时推出了agent2agent
2. 与MCP不同的是，A2A解决的是单个agent无法提供大而全的功能，MCP解决的是agent获取外部资源信息的能力
3. 他们都是采用了client及server两个端。 由client获取server服务中提供的服务能做什么，通过大模型判断是否调用该agent/mcp。server主要提供我能做什么，调用返回对应数据

### client
client通过agentServer暴露的服务器地址，获取各个agent能做什么，然后大模型判断后进行调用
1. 创建A2A card连接（server提供的服务IP:PORT）
2. 把card实例放入A2AClient中
3. A2AClient发送task 
```
import os
import sys

# 添加项目根目录到Python路径
current_dir = os.path.dirname(os.path.abspath(__file__))
project_root = os.path.abspath(os.path.join(current_dir, '..'))
sys.path.insert(0, project_root)

from common.client import A2AClient, A2ACardResolver
from common.types import TaskState, Task
from common.utils.push_notification_auth import PushNotificationReceiverAuth

import asyncio
from uuid import uuid4
import urllib

async def main():
    card_resolver = A2ACardResolver("http://localhost:10000")
    card = card_resolver.get_agent_card()

    payload = {
        "id": 1,
        "sessionId": 1,
        "acceptedOutputModes": ["text"],
        "message": {
            "role": "user",
            "parts": [
                {"type": "text","text": "张三的绩效是多少分"}
            ]
        }
    }

    client = A2AClient(agent_card=card)
    ret = client.send_task(payload=payload)
    print(ret.model_dump_json())
    
if __name__ == "__main__":
    asyncio.run(main())
```

### server
server中主要各个对象中主要为了描述当前agent有什么作用(通过agent_card)，提供服务让agentClient调用(通过take_manager)
1. 创建一个server服务,暴露IP:PORT
   * server
      * agent_card =>  Capabilities && skill
      * task_manager => agent
      * host 
      * port

```
from common.server import A2AServer
from common.types import AgentCard, AgentCapabilities, AgentSkill, MissingApiKeyError
from common.utils.push_notification_auth import PushNotificationSenderAuth
from agents.langgraph.task_manager import AgentTaskManager
from agents.langgraph.agent import CurrencyAgent
import click
import os
from dotenv import load_dotenv
from openai import OpenAI

class CurrencyAgent:
    SYSTEM_INSTRUCTION = ["text","text/plain"]    
    # 提供单次调用
    def invoke(self,query,sessionId) -> Str:
        client = OpenAI()
        final_answer = clientchat.completions.create(model="deepseek-r1",message=[{"role":"user","content":query}]).choices[0].message.content
        return {
            "is_task_complete": True,
            "require_user_input": False,
            "content": final_answer
        }
     
     # 提供stearm调用
    async def stream(self,query,sessionId) ->AsyncIterable[Dic[Str,Any]]:
        client = OpenAI()
        final_answer = clientchat.completions.create(model="deepseek-r1",message=[{"role":"user","content":query}]).choices[0].message.content
         yield {
           "is_task_complete": True,
            "require_user_input": False,
            "content": final_answer
        }
         

@click.command()
@click.option("--host","host",default="localhost")
@click.option("--port","port",default=10000)
def main(host,port):
    try:
        capabilities = AgentCapabilities(streaming=True,pushNotifications=True)
        skill = AgentSkill(
            id="skill1",
            name="员工绩效工具",
            description="查询员工的绩效信息",
            tags=["查询员工的绩效信息"],
            examples=["张三的几下是多少"],
            )
        agent_card = AgentCard(
            name = "员工绩效工具",
            description="查询员工的绩效信息"，
            url= f"http://{host}:{port}/"
            version="1.0.0.0"
            defaultInputModes=CurrencyAgent.SUPPORTED_CONTENT_TYPES,
            defaultOutputModes=CurrencyAgent.SUPPORTED_CONTENT_TYPES,
            capbilities=capabilities,
            skill =[skill],
        )
        
        notification_sender_auth = PushNotificationSenderAuth()
        notification_sender_auth.generate_jwk()
        server = A2AServer(
            agent_card=agent_card,
            task_manager= AgentTaskManger(agent=CurrencyAgent(),notification_serder_auth=notification_sender_auth)
            host=host,
            port=port,
        )
        server.app_add_router(
            "/.well-know/jwks.json",notification_sender_auth.handle_jwks_endpoint, methods=["GET"]
        )
        server.start()
    except MissingApiKeyError as e:
        exit(1)
    except Exception as e:
        exit(1)
        
if __name__ == "__main__":
    main()           
```
