## Claude

为什么说`AI cli`工具比`AI IDE`更强大的在于它可以直接操作shell等命令，无需被IDE或环境等限制。

### AI cli工具对比
* Claude code： "自治"的工作流引擎。提供最全的可定制的配置，可以在各种md文件中添加限制，目标等
* Codex cli： OPENAI公司的，在执行任何修改前，都会最清晰的方式展示变动
* Gemini cli： google公司的，得益于google天然搜索及大规模的上下文 
    
---

## Claude 搭建

### 环境准备
```sh
nvm install 22
npm install -g @anthropic-ai/claude-code
claude --version
```

### 替换中国模型
``` ~/.zshrc
export ANTHROPIC_BASE_URL="https://open.bigmodel.cn/api/anthropic"
export ANTHROPIC_AUTH_TOKEN="xxxxxxxxxxx"
```

### 修改配置
``` ~/.claude/settings.json
{
  "env": {
    "ANTHROPIC_DEFAULT_HAIKU_MODEL": "glm-4.5-air",
    "ANTHROPIC_DEFAULT_SONNET_MODEL": "glm-4.6",
    "ANTHROPIC_DEFAULT_OPUS_MODEL": "glm-4.6"
  }
}
```
* HAIKU: 最小模型
* SONNET: 中等模型
* OPUS: 大模型
#### 配置文件路径
* managed-settings.json: 由团队同意发放，优先级最高，不可覆盖 （当前路径）
* 命令行 -model 为临时覆盖，用于快速测试 （命令行）
* .claude/setting.local.json 用于个人项目的偏好 （当前路径）
* .claude/setting.json 所有团队贡献 （当前路径）
* ~/.claude/setting.json 全局配置 （根目录）

### 验证
```shell
claude -p "你好，用中文介绍go的channel,字数为50个字"
```
---

## Claude 使用

### 常用命令
* `/config` 设置
* `/status` 状态
* `/init` 初始化CLAUDE.md文件（用于长期记忆/项目规范）在初始化时会读取项目中的文件分析后生成
* `/memory` 修改CLAUDE.md文件（无需退出claude工具，同时指定哪个claude.md文件）
* `/review` 用于对PR代码升级
* `/pr_comments` 获取并显示氮气的PR所有评论
* `#` 添加到CLAUDE.md文件中
* `!` 执行shell命令
* `@` 添加文件到会话中
#### 会话管理
* `/clear` 清空上下文信息
* `/compact` 上下文生成摘要后，用摘要替换上下文（减少上下文token爆炸问题）
* `rewind` 让对话/操作回到对话历史的某一处
#### 环境配置
* `/config` 查看/修改配置
* `/PERMISSION` 管理ai工具“白名单”
* `/model[model]` 查看或切换模型
#### 元信息
* `/help` 查看命令
* `/status` 查看当前的模型，版本，账户信息等
* `/doctor` 检查claude是否健康，依赖是否完整，配置是否正确
* `/cost和/usage` cost显示当前会话的token消耗。usage显示套餐用例和速率限制
* `/feedback or /bug` 反馈官方的bug

### 自定义命令slash commands
* 命令存放位置
  * Project ./.claude/commands/xxx.md
  * User ~/.claude/commands/xxx.md
* 参数
  * $ARGUMENTS 占位符 根据指令后的参数带到md文件中生成信息提示词
  * $1 $2 $3 ... 对应指令后的参数1，参数2，参数3...带到md文件中生成信息提示词
* 元信息
  * description 这个描述这个工具的作用
  * argument-hint 参数
  * model 模型
  * allowed-tools 可以使用什么工具
```markdown
---
description: 这个描述这个工具的作用
argument-hint: [参数1是什么] [参数2是什么]
model: 使用哪个模型
allowed-tools: Bash(go test:*),Write,Bash(git add:*)
---

请根据`constitution.md`定义的规则进行xxx处理

**当前分支**
!`git branch --show-current`

**当前go版本**
!`go version`
```

### skill
自定义命令slash commands是需要用户手动在对话框中输入`/xxxxx`命令触发，而skill嵌入到提示词中描述这个技能有什么能力
* 处理方式
    1. 启动时加载所有可用的Skills，把元数据加载到提示语中
    2. 当用户提出一个任务时如果匹配到这个技能时，读取该技能下主体Skill文档
    3. 如果Skill文档中引用其他文件，只有在AI执行那一步时才会去加载引用文件
* 命令存放位置
    * Project ./.claude/skills/技能id名称/SKILL.md
    * User ~/.claude/skills/技能id名称/SKILL.md
* 元信息
    * name 描述这个技能的名称
    * description 这个描述这个技能的作用
```markdown
---
name: go语言代码专家
description: 你可以根据前端项目中的代码，分析后生成go的http服务器代码。
---

# go语言代码专家 Skill

## 核心能力
你可以读取前端项目的代码，分析代码中需要优化为网络请求的地方，生成go语言的http服务器代码

## 执行步骤
1. 创建main.go代码
2. 使用gin库进行http服务的搭建
3. 使用mysql的方式进行存储。mysql使用的是gorm的库进行获取。更新。删除等（注意，删除只允许使用update status = delete的方式，而不并不是正在的删除）
4. 检查代码后如果发现有错误就要及时修复，保证可以运行起来的go项目
5. 第三方库使用 go get 的方式获取
6. 代码分为handle.go ,service.go ,resposity.go.models.go这这个层级
```

### Subagent
目前agent无法成为一个超级agent，当你需要操作一个复杂的指令时（提升性能，同时要保证安全）可能本身是冲突的，那就交给多个子agent去做（每个subagent都是独立的上下文不印象其他或父agent）
* 命令
    * 可通过`/agents`命令创建（无需手动编写md文件）   
* 调用
  * `使用go语言写一个http样例 `go-code`使用这个subagent` (强制指定subagent)
* 命令存放位置
    * Project ./.claude/agnets/subagent_id名称.md
    * User ~/.claude/agnets/subagent_id名称.md
* 元信息
    * name 描述这个技能的名称
    * description 这个描述这个技能的作用
    * tools Read,Grep #可选项
    * model opus #可选项 opus,sonnet(默认),haiku,inherit(当前模型)
```markdown
---
name: go语言代码专家
description: 你可以根据前端项目中的代码，分析后生成go的http服务器代码。
tools: Read,Grep
model: inherit
---

# go语言代码专家 

## 核心能力
你可以读取前端项目的代码，分析代码中需要优化为网络请求的地方，生成go语言的http服务器代码

## 执行步骤
1. 创建main.go代码
2. 使用gin库进行http服务的搭建
3. 使用mysql的方式进行存储。mysql使用的是gorm的库进行获取。更新。删除等（注意，删除只允许使用update status = delete的方式，而不并不是正在的删除）
4. 检查代码后如果发现有错误就要及时修复，保证可以运行起来的go项目
5. 第三方库使用 go get 的方式获取
6. 代码分为handle.go ,service.go ,resposity.go.models.go这这个层级
```

### Hooks
`/hooks` 设定当触发某个操作时进行hook
* PreToolUse 工具使用前hook
* PostToolUse 工具时候后hook
* PostToolFailure 工具使用失败后hook
* Notification 当发通知时hook
* UserPromptSubmit 用户的提示词提交时
* SessionStart/SessionEnd 会话开始/结束时
* SubagentStart/SubagentStop subagent开始/结束时
* PreCompact 合并上下文前
* PermissionRequest 权限校验时
* Stop claude关闭时

### MCP
```shell
# -- 分隔符后面用于服务器启动
claude mcp add      --transport stdio --scope project mcp名称  -- python main.py
# transport http使用 add-json的方式添加, Authorization中添加token，${GITHUB_TOKEN}通过环境变量获取避免硬编码
claude mcp add-json --transport http  --scope project mcp名称 '{"type":"http","url":"https://api.xxxxxxxxxx.com/mcp/","headers":{"Authorization":"Bearer ${XXX_TOKEN}"}}' 
```
1. transport
   * http 与远端http通信
   * sse 与远端SSE通信（已废弃）
   * stdio 与本地标准输入输出通信
2. scope
   * user 存放在个人目录（~/.claude.json）
   * project 存放在项目根目录下（./mcp.json）
   * local 本地默认 (~/.claude.json)

#### 验证是否添加成功
* `/mcp` 命令可以查看是否连接成功，提供什么方法
* 在prompts可以使用`mcp__mcp名称__mcp工具名`调用该mcp工具

### Headless
Headless一般指软件工程中没有用户图形界面交互下运行，在claude场景中即不需要打开claude完成一问一答的方式，而是一次性的一问一答
```bash
# --allowedTools使用什么工具  --permission-mode使用什么权限模式  --output-format输出什么格式
claude -p "讲解一下go部署的优势及缺点"  --allowedTools="Bash,Read" --permission-mode acceptEdits --output-format text
# 将cat获取的内容通过管道符的方式传递给claude
cat error.log | claude -p "帮我分析这个日志里面的内容"  
```
--output-format输出响应格式
* text 文本默认（默认）
* json json格式 （可通过jq的工具进行快速获取json中的值）
* stream-json (如果任务很长，希望看到ai的实时进展，每一步都会输出到stdout) 注意：使用stream-json 必须带上 --verbose

### checkpointing
当操作后发现需要撤回之前的操作                                     
* `/rewind` 后选择倒流到哪个操作                              ### 其他
  * Restore code and conversation  回退代码跟对话          #### CLAUDE.md(操作指南)
  * Restore conversation  只回退对话 (代码信息还保留)           由于CLAUDE.md只是用于claude工具，如果工具替换后就需要使用其他的XXX.md。所以行业定义出来AGENTS.md用于存放通过的长期记忆/项目规范，但是目前还没被完全替代，使用方式
  * Restore code  只回退代码（会话之间的信息还保留）                 ```CLAUDE.md
注意                                                  ----- 通用长期记忆/项目规范 -----
  * 不跟踪bash命令的副作用（如 rm -rf ./*）                     @../AGENTS.md
  * 不跟踪外部编辑（claude没有记录）                             ----- Claude长期记忆/项目规范 -----
  * 它不能完全替代git，他们是互补关系而不是替换关系                       [角色]
                                                    你是一个程序员，这是一个前端项目

[基础]
这个是一个web的html项目，使用的是vue3.0框架
```

#### constitution.md(原则契约)
constitution.md拥有绝对的否决权。做什么需要参考这个”宪法“。这个是高度稳定一般不轻易修改
* 使用什么规范
* 必须遵循什么原则

### 权限体系 
AI自动操作与避免AI随便修改的平衡（效率与权限的平衡）
1. Permission Modes （shift+table 或 settings.json中设置defaultMode）
   * default 所有写操作都需要批准
   * plan 类似default，但更倾向ai制定行动计划，而不是直接执行或给出答案（只说不做）
   * acceptEdits 自动批准编辑，无需你批准，但是类似bash这种才需要你批准
   * bypassPermissions 跳过所有权限，自动执行所有操作
2. /permissions 权限规则
   * deny 最高否定权
   * allow 允许
   * ask 当前权限的默认行为
3. /sandbox 沙箱 (隔离读写权限，在当前的沙箱内，可以读写的权限到最大)
   * Sandbox BashTool,with auth-allow in accept edits mode 当你处在acceptEdits这个权限时，在边界内的bash命令不会询问，直接自动执行
   * Sandbox BashTool,with regular permissions 遵循permissions体系，只要没有命中allow规则就需要得到批准
```配置json
{
   "permissions": {
      "allow": [  // 允许
         "Read(README.md)",
         "Bash(go:version)",              // go version 
         "Bash(go:list:*)",               // go list xxx
         "WebFetch(domain:*.baidu.com)"   // xx.baidu.com
      ],
      "deny": [   // 禁止
         "Read(./**/*.md)",        // 相对路径遍历当前路径下所有的*.md文件
         "Read(./.env*)",          // 相对路径
         "Read(~/.ssh/*)",         // 用户主路径
         "Read(/*.json)",          // settings.json所在的目录下的xxx文件
         "Read(//etc/passwd)"      // 文件系统绝对路径
      ],  
      "ask": [    //询问
         "Write",
         "Edit",
         "MultiEdit"
      ],
      "defaultMode": "default"
   },
   "sanbox": {
      "autoAllowBashIfSandboxed": true,   // 自动授权策略
      "enabled": true                     // 启用整个沙箱隔离机制
   }
}
```
---

## 总结

### ./.claude目录
* `commands/` 自定义工具命令
* `skills/` skill技能
* `agnets/` Subagent子专家
* `hooks/` 触发hooks
* `setting.json` 存放claude配置文件信息(权限，模型等)
* `CLAUDE.md` 操作手册，记录项目基本信息等
* `constitution` 规范手册，记录需要使用什么规范生成

---

## 其他
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