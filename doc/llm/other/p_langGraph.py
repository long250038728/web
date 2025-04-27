from typing import TypedDict, Annotated
from langgraph.graph import START,END,StateGraph

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



def main():
    graph = StateGraph(State)
    graph.add_node("a",a)
    graph.add_node("b",b)
    graph.add_node("c",c)

    # 单边连接
    graph.add_edge(START,"a")
    graph.add_edge("a","b")
    graph.add_edge("b","c")
    graph.add_edge("c",END)

    # 条件边连接
    # graph.add_conditional_edges(...)

    app = graph.compile()

    #生成图片
    app.get_graph().draw_mermaid_png(output_file_path="flowchart.png")

    # 提问
    print(app.invoke({"queue": ["How old are you"]}))



if __name__ == '__main__':
    main()