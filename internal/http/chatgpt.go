package http

import (
	"bili_danmaku/internal/svc"
	"bytes"
	"context"
	"strconv"
	"strings"

	gogpt "github.com/sashabaranov/go-openai"
	"github.com/zeromicro/go-zero/core/logx"
)

func RequestChatgptRobot(msg string, svcCtx *svc.ServiceContext) (string, error) {
	// c := gogpt.NewClient(svcCtx.Config.ChatGPT.APIToken)
	cfg := gogpt.DefaultConfig(svcCtx.Config.ChatGPT.APIToken)
	cfg.BaseURL = svcCtx.Config.ChatGPT.APIUrl
	c := gogpt.NewClientWithConfig(cfg)
	ctx := context.Background()
	msgs := ""
	prompt := strings.ReplaceAll(svcCtx.Config.ChatGPT.Prompt, "{limit}", strconv.Itoa(svcCtx.Config.DanmuLen))
	req := gogpt.ChatCompletionRequest{
		Model: gogpt.GPT3Dot5Turbo0613,
		Messages: []gogpt.ChatCompletionMessage{
			{
				Role: gogpt.ChatMessageRoleAssistant,
				//Content: fmt.Sprintf("你是一个非常幽默的机器人助理，尽可能的在%v个字符内回答，不要使用emoji等表情符号，可以使用颜文字", svcCtx.Config.DanmuLen),
				Content: prompt,
			},
			{
				Role:    gogpt.ChatMessageRoleUser,
				Content: msg,
			},
		},
	}
	resp, err := c.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}
	logx.Infof("本次开销：%v tokens", resp.Usage.TotalTokens)
	for _, v := range resp.Choices {
		data := []byte(v.Message.Content)
		if bytes.HasPrefix(data, []byte{239, 188, 159}) {
			data = bytes.TrimPrefix(data, []byte{239, 188, 159})
		}
		data = bytes.ReplaceAll(data, []byte{10, 10}, []byte{})
		msgs += string(data)
	}
	return msgs, nil
}
