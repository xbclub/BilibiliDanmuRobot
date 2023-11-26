package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/xbclub/BilibiliDanmuRobot-Core/config"
	"github.com/xbclub/BilibiliDanmuRobot-Core/entity"
	"github.com/xbclub/BilibiliDanmuRobot-Core/handler"
	"github.com/xbclub/BilibiliDanmuRobot-Core/http"
	"github.com/xbclub/BilibiliDanmuRobot-Core/logic"
	"github.com/xbclub/BilibiliDanmuRobot-Core/svc"
	"github.com/xbclub/BilibiliDanmuRobot-Core/utiles"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var configFile = flag.String("f", "etc/bilidanmaku-api.yaml", "the config file")
var Version string
var cls handler.WsHandler

func main() {
	flag.Parse()
	var err error
	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())
	logx.MustSetup(c.Log)
	logx.DisableStat()
	logx.Infof("当前版本: %s", Version)
	logx.Infof("监听直播间: %d", c.RoomId)
	dir := "./token"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// Directory does not exist, create it
		err = os.Mkdir(dir, 0755)
		if err != nil {
			panic(fmt.Sprintf("无法创建token文件夹 请手动创建:%s", err))
		}
	}
	http.InitHttpClient()
	// 判断是否存在历史cookie
	if http.FileExists("token/bili_token.txt") && http.FileExists("token/bili_token.json") {
		err = http.SetHistoryCookie()
		if err != nil {
			if err = userlogin(); err != nil {
				logx.Errorf("用户登录失败：%v", err)
				return
			}
		}
		logx.Info("用户登录成功")
	} else {
		if err = userlogin(); err != nil {
			logx.Errorf("用户登录失败：%v", err)
			return
		}
		logx.Info("用户登录成功")
	}
	ctx := svc.NewServiceContext(c)
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	var interval = 10 * time.Second
	t := time.NewTimer(interval)
	defer t.Stop()
	var info *entity.RoomInitInfo
	var preStatus int

	input := make(chan string)
	if ctx.Config.CustomizeBullet {
		logx.Info("启动命令行输入线程...")
		go getTerminalInput(input)
	}
	logx.Info("正在检测直播间是否开播...")
	for {
		select {

		// 程序退出
		case <-sig:
			goto END

		case ta := <-input:
			if ctx.Config.CustomizeBullet {
				logic.PushToBulletSender(ta)
			}

		// 每1分钟检查一次直播间是否开播
		case <-t.C:
			t.Reset(interval)
			if info, err = http.RoomInit(ctx.Config.RoomId); err != nil || err != nil {
				logx.Infof("RoomInit错误：%v", err)
				continue
			}
			if info.Data.LiveStatus == entity.Live && preStatus == entity.NotStarted { // 由NotStarted到Live是开播
				logx.Infof("开播啦！%v", ctx.Config.RoomId)
				preStatus = entity.Live
				cls = handler.NewWsHandler()
				if cls == nil {
					os.Exit(1)
				}
				cls.StartWsClient() // 开启弹幕姬
			} else if info.Data.LiveStatus != entity.Live && preStatus == entity.Live { // 由Live到NotStarted是下播
				logx.Info("下播啦！")
				preStatus = entity.NotStarted
				cls.StopWsClient()
			}
		}
	}
END:

	return
}
func userlogin() error {
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
func getTerminalInput(input chan string) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		input <- scanner.Text()
	}
}
