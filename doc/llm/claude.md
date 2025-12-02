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
* `#` 添加到CLAUDE.md文件中
* `!` 执行shell命令
* `@` 添加文件到会话中


### 其他
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