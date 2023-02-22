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

	req := gogpt.CompletionRequest{
		Model:     gogpt.GPT3TextDavinci003,
		MaxTokens: svcCtx.Config.ChatGPT.MaxToken,
		Prompt:    "使用中文回答" + msg,
	}
	resp, err := c.CreateCompletion(ctx, req)
	if err != nil {
		return "", err
	}
	logx.Infof("本次开销：%v tokens", resp.Usage.TotalTokens)
	data := []byte(resp.Choices[0].Text)
	if bytes.HasPrefix(data, []byte{239, 188, 159}) {
		data = bytes.TrimPrefix(data, []byte{239, 188, 159})
	}
	data = bytes.ReplaceAll(data, []byte{10, 10}, []byte{})
	return string(data), nil
}
