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
	TalkRobotCmd       string            `json:",default=test"`
	FuzzyMatchCmd      bool              `json:",default=false"`
	RobotName          string            `json:",default=花花"`
	DanmuLen           int               `json:",default=20"`
	EntryEffect        bool              `json:",default=false"`
	EntryMsg           string            `json:",default=off"`
	WelcomeDanmu       []string          `json:",default='欢迎 {user} ~'"`
	ThanksGift         bool              `json:",default=false"`
	CustomizeBullet    bool              `json:",default=false"`
	InteractWord       bool              `json:",default=false"`
	InteractWordByTime bool              `json:",default=false"`
	WelcomeSwitch      bool              `json:",default=false"`
	WelcomeString      map[string]string `json:",optional"`
	RobotMode          string            `json:",default=QingYunKe,options=QingYunKe|ChatGPT"`
	ChatGPT            struct {
		APIToken string `json:",optional"`
	}
	CronDanmu     bool `json:",default=false"`
	CronDanmuList []struct {
		Cron   string   `json:",optional"`
		Random bool     `json:",default=false"`
		Danmu  []string `json:",optional"`
	} `json:",optional"`
	WelcomeDanmuByTime []struct {
		Enabled bool     `json:",optional"`
		Key     string   `json:",optional"`
		Random  bool     `json:",default=false"`
		Danmu   []string `json:",optional"`
	} `json:",optional"`
	FocusDanmu           []string `json:",optional"`
	PKNotice             bool     `json:",default=true"`
	WelcomeBlacklistWide []string `json:",optional"`
	WelcomeBlacklist     []string `json:",optional"`
}
