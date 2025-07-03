package client

import (
	"bufio"
	"context"
	"fmt"
	llm2 "github.com/long250038728/web/cmd/chat/llm"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

type ChatCorn struct {
}

func NewChatCorn() *ChatCorn {
	return &ChatCorn{}
}

func (c *ChatCorn) Chat() *cobra.Command {
	return &cobra.Command{
		Use:   "chat",
		Short: "聊天",
		Run: func(cmd *cobra.Command, args []string) {
			chat()
		},
	}
}

const ExitCommand = "/exit"

func chat() {
	prompt := "你是一个go开发高手,对于用户的问题你都可以精准回答。返回格式为json, json格式如下{\"think\":\"xxxxxx\" ,\"message\":\"xxxx\",\"run\":[\"xxxxx\"]}, json字段解释think返回的是字符串为思考内容,message返回的是字符串为返回的文字信息，run字段返回的是linux的命令数组字符串，可用于linux系统调用"

	assistant := newAssistant()
	chat, err := llm2.NewChat(&llm2.Config{Model: "deepseek-r1:32b", BaseURL: "http://159.75.100.193:6399/v1"}, llm2.NewConversationMemoryLocal(prompt))
	if err != nil {
		assistant.echo(err.Error())
		return
	}
	assistant.echo("你好，欢迎进入我们的聊天对话，输入'/exit' 则退出聊天功能")

	for {
		message := assistant.getInputMessage()

		// 检查是否退出
		if ExitCommand == message {
			assistant.echo("session ended.")
			break
		}

		resp, err := chat.Chat(context.Background(), message)
		if err != nil {
			assistant.echo(err.Error())
			break
		}

		// 输入结果
		assistant.echo(resp)
	}
}

// ============================================

type assistant struct {
	reader *bufio.Reader
}

func newAssistant() *assistant {
	return &assistant{reader: bufio.NewReader(os.Stdin)}
}

func (a *assistant) echo(say string) {
	fmt.Println("机器人助理: ", say)
}

func (a *assistant) getInputMessage() string {
	fmt.Print("> ")
	input, _ := a.reader.ReadString('\n')
	return strings.TrimSpace(input)
}
