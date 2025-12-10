## Claude
为什么说`AI cli`工具比`AI IDE`更强大的在于它可以直接操作shell等命令，无需被IDE或环境等限制。

### AI cli工具对比
* Claude code： "自治"的工作流引擎。提供最全的可定制的配置，可以在各种md文件中添加限制，目标等
* Codex cli： OPENAI公司的，在执行任何修改前，都会最清晰的方式展示变动
* Gemini cli： google公司的，得益于google天然搜索及大规模的上下文 

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


### 自定义命令
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


### 其他
#### CLAUDE.md(操作指南)
由于CLAUDE.md只是用于claude工具，如果工具替换后就需要使用其他的XXX.md。所以行业定义出来AGENTS.md用于存放通过的长期记忆/项目规范，但是目前还没被完全替代，使用方式
```CLAUDE.md
----- 通用长期记忆/项目规范 -----
@../AGENTS.md
----- Claude长期记忆/项目规范 -----
[角色]
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
```
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
   
   