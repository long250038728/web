# plugin

## SDD
### spec/plan/task的关系
* spec.md 以产品经理的视角分析
  * 这个一个什么产品，产品的定位是什么
  * 产品的功能是什么，实现了用户的什么目标
  * 用户群体是什么，风格是什么
  * 有什么是功能是必须有，什么功能是一定不能有，有什么页面，每个页面具体实现什么功能
* plan.md 以架构师的视角分析
  * 使用什么技术栈，需要遵循开发中的什么规范及原则
  * 架构图分层设计
* task.md 以开发者的视角分析
  * 分析及明确实现的步骤及任务（步骤1做什么，步骤2做什么）
  * 主要注意的事项


## superpowers （解决怎么做得正确）- 场景单线程完成
防止 agent 直接跳入写代码，强制执行开发方法论（设计 → 计划 → TDD → 代码审查 → 合并），核心目标是质量和规范。
```
#  Anthropic marketplace 默认已经加入，所以不需要 /plugin marketplace add xxxx
/plugin install superpowers@claude-plugins-official
```
1. /superpowers:brainstorming 需求:  进行头脑风暴(澄清需求) （生成docs/superpowers/spec/2026-04-28-weight-edit.md）
2. /superpowers:writing-plans: 编写实现计划（生成docs/superpowers/plans/2026-04-28-weight-edit.md）
3. /superpowers:subagent-driven-development: 执行实现&&代码审查(子代理)

## oh-my-claude-code （解决怎么做得更快）- 场景多代理并行处理
把一个 agent 变成多个 agent，通过并行执行、角色分离、验证循环和可观测性（HUD、成本统计、状态追踪）来工作，
```
/plugin marketplace add https://github.com/Yeachan-Heo/oh-my-claudecode
/plugin install oh-my-claudecode 
```
* /oh-my-claudecode:setup 初始化项目
  1. 初始化claude.md
  2. 初始化skill
  3. 安装mcp & omc插件
* /oh-my-claudecode:plan  帮我设计一个可以用vue的登录注册修改密码，登录后带上access_token ，后端使用go语言
  1. 询问需求澄清
  2. 创建plans计划 (.omc/plans/vue-go-auth.md)
* /oh-my-claudecode:ralph 根据计划执行（只有执行完才会停止）同时如果中断了重新进入可以继续执行该命令继续完成
  1. 创建`.omc/state/sessions/ralph-vue-go-auth/prd.json`  (知道有哪些 Story)
  2. 创建`.omc/state/sessions/ralph-vue-go-auth/progress.txt`  （知道哪些已完成）
* /oh-my-claudecode:ulw 自动拆解任务，并行派出多个 specialist agent


### setup做的内容
```
  新增文件：

  ┌───────────────────────────────────────┬──────────────────────────────────────────────────────┐
  │                 文件                  │                         说明                         │
  ├───────────────────────────────────────┼──────────────────────────────────────────────────────┤
  │ .claude/CLAUDE.md                     │ OMC 核心行为配置，让 Claude 具备并行调度、规划等能力 │
  ├───────────────────────────────────────┼──────────────────────────────────────────────────────┤
  │ .claude/skills/omc-reference/SKILL.md │ OMC 技能参考文档                                     │
  ├───────────────────────────────────────┼──────────────────────────────────────────────────────┤
  │ ~/.claude/hud/omc-hud.mjs             │ HUD 状态栏脚本                                       │
  ├───────────────────────────────────────┼──────────────────────────────────────────────────────┤
  │ ~/.claude/hud/lib/config-dir.mjs      │ HUD 依赖模块                                         │
  ├───────────────────────────────────────┼──────────────────────────────────────────────────────┤
  │ ~/.claude/.omc-config.json            │ OMC 偏好配置（执行模式、任务工具、团队配置）         │
  └───────────────────────────────────────┴──────────────────────────────────────────────────────┘

  修改的文件：

  ┌──────────────────────────┬─────────────────────────────────────────────────────────────────────┐
  │           文件           │                              改动内容                               │
  ├──────────────────────────┼─────────────────────────────────────────────────────────────────────┤
  │ ~/.claude/settings.json  │ 新增 statusLine（HUD）和 env.CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1 │
  ├──────────────────────────┼─────────────────────────────────────────────────────────────────────┤
  │ ~/.claude.json（项目级） │ 新增 Context7 MCP server 配置                                       │
  └──────────────────────────┴─────────────────────────────────────────────────────────────────────┘

```

## everything-claude-code （怎么保证质量和安全不出错）
在代码生成中加载对应语言的技能获取规范,常见的是自动触发无需手动调用
```
/plugin marketplace add https://github.com/affaan-m/everything-claude-code
/plugin install everything-claude-code@everything-claude-code
```

