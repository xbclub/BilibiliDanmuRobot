package logic

import (
	"bili_danmaku/internal/bullet_girl"
	"bili_danmaku/internal/errs"
	"bili_danmaku/internal/http"
	"bili_danmaku/internal/svc"
	entity "bili_danmaku/internal/types"
	"bili_danmaku/internal/utiles"
	"bufio"
	"context"
	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var sendBulletCtx context.Context
var sendBulletCancel context.CancelFunc

// var timingBulletCtx context.Context
var timingBulletCancel context.CancelFunc
var robotBulletCtx context.Context
var robotBulletCancel context.CancelFunc
var catchBulletCtx context.Context
var catchBulletCancel context.CancelFunc
var handleBulletCtx context.Context
var handleBulletCancel context.CancelFunc
var thanksGiftCtx context.Context
var thankGiftCancel context.CancelFunc
var ineterractCtx context.Context
var ineterractCancel context.CancelFunc
var corndanmu *cron.Cron

type Bili_danmakuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBili_danmakuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Bili_danmakuLogic {
	return &Bili_danmakuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func getTerminalInput(input chan string) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		input <- scanner.Text()
	}
}

func (l *Bili_danmakuLogic) Bili_danmaku_Start() {
	var err error
	http.InitHttpClient()
	// 判断是否存在历史cookie
	if http.FileExists("token/bili_token.txt") && http.FileExists("token/bili_token.json") {
		err = http.SetHistoryCookie()
		if err != nil {
			if err = l.userlogin(); err != nil {
				logx.Errorf("用户登录失败：%v", err)
				return
			}
		}
		logx.Info("用户登录成功")
	} else {
		if err = l.userlogin(); err != nil {
			logx.Errorf("用户登录失败：", err)
			return
		}
		logx.Info("用户登录成功")
	}

	// 准备select中用到的变量
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	var interval = 10 * time.Second
	t := time.NewTimer(interval)
	defer t.Stop()
	var info *entity.RoomInitInfo
	var preStatus int

	input := make(chan string)
	if l.svcCtx.Config.CustomizeBullet {
		logx.Info("启动命令行输入线程...")
		go getTerminalInput(input)
	}
	logx.Info("正在检测直播间是否开播...")

	// 循环监听直播间情况
	for {
		select {

		// 程序退出
		case <-sig:
			goto END

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
			if info.Data.LiveStatus == entity.Live && preStatus == entity.NotStarted { // 由NotStarted到Live是开播
				logx.Infof("开播啦！%v", l.svcCtx.Config.RoomId)
				preStatus = entity.Live
				// 弹幕姬各goroutine上下文
				sendBulletCtx, sendBulletCancel = context.WithCancel(context.Background())
				//timingBulletCtx, timingBulletCancel = context.WithCancel(context.Background())
				robotBulletCtx, robotBulletCancel = context.WithCancel(context.Background())
				catchBulletCtx, catchBulletCancel = context.WithCancel(context.Background())
				handleBulletCtx, handleBulletCancel = context.WithCancel(context.Background())
				thanksGiftCtx, thankGiftCancel = context.WithCancel(context.Background())
				ineterractCtx, ineterractCancel = context.WithCancel(context.Background())
				l.StartBulletGirl(sendBulletCtx,
					//timingBulletCtx,
					robotBulletCtx, catchBulletCtx, handleBulletCtx, thanksGiftCtx) // 开启弹幕姬
			} else if info.Data.LiveStatus == entity.NotStarted && preStatus == entity.Live { // 由Live到NotStarted是下播
				logx.Info("下播啦！")
				preStatus = entity.NotStarted
				if sendBulletCancel != nil {
					sendBulletCancel()
				}
				if timingBulletCancel != nil {
					timingBulletCancel()
				}
				if robotBulletCancel != nil {
					robotBulletCancel()
				}
				if catchBulletCancel != nil {
					catchBulletCancel()
				}
				if handleBulletCancel != nil {
					handleBulletCancel()
				}
				if thankGiftCancel != nil {
					if l.svcCtx.Config.ThanksGift {
						thankGiftCancel()
					}
				}
				if ineterractCancel != nil {
					ineterractCancel() // 关闭弹幕姬goroutine
				}
				if corndanmu != nil {
					l.danmustop()
				}
			}
		}
	}
END:

	return
}
func (l *Bili_danmakuLogic) Bili_danmaku_Stop() {
	defer func() {
		if sendBulletCancel != nil {
			sendBulletCancel()
		}
		if timingBulletCancel != nil {
			timingBulletCancel()
		}
		if robotBulletCancel != nil {
			robotBulletCancel()
		}
		if catchBulletCancel != nil {
			catchBulletCancel()
		}
		if handleBulletCancel != nil {
			handleBulletCancel()
		}
		if thankGiftCancel != nil {
			if l.svcCtx.Config.ThanksGift {
				thankGiftCancel()
			}
		}
		if ineterractCancel != nil {
			ineterractCancel()
		}
		if corndanmu != nil {
			l.danmustop()
		}
	}()
}
func (l *Bili_danmakuLogic) userlogin() error {
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
func (l *Bili_danmakuLogic) StartBulletGirl(sendBulletCtx,
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

	// 开启礼物感谢
	if l.svcCtx.Config.ThanksGift {
		go bullet_girl.ThanksGift(thanksGiftCtx)
		logx.Info("礼物感谢已开启")
	}

	go bullet_girl.Interact(ineterractCtx)
	if l.svcCtx.Config.CronDanmu {
		l.danmustart()
		logx.Info("定时弹幕已开启")
	}
}

// 定时弹幕功能
func (l *Bili_danmakuLogic) danmustart() {
	corndanmu = cron.New(cron.WithParser(cron.NewParser(
		cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow,
	)))
	for i, danmu := range l.svcCtx.Config.CronDanmuList {
		danmus := danmu
		_, err := corndanmu.AddFunc(danmus.Cron, func() {
			bullet_girl.PushToBulletSender(danmus.Danmu)
		})
		if err != nil {
			logx.Errorf("第%d条定时弹幕配置出现错误: %v", i+1, err)
		}
	}
	corndanmu.Start()
}
func (l *Bili_danmakuLogic) danmustop() {
	corndanmu.Stop()
	corndanmu = nil
	logx.Info("定时弹幕已关闭")
}
