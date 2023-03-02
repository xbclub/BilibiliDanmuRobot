package http

import (
	"bili_danmaku/internal/svc"
	"bytes"
	"context"
	gogpt "github.com/sashabaranov/go-gpt3"
	"github.com/zeromicro/go-zero/core/logx"
)

func RequestChatgptRobot(msg string, svcCtx *svc.ServiceContext) (string, error) {
	c := gogpt.NewClient(svcCtx.Config.ChatGPT.APIToken)
	ctx := context.Background()
	msgs := ""
	req := gogpt.ChatCompletionRequest{
		Model: gogpt.GPT3Dot5Turbo,
		Messages: []gogpt.ChatCompletionMessage{
			{
				Role:    "user",
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
