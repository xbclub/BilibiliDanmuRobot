package config

import (
	"github.com/zeromicro/go-zero/core/logx"
)

type Config struct {
	//rest.RestConf
	Log         logx.LogConf
	RoomId      int    `json:",default=4699397"`
	WsServerUrl string `json:",default=wss://broadcastlv.chat.bilibili.com:2245/sub"`
	//QrCodePath   string `json:"qr_code_path"`
	TalkRobotCmd string `json:",default=test"`
	RobotName    string `json:",default=花花"`
	DanmuLen     int    `json:",default=20"`
	EntryMsg     string `json:",default=off"`
}
