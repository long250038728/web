package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sashabaranov/go-openai"
)

//go get github.com/sashabaranov/go-openai

type AddTool struct {
}

func (t *AddTool) Name() string {
	return "AddTools"
}

func (t *AddTool) Function() *openai.FunctionDefinition {
	return &openai.FunctionDefinition{
		Name:        t.Name(),
		Description: ` Use this tool for addition calculations. example: 1+2 =? then Action Input is: 1,2 `,
		Parameters:  `{"type":"object","properties":{"numbers":{"type":"array","items":{"type":"integer"}}}}`,
	}
}

func (t *AddTool) Sum(call *openai.ToolCall) (int, error) {
	// 执行加法操作
	sum := 0

	// 工具的参数应该在 `call.Function.Parameters`，解析参数
	var params struct {
		Numbers []int `json:"numbers"`
	}
	// 如果工具参数为空或者解析失败，返回错误
	if call.Function.Arguments == "" {
		return sum, fmt.Errorf("tool parameters are empty")
	}
	err := json.Unmarshal([]byte(call.Function.Arguments), &params)
	if err != nil {
		return sum, fmt.Errorf("failed to parse tool parameters: %w", err)
	}

	for _, num := range params.Numbers {
		sum += num
	}

	return sum, nil
}

//====================================================================================

type QwChat struct {
	client  *openai.Client
	message []openai.ChatCompletionMessage
}

type Opt func(c *QwChat)

func SetMessage(message []openai.ChatCompletionMessage) Opt {
	return func(c *QwChat) {
		c.message = message
	}
}

func NewOpenAiClient(opts ...Opt) *QwChat {
	config := openai.DefaultConfig("sk-bff318c2fa9e4eceb6c292e2990f0dfc")
	config.BaseURL = "https://dashscope.aliyuncs.com/compatible-mode/v1"
	return &QwChat{
		client:  openai.NewClientWithConfig(config),
		message: []openai.ChatCompletionMessage{},
	}
}

func (c *QwChat) Chat(ctx context.Context, msg string) (string, error) {
	c.message = append(c.message, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: msg,
	})
	add := AddTool{}

	resp, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:      "qwen-plus",
			Messages:   c.message,
			ToolChoice: "auto", // 工具选择方式让大模型自己根据实际情况选择是否调用工具
			Tools: []openai.Tool{
				{
					Type:     "function",
					Function: add.Function(),
				},
			},
		},
	)
	if err != nil {
		return "", err
	}

	toolCalls := resp.Choices[0].Message.ToolCalls
	if toolCalls != nil {
		// 遍历工具调用
		for _, call := range toolCalls {
			// 判断是否为函数类型工具
			if call.Type == "function" {
				// 判断调用的工具名称是否为 "AddTools"
				if call.Function.Name == add.Name() {

					val, err := add.Sum(&call)
					if err != nil {
						return "", err
					}
					// 将工具的执行结果作为对话消息返回
					toolResponse := fmt.Sprintf("The result of addition is: %d", val)
					c.message = append(c.message, openai.ChatCompletionMessage{
						Role:    openai.ChatMessageRoleAssistant,
						Content: toolResponse,
					})

					// 返回结果
					return toolResponse, nil
				}
			}
		}

		// 如果工具调用未被处理
		return "", fmt.Errorf("no valid tool call handled")
	}

	c.message = append(c.message, resp.Choices[0].Message)
	return resp.Choices[0].Message.Content, nil
}

func (c *QwChat) ChatStream(ctx context.Context, msg string) (chan string, error) {
	ch := make(chan string, 5)

	c.message = append(c.message, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: msg,
	})

	resp, err := c.client.CreateChatCompletionStream(
		ctx,
		openai.ChatCompletionRequest{
			Model:    "qwen-plus",
			Messages: c.message,
		},
	)
	if err != nil {
		return nil, err
	}

	go func() {
		defer resp.Close()
		for {
			receivedResponse, streamErr := resp.Recv()
			if streamErr != nil {
				close(ch)
				return
			}
			ch <- receivedResponse.Choices[0].Delta.Content
		}
	}()
	return ch, nil
}
