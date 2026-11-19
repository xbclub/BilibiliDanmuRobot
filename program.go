package main

import (
	"context"
	"fmt"
	"github.com/xbclub/BilibiliDanmuRobot-Core/config"
	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
	"github.com/xbclub/BilibiliDanmuRobot-Core/handler"
	"github.com/xbclub/BilibiliDanmuRobot-Core/biliapi"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
	"github.com/xbclub/BilibiliDanmuRobot-Core/qrcode"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
	"sync"
	"time"
)

type Program struct {
	running    bool
	svcCtx     *svc.ServiceContext
	workCtx    context.Context
	workCancel context.CancelFunc
	cls        handler.WsHandler
	mutex      *sync.Mutex
}

// emitEvent 由 App.startup 注入的 wails 事件发射器；为 nil 时（无窗口环境）静默跳过
var emitEvent func(name string, data ...interface{})

// emitLoginRequired 后端检测到「未登录但需要登录态」时通知前端弹提示
func emitLoginRequired() {
	logx.Error("请先扫码登录后再启动机器人")
	if emitEvent != nil {
		emitEvent("login:required")
	}
}

// emitLoginExpired 登录态被服务端判定失效（开播校验发现），已自动本地注销，
// 通知前端提示并回到扫码登录页
func emitLoginExpired() {
	if emitEvent != nil {
		emitEvent("login:expired")
	}
}

// loginValid 登录态两级校验：本地 cookie 齐备性快筛（完全没登录时不打请求），
// 再调服务端 nav 接口确认真实有效——cookie 可能已在服务端过期，纯本地检查拦不住。
// 每次调用产生一个网络请求，仅在低频时机使用（启动前 / 开播连接前各一次）
func loginValid() bool {
	if !biliapi.LoggedIn() {
		return false
	}
	info := biliapi.GetUserInfo()
	return info != nil && info.Islogin
}

func NewProgram() *Program {
	return &Program{cls: handler.NewWsHandler(),
		mutex: new(sync.Mutex),
	}
}

func (p *Program) Start() bool {
	p.cls = handler.NewWsHandler()
	if p.running {
		logx.Info("已启动跳过")
		return true
	}
	// 启动前登录校验（服务端确认）：未登录/已过期直接拒绝。弹幕 WS 游客也能
	// 连上，但 EntryMsg 与后续所有弹幕发送都需要 bili_jct 登录态，
	// 未登录时连上只是"看起来正常"空跑
	if !loginValid() {
		logx.Error("启动失败：登录态无效（未登录或已过期），请先扫码登录")
		emitLoginRequired()
		return false
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

// Stop 手动停止：只停 WS 连接与监控循环（Done 分支即 StopWsClient 语义）。
// 弹幕处理各模块/发送队列不回收——没有 WS 事件流入时它们空转待命，
// 再次 Start 时无需整体重建
func (p *Program) Stop() bool {
	p.mutex.Lock()
	if p.workCancel != nil {
		p.workCancel()
		p.workCancel = nil
		p.mutex.Unlock()
		// 有界等待监控 goroutine 退出（无 token 早退路径不会置 running=false，不能死等）
		deadline := time.Now().Add(5 * time.Second)
		for p.running && time.Now().Before(deadline) {
			time.Sleep(50 * time.Millisecond)
		}
		if p.running {
			// 监控 goroutine 大概率卡在开播检测分支：此时强置 running=false，
			// goroutine 会在下一次 select 命中 Done 自行退出，但期间 Start 会双跑，记日志备查
			logx.Errorf("等待监控 goroutine 退出超时（5s），强制标记为已停止；期间请勿 Start（会双跑）")
		}
		p.mutex.Lock()
	}
	p.running = false
	p.mutex.Unlock()
	logx.Info("弹幕姬已停止（WS 连接已断开）")
	return true
}

// Logout 退出登录：与普通 Stop 不同，这里必须全面停止——
// 1. Stop：断开 WS 连接（连接持有旧 cookie 快照，注销前必须先停）与监控循环
// 2. StopChanel：停掉弹幕解析/欢迎/感谢/机器人等处理模块、中央发送队列、定时弹幕
// 3. Core Logout：best-effort 服务端注销 → 清空进程内 cookie → 删除本地凭证
// 完成后前端回到扫码登录页重新走登录流程
func (p *Program) Logout() error {
	p.Stop()
	p.mutex.Lock()
	if p.cls != nil {
		p.cls.StopChanel()
		// 兜底再停一次 WS：Stop() 走满 5s 超时路径时 Done 分支可能尚未执行，
		// 而「注销前必须断开持有旧 cookie 的连接」是硬性顺序。
		// client.Stop() 为 context cancel、corndanmu.Stop() 幂等，重复调用安全
		p.cls.StopWsClient()
	}
	p.mutex.Unlock()
	return biliapi.Logout()
}
func (p *Program) Restart() bool {
	p.mutex.Lock()
	if p.cls != nil {
		err := p.cls.ReloadConfig()
		if err != nil {
			logx.Error(err)
			return false
		}
		logx.Info("重载配置已完成")
	}
	p.mutex.Unlock()
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
	if !biliapi.FileExists("token/bili_token.txt") || !biliapi.FileExists("token/bili_token.json") {
		return
	}
	var err error
	biliapi.InitHttpClient()

	// 准备select中用到的变量
	//sig := make(chan os.Signal)
	//signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	var interval = 10 * time.Second
	t := time.NewTimer(interval)
	defer t.Stop()
	var info *entity.RoomInitInfo
	var preStatus int

	logx.Info("正在检测直播间是否开播...")
	l.cls = handler.NewWsHandler()
	l.cls.InitStartWsClient()
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
			if info, err = biliapi.RoomInit(l.cls.GetSvc().Config.RoomId); err != nil {
				logx.Infof("RoomInit错误：%v", err)
				t.Reset(interval)
				continue
			}
			if info.Data.LiveStatus == entity.Live && preStatus == entity.NotStarted { // 由NotStarted到Live是开播
				// 开播连接前登录校验（服务端确认，仅此一次，不轮询）。
				// 能走到这里说明启动时登录有效，此刻失效即登录已过期——
				// 直接走本地注销流程：停全部任务、清 cookie、删凭证、回登录页；
				// 过期会话调服务端注销注定失败，跳过该接口（LogoutLocal）
				if !loginValid() {
					logx.Errorf("开播时校验登录态已失效，停止机器人并执行本地注销")
					emitLoginExpired()
					l.running = false
					if l.cls != nil {
						l.cls.StopChanel()
						l.cls.StopWsClient()
					}
					biliapi.LogoutLocal()
					return
				}
				logx.Infof("开播啦！%v", l.cls.GetSvc().Config.RoomId)
				l.mutex.Lock()
				l.cls.StopWsClient()
				err := l.cls.StartWsClient()
				if err != nil {
					logx.Error(err)
					t.Reset(interval)
					continue
				}
				preStatus = entity.Live
				l.mutex.Unlock()
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
	biliapi.InitHttpClient()
	var loginUrl *entity.LoginUrl
	if loginUrl, err = biliapi.GetLoginUrl(); err != nil {
		logx.Error(err)
		return err
	}

	if err = qrcode.Generate(loginUrl.Data.Url); err != nil {
		logx.Error(err)
		return err
	}

	if _, err = biliapi.GetLoginInfo(loginUrl.Data.OauthKey); err != nil {
		logx.Error(err)
		return err
	}

	return err
}
