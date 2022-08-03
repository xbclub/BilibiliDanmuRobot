package http

import (
	entity "bili_danmaku/internal/types"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/zeromicro/go-zero/core/logx"
)

// 调用青云客机器人api
func RequestQingyunkeRobot(msg string) (string, error) {
	var err error
	var url = "http://api.qingyunke.com/api.php?key=free&appid=0&msg=" + msg
	var resp *resty.Response

	if resp, err = cli.R().
		SetHeader("Content-Type", "utf-8").
		Get(url); err != nil {
		logx.Error("请求qingyunke机器人接口失败：", err)
		return "", err
	}

	r := &entity.QinugyunkeRobotReplay{}
	err = json.Unmarshal(resp.Body(), r)

	return r.Content, err
}
