package logic

import (
	"bili_danmaku/internal/bullet_girl"
	"bili_danmaku/internal/errs"
	"bili_danmaku/internal/http"
	"bili_danmaku/internal/svc"
	entity "bili_danmaku/internal/types"
	"bili_danmaku/internal/utiles"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var sendBulletCtx context.Context
var sendBulletCancel context.CancelFunc
var timingBulletCtx context.Context
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

func (l *Bili_danmakuLogic) Bili_danmaku_Start() {
	// todo: add your logic here and delete this line
	var err error
	http.InitHttpClient()
	// 判断是否存在历史cookie
	if http.FileExists("token/bili_token.txt") && http.FileExists("token/bili_token.json") {
		err = http.SetHistoryCookie()
		if err != nil {
			if err = l.userlogin(); err != nil {
				logx.Errorf("用户登录失败：", err)
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

	log.Println("正在检测直播间是否开播...")

	// 循环监听直播间情况
	for {
		select {

		// 程序退出
		case <-sig:
			goto END

		// 每1分钟检查一次直播间是否开播
		case <-t.C:
			t.Reset(interval)
			if info, err = http.RoomInit(l.svcCtx); err != nil || err == errs.RoomIdNotExistErr {
				log.Println("RoomInit错误：", err)
				continue
			}
			if info.Data.LiveStatus == entity.Live && preStatus == entity.NotStarted { // 由NotStarted到Live是开播
				log.Println("开播啦！", l.svcCtx.Config.RoomId)
				preStatus = entity.Live
				// 弹幕姬各goroutine上下文
				sendBulletCtx, sendBulletCancel = context.WithCancel(context.Background())
				timingBulletCtx, timingBulletCancel = context.WithCancel(context.Background())
				robotBulletCtx, robotBulletCancel = context.WithCancel(context.Background())
				catchBulletCtx, catchBulletCancel = context.WithCancel(context.Background())
				handleBulletCtx, handleBulletCancel = context.WithCancel(context.Background())
				thanksGiftCtx, thankGiftCancel = context.WithCancel(context.Background())
				ineterractCtx, ineterractCancel = context.WithCancel(context.Background())
				l.StartBulletGirl(sendBulletCtx, timingBulletCtx, robotBulletCtx, catchBulletCtx, handleBulletCtx, thanksGiftCtx) // 开启弹幕姬
			} else if info.Data.LiveStatus == entity.NotStarted && preStatus == entity.Live { // 由Live到NotStarted是下播
				log.Println("下播啦！")
				preStatus = entity.NotStarted
				sendBulletCancel()
				timingBulletCancel()
				robotBulletCancel()
				catchBulletCancel()
				handleBulletCancel()
				thankGiftCancel()
				ineterractCancel() // 关闭弹幕姬goroutine
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
			thankGiftCancel()
		}
		if ineterractCancel != nil {
			ineterractCancel()
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
func (l *Bili_danmakuLogic) StartBulletGirl(sendBulletCtx, timingBulletCtx, robotBulletCtx, catchBulletCtx, handleBulletCtx, thanksGiftCtx context.Context) {
	if l.svcCtx.Config.EntryMsg != "off" {
		http.Send(l.svcCtx.Config.EntryMsg, l.svcCtx)
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
	go bullet_girl.ThanksGift(thanksGiftCtx)
	logx.Info("礼物感谢已开启")
	go bullet_girl.Interact(ineterractCtx)
	// 指定弹幕定时任务
	//time.Sleep(time.Second) // 现开启定时任务弹幕再推送，这个方法很low，暂且这样吧
	//bullet_girl.PushToBulletEvent(
	//	bullet_girl.NewBulletEvent(
	//		bullet_girl.Save, bullet_girl.NewBulletTask(
	//			bullet_girl.NewBullet("ios请到哔哩哔哩直播姬公众号投喂哦～", "*/9 * * * * *"))))
	//bullet_girl.PushToBulletEvent(
	//	bullet_girl.NewBulletEvent(
	//		bullet_girl.Save, bullet_girl.NewBulletTask(
	//			bullet_girl.NewBullet("喜欢主播可以加入粉丝团哦～", "*/7 * * * * *"))))
	//bullet_girl.PushToBulletEvent(
	//	bullet_girl.NewBulletEvent(
	//		bullet_girl.Save, bullet_girl.NewBulletTask(
	//			bullet_girl.NewBullet("主播今天很可爱哦！干巴爹！", "*/17 * * * * *"))))
	//bullet_girl.PushToBulletEvent(
	//	bullet_girl.NewBulletEvent(
	//		bullet_girl.Save, bullet_girl.NewBulletTask(
	//			bullet_girl.NewBullet("无聊的同学可以找橘子聊天喔！", "*/23 * * * * *"))))
}
