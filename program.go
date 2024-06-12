package main

import (
	"context"
	"fmt"
	"github.com/xbclub/BilibiliDanmuRobot-Core/config"
	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
	"github.com/xbclub/BilibiliDanmuRobot-Core/handler"
	"github.com/xbclub/BilibiliDanmuRobot-Core/http"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
	"github.com/xbclub/BilibiliDanmuRobot-Core/utiles"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
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
	return &Program{cls: handler.NewWsHandler()}
}

func (p *Program) Start() bool {
	p.cls = handler.NewWsHandler()
	if p.running {
		logx.Info("已启动跳过")
		return true
	}
	var c config.Config
	if err := Mustload(&c); err != nil {
		logx.Error(err)
		return false
	}
	logx.MustSetup(c.Log)
	logx.DisableStat()
	dbdir, err := os.Stat(c.DBPath)
	if os.IsNotExist(err) || !dbdir.IsDir() {
		err = os.MkdirAll(c.DBPath, 0777)
		if err != nil {
			logx.Errorf("数据库文件夹创建失败：%s", c.DBPath)
			panic(fmt.Sprintf("无法创建数据库文件夹 请手动创建:%s", err))
		}
	}
	p.svcCtx = svc.NewServiceContext(c)
	p.workCtx, p.workCancel = context.WithCancel(context.Background())
	p.running = true
	go p.Bili_danmaku_Start(p.workCtx)
	return true
}

func (p *Program) Stop() bool {
	//p.locked.Lock()
	//if p.workCancel != nil {
	//	p.workCancel()
	//}
	//for p.running {
	//	time.Sleep(1 * time.Second)
	//	logx.Info("等待机器关闭中.....")
	//}
	logx.Info("手动停止已经不让用了 别请求了")
	//p.locked.Unlock()
	return false
}
func (p *Program) Restart() bool {
	if p.cls != nil {
		err := p.cls.ReloadConfig()
		if err != nil {
			logx.Error(err)
			return false
		}
		logx.Info("重载配置已完成")
	}

	//if err != nil {
	//	logx.Error(err)
	//	return false
	//}
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
			if l.cls != nil {
				l.cls.StopWsClient()
			}
			l.running = false
			logx.Info("弹幕姬进程停止")

		// 每10秒检查一次直播间是否开播
		case <-t.C:
			if info, err = http.RoomInit(l.cls.GetSvc().Config.RoomId); err != nil {
				logx.Infof("RoomInit错误：%v", err)
				t.Reset(interval)
				continue
			}
			if info.Data.LiveStatus == entity.Live && preStatus == entity.NotStarted { // 由NotStarted到Live是开播
				logx.Infof("开播啦！%v", l.svcCtx.Config.RoomId)
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
				l.cls.SayGoodbye()
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
