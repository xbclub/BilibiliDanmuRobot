package main

import (
	"bili_danmaku/internal/bullet_girl"
	"bili_danmaku/internal/config"
	"bili_danmaku/internal/errs"
	"bili_danmaku/internal/http"
	"bili_danmaku/internal/svc"
	types2 "bili_danmaku/internal/types"
	"bili_danmaku/internal/utiles"
	"context"
	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/logx"
	"math/rand"
	"time"
)

type Program struct {
	sendBulletCtx       context.Context
	sendBulletCancel    context.CancelFunc
	timingBulletCancel  context.CancelFunc
	robotBulletCtx      context.Context
	robotBulletCancel   context.CancelFunc
	catchBulletCtx      context.Context
	catchBulletCancel   context.CancelFunc
	handleBulletCtx     context.Context
	handleBulletCancel  context.CancelFunc
	thanksGiftCtx       context.Context
	thankGiftCancel     context.CancelFunc
	ineterractCtx       context.Context
	ineterractCancel    context.CancelFunc
	pkCtx               context.Context
	pkCancel            context.CancelFunc
	corndanmu           *cron.Cron
	mapCronDanmuSendIdx map[int]int
	running             bool
	svcCtx              *svc.ServiceContext
	workCtx             context.Context
	workCancel          context.CancelFunc
}

func NewProgram() *Program {
	return &Program{
		mapCronDanmuSendIdx: make(map[int]int),
	}
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
	var err error
	http.InitHttpClient()

	// 准备select中用到的变量
	//sig := make(chan os.Signal)
	//signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	var interval = 10 * time.Second
	t := time.NewTimer(interval)
	defer t.Stop()
	var info *types2.RoomInitInfo
	var preStatus int

	input := make(chan string)
	logx.Info("正在检测直播间是否开播...")

	// 循环监听直播间情况
	for l.running {
		select {

		// 程序退出
		case <-workctx.Done():
			if l.sendBulletCancel != nil {
				l.sendBulletCancel()
			}
			if l.timingBulletCancel != nil {
				l.timingBulletCancel()
			}
			if l.robotBulletCancel != nil {
				l.robotBulletCancel()
			}
			if l.catchBulletCancel != nil {
				l.catchBulletCancel()
			}
			if l.handleBulletCancel != nil {
				l.handleBulletCancel()
			}
			if l.thankGiftCancel != nil {
				if l.svcCtx.Config.ThanksGift {
					l.thankGiftCancel()
				}
			}
			if l.ineterractCancel != nil {
				l.ineterractCancel()
			}
			if l.pkCancel != nil {
				l.pkCancel()
			}
			if l.corndanmu != nil {
				l.danmustop()
			}
			l.running = false
			logx.Info("弹幕姬进程停止")

		case ta := <-input:
			if l.svcCtx.Config.CustomizeBullet {
				bullet_girl.PushToBulletSender(ta)
			}

		// 每1分钟检查一次直播间是否开播
		case <-t.C:
			t.Reset(interval)
			if info, err = http.RoomInit(l.svcCtx.Config.RoomId); err != nil || err == errs.RoomIdNotExistErr {
				logx.Infof("RoomInit错误：%v", err)
				continue
			}
			if info.Data.LiveStatus == types2.Live && preStatus == types2.NotStarted { // 由NotStarted到Live是开播
				logx.Infof("开播啦！%v", l.svcCtx.Config.RoomId)
				preStatus = types2.Live
				// 弹幕姬各goroutine上下文
				l.sendBulletCtx, l.sendBulletCancel = context.WithCancel(context.Background())
				//timingBulletCtx, timingBulletCancel = context.WithCancel(context.Background())
				l.robotBulletCtx, l.robotBulletCancel = context.WithCancel(context.Background())
				l.catchBulletCtx, l.catchBulletCancel = context.WithCancel(context.Background())
				l.handleBulletCtx, l.handleBulletCancel = context.WithCancel(context.Background())
				l.thanksGiftCtx, l.thankGiftCancel = context.WithCancel(context.Background())
				l.ineterractCtx, l.ineterractCancel = context.WithCancel(context.Background())
				l.pkCtx, l.pkCancel = context.WithCancel(context.Background())
				l.StartBulletGirl(l.sendBulletCtx,
					//timingBulletCtx,
					l.robotBulletCtx, l.catchBulletCtx, l.handleBulletCtx, l.thanksGiftCtx) // 开启弹幕姬
			} else if info.Data.LiveStatus == types2.NotStarted && preStatus == types2.Live { // 由Live到NotStarted是下播
				logx.Info("下播啦！")
				preStatus = types2.NotStarted
				if l.sendBulletCancel != nil {
					l.sendBulletCancel()
				}
				if l.timingBulletCancel != nil {
					l.timingBulletCancel()
				}
				if l.robotBulletCancel != nil {
					l.robotBulletCancel()
				}
				if l.catchBulletCancel != nil {
					l.catchBulletCancel()
				}
				if l.handleBulletCancel != nil {
					l.handleBulletCancel()
				}
				if l.thankGiftCancel != nil {
					if l.svcCtx.Config.ThanksGift {
						l.thankGiftCancel()
					}
				}
				if l.ineterractCancel != nil {
					l.ineterractCancel() // 关闭弹幕姬goroutine
				}
				if l.pkCancel != nil {
					l.pkCancel()
				}
				if l.corndanmu != nil {
					l.danmustop()
				}
			}
		}
	}

	return
}
func (l *Program) Bili_danmaku_Stop() {
	defer func() {
		if l.sendBulletCancel != nil {
			l.sendBulletCancel()
		}
		if l.timingBulletCancel != nil {
			l.timingBulletCancel()
		}
		if l.robotBulletCancel != nil {
			l.robotBulletCancel()
		}
		if l.catchBulletCancel != nil {
			l.catchBulletCancel()
		}
		if l.handleBulletCancel != nil {
			l.handleBulletCancel()
		}
		if l.thankGiftCancel != nil {
			if l.svcCtx.Config.ThanksGift {
				l.thankGiftCancel()
			}
		}
		if l.ineterractCancel != nil {
			l.ineterractCancel()
		}
		if l.pkCancel != nil {
			l.pkCancel()
		}
		if l.corndanmu != nil {
			l.danmustop()
		}
	}()
}
func (l *Program) userlogin() error {
	var err error
	http.InitHttpClient()
	var loginUrl *types2.LoginUrl
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
func (l *Program) StartBulletGirl(sendBulletCtx,
	//timingBulletCtx,
	robotBulletCtx, catchBulletCtx, handleBulletCtx, thanksGiftCtx context.Context) {
	if l.svcCtx.Config.EntryMsg != "off" {
		err := http.Send(l.svcCtx.Config.EntryMsg, l.svcCtx)
		if err != nil {
			logx.Error(err)
		}
	}
	// 开启弹幕推送
	go bullet_girl.StartSendBullet(sendBulletCtx, l.svcCtx)
	logx.Info("弹幕推送已开启...")

	// 开启定时弹幕任务
	//go bullet_girl.StartTimingBullet(timingBulletCtx)
	//logx.Info("定时弹幕已开启...")

	// 开启弹幕机器人
	go bullet_girl.StartBulletRobot(robotBulletCtx, l.svcCtx)
	logx.Info("弹幕机器人已开启")

	// 开启弹幕抓取
	go bullet_girl.StartCatchBullet(catchBulletCtx, l.svcCtx)
	logx.Info("弹幕抓取已开启...")

	// 开启弹幕处理
	go bullet_girl.HandleBullet(handleBulletCtx, l.svcCtx)
	logx.Info("弹幕处理已开启...")
	// 开启pk
	if l.svcCtx.Config.PKNotice {
		go bullet_girl.PK(l.pkCtx)
	}
	// 开启礼物感谢
	if l.svcCtx.Config.ThanksGift {
		go bullet_girl.ThanksGift(thanksGiftCtx)
		logx.Info("礼物感谢已开启")
	}

	go bullet_girl.Interact(l.ineterractCtx)
	if l.svcCtx.Config.CronDanmu {
		l.danmustart()
		logx.Info("定时弹幕已开启")
	}
}

// 定时弹幕功能
func (l *Program) danmustart() {
	if l.svcCtx.Config.CronDanmuList != nil && len(l.svcCtx.Config.CronDanmuList) > 0 {
		if l.corndanmu == nil {
			l.corndanmu = cron.New(cron.WithParser(cron.NewParser(
				cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow,
			)))
		}
		rand.Seed(time.Now().UnixNano())
		for i, danmu := range l.svcCtx.Config.CronDanmuList {
			if danmu.Danmu != nil {
				danmus := danmu
				_, err := l.corndanmu.AddFunc(danmus.Cron, func() {
					if len(danmus.Danmu) > 0 {
						if danmus.Random {
							bullet_girl.PushToBulletSender(danmus.Danmu[rand.Intn(len(danmus.Danmu))])
						} else {
							_, ok := l.mapCronDanmuSendIdx[i]
							if !ok {
								l.mapCronDanmuSendIdx[i] = 0
							}
							l.mapCronDanmuSendIdx[i] = l.mapCronDanmuSendIdx[i] + 1
							bullet_girl.PushToBulletSender(danmus.Danmu[l.mapCronDanmuSendIdx[i]%len(danmus.Danmu)])
						}
					}
				})
				if err != nil {
					logx.Errorf("第%d条定时弹幕配置出现错误: %v", i+1, err)
				}
			}
		}
		l.corndanmu.Start()
	}
}

func (l *Program) danmustop() {
	l.corndanmu.Stop()
	l.corndanmu = nil
	logx.Info("定时弹幕已关闭")
}

//func main() {
//	program := NewProgram()
//
//	err := program.Start()
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	program.monitor()
//	time.Sleep(time.Second * 10) // 模拟程序运行10秒钟
//
//	err = program.Stop()
//	program.monitor()
//	fmt.Println(program.buf)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//}
