package http

import (
	"bili_danmaku/internal/errs"
	"bili_danmaku/internal/svc"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/zeromicro/go-zero/core/logx"

	//entity "bili_danmaku/internal/types"
	entity "bili_danmaku/internal/types"
	"strconv"
)

func RoomInit(svcCtx *svc.ServiceContext) (*entity.RoomInitInfo, error) {
	var err error
	var resp *resty.Response
	var url = "https://api.live.bilibili.com/room/v1/Room/room_init?id=" + strconv.Itoa(svcCtx.Config.RoomId)

	if resp, err = cli.R().
		SetHeader("user-agent", userAgent).
		Get(url); err != nil {
		logx.Error("请求room_init失败：", err)
		return nil, err
	}

	// 先解析响应状态
	status := &entity.RoomInitStatus{}
	if err = json.Unmarshal(resp.Body(), status); err != nil {
		logx.Error("Unmarshal失败：", err, "body:", string(resp.Body()))
		return nil, err
	}

	// 在解析房间状态
	r := &entity.RoomInitInfo{}
	if status.Code == 0 {
		if err = json.Unmarshal(resp.Body(), r); err != nil {
			logx.Error("Unmarshal失败：", err, "body:", string(resp.Body()))
			return nil, err
		}
	}

	// 太长时间下播，房间号可能会消失，请求响应的code=60004
	if status.Code == 60004 {
		return nil, errs.RoomIdNotExistErr
	}
	return r, err
}
