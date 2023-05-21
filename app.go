package main

import (
	"bili_danmaku/internal/http"
	entity "bili_danmaku/internal/types"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
	"strings"
	"time"
)

// App struct
type App struct {
	ctx         context.Context
	login       chan bool
	loginurl    *entity.LoginUrl
	loginstatus int
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
func (l *App) Userlogin() *entity.LoginUrl {
	var err error
	http.InitHttpClient()
	if l.login != nil {
		l.Stopwork()
	}
	l.login = make(chan bool)

	if l.loginurl, err = http.GetLoginUrl(); err != nil {
		logx.Error(err)
		return nil
	}
	go l.work()
	return l.loginurl
	//if err = utiles.GenerateQr(loginUrl.Data.Url); err != nil {
	//	logx.Error(err)
	//	return err
	//}
}
func (l *App) Getlogin() bool {
	if l.loginstatus == 1 {
		l.Stopwork()
		return true
	}
	return false
}
func (l *App) work() {
	var err error
	var url = "https://passport.bilibili.com/qrcode/getLoginInfo?oauthKey=" + l.loginurl.Data.OauthKey
	var resp *resty.Response
	var data *entity.LoginInfoData
	var file *os.File
	var CookieStr string
	var CookieList = make(map[string]string)
	cli := resty.New()
	pre := &entity.LoginInfoPre{}
	logx.Info("等待扫码登录...")
	for {
		select {
		case <-l.login:
			fmt.Println("Worker stopped")
			return
		default:
			if l.loginstatus == 1 {
				logx.Info("登录已成功等待线程退出")
				time.Sleep(5 * time.Second)
				continue
			}
			logx.Info("等待扫码中")

			if resp, err = cli.R().
				SetHeader("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.64 Safari/537.36").
				Post(url); err != nil {
				logx.Error("请求getLoginInfo失败：", err)
				l.loginstatus = 3
			}

			if err = json.Unmarshal(resp.Body(), pre); err != nil {
				logx.Error("Unmarshal失败：", err, "body:", string(resp.Body()))
				l.loginstatus = 3
			}

			if pre.Status {

				data = &entity.LoginInfoData{}
				if err = json.Unmarshal(resp.Body(), data); err != nil {
					logx.Error("Unmarshal失败：", err, "body:", string(resp.Body()))
					l.loginstatus = 3
				}
				logx.Info("登录成功！")
				for _, v := range resp.Header().Values("Set-Cookie") {
					pair := strings.Split(v, ";")
					kv := strings.Split(pair[0], "=")
					CookieList[kv[0]] = kv[1]
					CookieStr += pair[0] + ";"
				}
				//使用追加模式打开文件
				file, err = os.OpenFile("token/bili_token.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
				if err != nil {
					logx.Errorf("打开文件错误：", err)
				}
				file.WriteString(CookieStr)
				file.Close()
				//使用追加模式打开文件
				file, err = os.OpenFile("token/bili_token.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
				if err != nil {
					logx.Errorf("打开文件错误：", err)
					l.loginstatus = 3
				}
				tokenstr, _ := json.Marshal(CookieList)
				file.WriteString(string(tokenstr))
				file.Close()
				l.loginstatus = 1
				//l.Stopwork()
			}

			time.Sleep(5 * time.Second)
		}
	}
}
func (l *App) Stopwork() {
	l.login <- true
}
