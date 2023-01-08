package http

import (
	"bili_danmaku/internal/svc"
	entity "bili_danmaku/internal/types"
	"encoding/json"
	"errors"
	"github.com/avast/retry-go/v4"
	"github.com/go-resty/resty/v2"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
	"time"
)

//func GetDanmuInfo(svcCtx *svc.ServiceContext) (*entity.ResponseBulletInfo, error) {
//	var err error
//	var resp *resty.Response
//	var url = "https://api.live.bilibili.com/xlive/web-room/v1/index/getDanmuInfo?id=" + strconv.Itoa(svcCtx.Config.RoomId) + "&type=0"
//	r := &entity.ResponseBulletInfo{}
//	if resp, err = cli.R().
//		SetHeader("user-agent", userAgent).
//		Get(url); err != nil {
//		logx.Error("请求getDanmuInfo失败：", err)
//		return nil, err
//	}
//	if err = json.Unmarshal(resp.Body(), r); err != nil {
//		logx.Error("Unmarshal失败：", err, "body:", string(resp.Body()))
//		return nil, err
//	}
//
//	return r, nil
//}

func Send(msg string, svcCtx *svc.ServiceContext) error {
	var err error
	var url = "https://api.live.bilibili.com/msg/send"
	var resp *resty.Response
	var respdata *entity.DanmuResp = new(entity.DanmuResp)
	m := make(map[string]string)
	m["bubble"] = "5"
	m["msg"] = msg
	m["color"] = "4546550"
	m["mode"] = "4"
	m["fontsize"] = "25"
	m["rnd"] = strconv.FormatInt(time.Now().Unix(), 10)
	m["roomid"] = strconv.Itoa(svcCtx.Config.RoomId)
	m["csrf"] = CookieList["bili_jct"]
	m["csrf_token"] = CookieList["bili_jct"]
	err = retry.Do(func() error {
		if resp, err = cli.R().
			SetHeader("user-agent", userAgent).
			SetHeader("cookie", CookieStr).
			SetFormData(m).
			Post(url); err != nil {
			logx.Errorf("请求send失败：%v", err)
			return err
		}
		//respdata := entity.DanmuResp{}
		err = json.Unmarshal(resp.Body(), respdata)
		if err != nil {
			//fmt.Println(string(resp.Body()))
			logx.Errorf("send弹幕响应解析失败:%v", err)
			return nil
		}
		if respdata.Code == 10031 {
			logx.Infof("请求send失败:%s", respdata.Msg)
			return errors.New(respdata.Msg)
		}
		return nil
	}, retry.Attempts(3), retry.DelayType(retry.BackOffDelay))
	if err != nil {
		logx.Error(err)
	}
	return nil
}
