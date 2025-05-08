# uv add langgraph
# uv add pyppeteer

from typing import TypedDict, Annotated
from langgraph.graph import START,END,StateGraph
from langchain_core.runnables.graph_mermaid import MermaidDrawMethod

def add_messages(old: list[str], new: list[str]) -> list[str]:
    return old + new

class State(TypedDict):
    queue:Annotated[list[str], add_messages]
    result:Annotated[list[str], add_messages]


# =================================================================


def a(state):
    print(state,"a")
    return {"result": ["this is a"]}

def b(state):
    print(state,"b")
    return {"result":["this is b"]}

def c(state):
    print(state,"c")
    return {"result":["this is c"]}

def conditional_state(state):
    if len(state["result"]) < 5:
        return "LOOP"
    return "NEXT"

def main():
    graph = StateGraph(State)
    graph.add_node("a",a)
    graph.add_node("b",b)
    graph.add_node("c",c)

    # 单边连接
    graph.add_edge(START,"a")


    # 条件边连接
    graph.add_conditional_edges(
        "a",                           # 上一节点 a
        conditional_state,                    # 逻辑判断会返回key
        {                           # dict是key value格式 根据返回的key选择是哪个node
            "LOOP":"b",
            "NEXT":"c"
        }
    )

    graph.add_edge("b", "a") # b 执行完后回到 a 形成循环

    graph.add_edge("c",END)

    app = graph.compile()

    #生成图片
    app.get_graph().draw_mermaid_png(output_file_path="flowchart.png",draw_method=MermaidDrawMethod.PYPPETEER )

    # 提问
    print(app.invoke({"queue": ["How old are you"]}))



if __name__ == '__main__':
    main()