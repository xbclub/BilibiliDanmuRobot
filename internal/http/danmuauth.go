package http

import (
	types2 "bili_danmaku/internal/types"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/zeromicro/go-zero/core/logx"
)

func GetDanmuToken(roomid int) (danmuAuthDatas *types2.DanmuAuthData, err error) {
	var url = fmt.Sprintf("https://api.live.bilibili.com/xlive/web-room/v1/index/getDanmuInfo?id=%v", roomid)
	var resp *resty.Response
	if resp, err = cli.R().
		SetHeader("user-agent", userAgent).
		Get(url); err != nil {
		logx.Error("弹幕流秘钥获取失败：", err)
		return nil, err
	}

	// 先解析响应状态
	danmuAuthDatas = &types2.DanmuAuthData{}
	if err = json.Unmarshal(resp.Body(), danmuAuthDatas); err != nil {
		logx.Error("Unmarshal失败：", err, "body:", string(resp.Body()))
		return nil, err
	}
	if danmuAuthDatas.Code != 0 {
		logx.Error(danmuAuthDatas.Message)
		return nil, errors.New(danmuAuthDatas.Message)
	}
	return danmuAuthDatas, nil
}
