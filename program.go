package main

import (
	"context"
	"github.com/xbclub/BilibiliDanmuRobot-Core/config"
	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
	"github.com/xbclub/BilibiliDanmuRobot-Core/handler"
	"github.com/xbclub/BilibiliDanmuRobot-Core/http"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
	"github.com/xbclub/BilibiliDanmuRobot-Core/utiles"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type Program struct {
	running    bool
	svcCtx     *svc.ServiceContext
	workCtx    context.Context
	workCancel context.CancelFunc
	cls        handler.WsHandler
}

func NewProgram() *Program {
	return &Program{}
}

func (p *Program) Start() bool {
	if p.running {
		return true
		logx.Info("已启动跳过")
	}
	var c config.Config
	if err := Mustload(&c); err != nil {
		logx.Error(err)
		return false
	}
	logx.MustSetup(c.Log)
	logx.DisableStat()
	p.svcCtx = svc.NewServiceContext(c)
	p.workCtx, p.workCancel = context.WithCancel(context.Background())
	p.running = true
	go p.Bili_danmaku_Start(p.workCtx)
	return true
}

func (p *Program) Stop() bool {
	if p.workCancel != nil {
		p.workCancel()
	}
	return true
}

func (p *Program) Monitor() bool {
	return p.running
}

func (l *Program) Bili_danmaku_Start(workctx context.Context) {
	if !http.FileExists("token/bili_token.txt") || !http.FileExists("token/bili_token.json") {
		return
	}
	var err error
	http.InitHttpClient()

	// 准备select中用到的变量
	//sig := make(chan os.Signal)
	//signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	var interval = 10 * time.Second
	t := time.NewTimer(interval)
	defer t.Stop()
	var info *entity.RoomInitInfo
	var preStatus int

	logx.Info("正在检测直播间是否开播...")

	// 循环监听直播间情况
	for l.running {
		select {

		// 程序退出
		case <-workctx.Done():
			l.cls.StopWsClient()
			l.running = false
			logx.Info("弹幕姬进程停止")

		// 每10秒检查一次直播间是否开播
		case <-t.C:
			if info, err = http.RoomInit(l.svcCtx.Config.RoomId); err != nil {
				logx.Infof("RoomInit错误：%v", err)
				t.Reset(interval)
				continue
			}
			if info.Data.LiveStatus == entity.Live && preStatus == entity.NotStarted { // 由NotStarted到Live是开播
				logx.Infof("开播啦！%v", l.svcCtx.Config.RoomId)

				l.cls = handler.NewWsHandler()
				if l.cls == nil {
					t.Reset(interval)
					continue
				}
				err := l.cls.StartWsClient()
				if err != nil {
					logx.Error(err)
					t.Reset(interval)
					continue
				}
				preStatus = entity.Live
				// 开启弹幕姬
			} else if info.Data.LiveStatus != entity.Live && preStatus == entity.Live { // 由Live到NotStarted是下播
				logx.Info("下播啦！")
				preStatus = entity.NotStarted
				l.cls.StopWsClient()
				//l.danmustop()
			}
			t.Reset(interval)
		}
	}

	return
}
func (l *Program) userlogin() error {
	var err error
	http.InitHttpClient()
	var loginUrl *entity.LoginUrl
	if loginUrl, err = http.GetLoginUrl(); err != nil {
		logx.Error(err)
		return err
	}

	if err = utiles.GenerateQr(loginUrl.Data.Url); err != nil {
		logx.Error(err)
		return err
	}

	if _, err = http.GetLoginInfo(loginUrl.Data.OauthKey); err != nil {
		logx.Error(err)
		return err
	}

	return err
}
